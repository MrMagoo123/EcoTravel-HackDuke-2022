package checkmybus

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	http "github.com/saucesteals/fhttp"
	"github.com/tidwall/gjson"
	"github.com/twinj/uuid"
	"io"
	"net/url"
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

func (r *checkmybusWorker) fetchResults(query *shared.AggregatorQuery) ([]*shared.AggregatorResponse, error) {
	r.statusC("Fetching Results")

	// help!!!
	payload := "originIsCity=" + "true" + "&" +
		"originIsAirport=" + "false" + "&" +
		"destinationIsCity=" + "true" + "&" +
		"destinationIsAirport=" + "false" + "&" +
		"latitudeFrom=" + fmt.Sprintf("%f", query.StartLocation.Latitude) + "&" +
		"longitudeFrom=" + fmt.Sprintf("%f", query.StartLocation.Longitude) + "&" +
		"nameFrom=" + url.QueryEscape(query.StartLocation.Formatted) + "&" +
		"latitudeTo=" + fmt.Sprintf("%f", query.EndLocation.Latitude) + "&" +
		"longitudeTo=" + fmt.Sprintf("%f", query.EndLocation.Longitude) + "&" +
		"nameTo=" + url.QueryEscape(query.EndLocation.Formatted) + "&" +
		"departureDate=" + "2022-10-26" + "&" +
		"queryId=" + uuid.NewV4().String() + "&" +
		"currency=USD&culture=en-US&searchRadius=15&adults=1&children0&__RequestVerificationToken=" + url.QueryEscape(r.requestVerificationToken)

	req, err := http.NewRequest("POST", "https://www.checkmybus.com/api/search", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("content-length", ""),
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("user-agent", r.browserConfig.UserAgent),
		header("content-type", "application/x-www-form-urlencoded; charset=UTF-8"),
		header("accept", "*/*"),
		header("x-requested-with", "XMLHttpRequest"),
		//header("request-id"),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("origin", "https://www.checkmybus.com"),
		header("sec-fetch-site", "same-site"),
		header("sec-fetch-mode", "cors"),
		header("sec-fetch-dest", "empty"),
		header("referer", "https://www.checkmybus.com/search"),
		header("accept-encoding", "gzip, deflate, br"),
		header("accept-language", r.browserConfig.AcceptLanguage),
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
	hSplit := strings.Split(duration, "h ")
	if len(hSplit) < 2 {
		return -1
	}

	hours, _ := strconv.Atoi(hSplit[0])
	minutes, _ := strconv.Atoi(strings.Split(hSplit[1], "m")[0])

	return (hours*60 + minutes) * 60
}

func parseResults(body []byte) ([]*shared.AggregatorResponse, error) {
	var responses []*shared.AggregatorResponse
	stations := make(map[int]string)

	jsonBody := gjson.ParseBytes(body)
	jsonBody.Get("Data.Stations").ForEach(func(key, value gjson.Result) bool {
		stations[int(value.Get("LocationCode").Int())] = strings.Title(strings.ToLower(fmt.Sprintf("%s (%s)", value.Get("Name").String(), value.Get("StreetAddress1").String())))

		return true
	})

	jsonBody.Get("ConnectionResultModels").ForEach(func(key, value gjson.Result) bool {
		value = value.Get("ConnectionResultItems.0")

		// Jan 2, 2006 at 3:04pm (MST)
		departureTime, _ := time.Parse("2006-01-02T15:04:00", value.Get("DepartureDateTime").String())
		// kms lol
		arrivalTime, _ := time.Parse("January 02 3:04 PM 2006", value.Get("ArrivalLongDate").String()+" "+value.Get("Arrival").String()+" 2022")
		if arrivalTime.Unix() < 0 {
			//fmt.Println(value)
		}

		responses = append(responses, &shared.AggregatorResponse{
			Source:      "c-1",
			Carrier:     value.Get("CompanyName").String(),
			Type:        "Bus",
			Description: value.Get("CompanyName").String() + " Route",

			StartLocation: value.Get("OriginStationName").String(),
			EndLocation:   value.Get("DestinationStationName").String(),

			StartTime: int(departureTime.Unix()),
			EndTime:   int(arrivalTime.Unix()),

			Price:    int(value.Get("Price.Value").Float()),
			Duration: parseDuration(value.Get("Duration").String()),
		})

		return true
	})

	return responses, nil
}

func (r *checkmybusWorker) fetchHome() error {
	r.statusC("Fetching Session")

	req, err := http.NewRequest("GET", "https://www.checkmybus.com/search", nil)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	r.requestVerificationToken, _ = doc.Find("input[name='__RequestVerificationToken']").First().Attr("value")

	return nil
}
