package megabus

import (
	"encoding/json"
	"fmt"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	http "github.com/saucesteals/fhttp"
	"github.com/tidwall/gjson"
	"io"
	"strconv"
	"strings"
	"time"
)

func header(input ...string) []string {
	var output []string
	for _, i := range input {
		output = append(output, i)
	}
	return output
}

func (r *megabusWorker) fetchResults(query *shared.AggregatorQuery) ([]*shared.AggregatorResponse, error) {
	originID := r.findClosestNode(query.StartLocation.Latitude, query.StartLocation.Longitude)
	destinationID := r.findClosestNode(query.EndLocation.Latitude, query.EndLocation.Longitude)

	r.statusC(fmt.Sprintf("Fetching Results (%d -> %d)", originID, destinationID))

	req, err := http.NewRequest("GET", "https://us.megabus.com/journey-planner/api/journeys?originId="+strconv.Itoa(originID)+"&destinationId="+strconv.Itoa(destinationID)+"&departureDate=2022-10-25&totalPassengers=1&concessionCount=0&nusCount=0&otherDisabilityCount=0&wheelchairSeated=0&pcaCount=0&days=1", nil)
	if err != nil {
		return nil, err
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("Host", "us.megabus.com"),
		header("Connection", "keep-alive"),
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("Accept", "application/json, text/plain, */*"),
		header("X-XSRF-Token", "zliDY-rd7e0Ac0GER6tEzOmIRQDNglP85kjcFEwBc6KWsdu9tYC2wUG25CWo29Zdwp-mjKhkg5DaMM9_W_34OuDjHiSBpa_BX6TovNmuHmI1"),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("User-Agent", r.browserConfig.UserAgent),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("Sec-Fetch-Site", "same-origin"),
		header("Sec-Fetch-Mode", "cors"),
		header("Sec-Fetch-Dest", "empty"),
		header("Referer", "https://us.megabus.com/journey-planner/journeys?days=1&concessionCount=0&departureDate=2022-10-25&destinationId=545&inboundDepartureDate=2022-10-25&inboundOtherDisabilityCount=0&inboundPcaCount=0&inboundWheelchairSeated=0&nusCount=0&originId=542&otherDisabilityCount=0&pcaCount=0&totalPassengers=1&wheelchairSeated=0"),
		header("Accept-Encoding", "gzip, deflate, br"),
		header("Accept-Language", r.browserConfig.AcceptLanguage),
	})

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error Fetching Results (%v)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseResults(body)
}

// parseDuration parses duration in Megabus format of PT4H15 to minutes.
func parseDuration(duration string) int {
	dSplit := strings.Split(duration, "PT")
	if len(dSplit) < 1 {
		return -1
	}

	hSplit := strings.Split(dSplit[1], "H")
	if len(hSplit) < 2 {
		return -1
	}

	hours, _ := strconv.Atoi(hSplit[0])
	minutes, _ := strconv.Atoi(hSplit[1])

	return (hours*60 + minutes) * 60
}

func parseResults(body []byte) ([]*shared.AggregatorResponse, error) {
	var responses []*shared.AggregatorResponse

	jsonBody := gjson.ParseBytes(body)
	jsonBody.Get("journeys").ForEach(func(key, value gjson.Result) bool {
		// Jan 2, 2006 at 3:04pm (MST)
		departureTime, _ := time.Parse("2006-01-02T15:04:00", value.Get("departureDateTime").String())
		arrivalTime, _ := time.Parse("2006-01-02T15:04:00", value.Get("arrivalDateTime").String())

		responses = append(responses, &shared.AggregatorResponse{
			Source:      "m-1",
			Description: "Megabus Route " + value.Get("routeName").String(),

			Carrier: "Megabus",
			Type:    "Bus",

			StartLocation: value.Get("origin.stopName").String(),
			EndLocation:   value.Get("destination.stopName").String(),

			StartTime: int(departureTime.Unix()),
			EndTime:   int(arrivalTime.Unix()),

			Price:    int(value.Get("price").Int()),
			Duration: parseDuration(value.Get("duration").String()),
		})

		return true
	})

	return responses, nil
}

func (r *megabusWorker) fetchNodes() error {
	r.statusC("Fetching Nodes")

	req, err := http.NewRequest("GET", "https://us.megabus.com/journey-planner/api/origin-cities", nil)
	if err != nil {
		return err
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("upgrade-insecure-requests", "1"),
		header("user-agent", r.browserConfig.UserAgent),
		header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"),
		header("sec-fetch-site", "none"),
		header("sec-fetch-mode", "navigate"),
		header("sec-fetch-user", "?1"),
		header("sec-fetch-dest", "document"),
		header("accept-encoding", "indentity"),
		header("accept-language", r.browserConfig.AcceptLanguage),
	})

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &r.nodes)
}

func (r *megabusWorker) fetchHome() error {
	r.statusC("Fetching Session")

	req, err := http.NewRequest("GET", "https://www.megabus.com/", nil)
	if err != nil {
		return err
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("upgrade-insecure-requests", "1"),
		header("user-agent", r.browserConfig.UserAgent),
		header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"),
		header("sec-fetch-site", "none"),
		header("sec-fetch-mode", "navigate"),
		header("sec-fetch-user", "?1"),
		header("sec-fetch-dest", "document"),
		header("accept-encoding", "gzip, deflate, br"),
		header("accept-language", r.browserConfig.AcceptLanguage),
	})

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Errorf("Error Fetching Session (%v)", resp.StatusCode)
	}

	return nil
}
