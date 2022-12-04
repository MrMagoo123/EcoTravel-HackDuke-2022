package amtrak

import (
	"bytes"
	"encoding/json"
	"errors"
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

var (
	MaxRetries = errors.New("Max Retries Reached")
)

func header(input ...string) []string {
	var output []string
	for _, i := range input {
		output = append(output, i)
	}
	return output
}

func (r *amtrakWorker) fetchResults(query *shared.AggregatorQuery) ([]*shared.AggregatorResponse, error) {
	r.statusC(fmt.Sprintf("Fetching Results (%d)", r.akamaiSolveIncrement))
	offset := -1
	retries := 0

	destination, err := r.fetchNearbyStation(query.EndLocation.Formatted, 1)
	if err != nil {
		return nil, err
	}

noResults:
	if r.akamaiSolveIncrement >= 2 {
		err := r.solveAkamai()
		if err != nil {
			return nil, err
		}

		r.akamaiSolveIncrement = 0
	}

	offset++
	retries++

	if retries > 3 {
		return nil, MaxRetries
	}

	origin, err := r.fetchNearbyStation(query.StartLocation.Formatted, offset)
	if err != nil {
		return nil, err
	}

	journey := newJourney(origin, destination)
	payload, err := json.Marshal(&journey)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://www.amtrak.com/dotcom/journey-solution-option", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	r.akamaiSolveIncrement++

	req.Header = funcs.BuildHeaders([][]string{
		header("Host", "www.amtrak.com"),
		header("Connection", "keep-alive"),
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("Accept", "*/*"),
		header("Content-Type", "application/json"),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("User-Agent", r.browserConfig.UserAgent),
		header("Ocp-Apim-Subscription-Key", "B4C0755B-228A-48C9-892D-ED258055D763"),
		header("x-amtrak-trace-id", uuid.NewV4().String()),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("Origin", "https://www.amtrak.com"),
		header("Sec-Fetch-Site", "same-site"),
		header("Sec-Fetch-Mode", "cors"),
		header("Sec-Fetch-Dest", "empty"),
		header("Referer", "https://www.amtrak.com/"),
		header("Accept-Encoding", "gzip, deflate, br"),
		header("Accept-Language", r.browserConfig.AcceptLanguage),
	})

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 422 {
			goto noResults
		}

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
	dSplit := strings.Split(duration, "P0DT")
	if len(dSplit) < 2 {
		return -1
	}

	hSplit := strings.Split(dSplit[1], "H")
	if len(hSplit) < 2 {
		return -1
	}

	hours, _ := strconv.Atoi(hSplit[0])
	minutes, _ := strconv.Atoi(strings.Split(hSplit[1], "M")[0])

	return (hours*60 + minutes) * 60
}

func parseResults(body []byte) ([]*shared.AggregatorResponse, error) {
	var responses []*shared.AggregatorResponse

	jsonBody := gjson.ParseBytes(body)
	jsonBody.Get("data.journeySolutionOption.journeyLegs.0.journeyLegOptions").ForEach(func(key, value gjson.Result) bool {
		// Jan 2, 2006 at 3:04pm (MST)
		departureTime, _ := time.Parse("2006-01-02T15:04:00", value.Get("origin.schedule.departureDateTime").String())
		arrivalTime, _ := time.Parse("2006-01-02T15:04:00", value.Get("destination.schedule.arrivalDateTime").String())

		responses = append(responses, &shared.AggregatorResponse{
			Source:      "a-1",
			Description: "Amtrak Route (" + value.Get("origin.code").String() + " to " + value.Get("destination.code").String() + ")",

			Carrier: "Amtrak",
			Type:    "Train",

			Price:    int(value.Get("reservableAccommodations.0.accommodationFare.dollarsAmount.total").Float()),
			Duration: parseDuration(value.Get("elapsedTime").String()),

			StartLocation: value.Get("destination.name").String(),
			EndLocation:   value.Get("origin.name").String(),

			StartTime: int(departureTime.Unix()),
			EndTime:   int(arrivalTime.Unix()),
		})

		return true
	})

	return responses, nil
}

func (r *amtrakWorker) fetchNearbyStation(location string, offset int) (string, error) {
	req, err := http.NewRequest("GET", "https://www.amtrak.com/services/MapDataService/AutoCompleterArcgis/getResponseList?searchTerm="+url.QueryEscape(location), nil)
	if err != nil {
		return "", err
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("Host", "www.amtrak.com"),
		header("Connection", "keep-alive"),
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("Accept", "*/*"),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("User-Agent", r.browserConfig.UserAgent),
		header("Ocp-Apim-Subscription-Key", "B4C0755B-228A-48C9-892D-ED258055D763"),
		header("x-amtrak-trace-id", uuid.NewV4().String()),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("Origin", "https://www.amtrak.com"),
		header("Sec-Fetch-Site", "same-site"),
		header("Sec-Fetch-Mode", "cors"),
		header("Sec-Fetch-Dest", "empty"),
		header("Referer", "https://www.amtrak.com/"),
		header("Accept-Encoding", "gzip, deflate, br"),
		header("Accept-Language", r.browserConfig.AcceptLanguage),
	})

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error Fetching Results (%v)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// lol fuck
	// amtrak api returns random shit trainstations that are not even remotely operational
	// if you search for new york, new york state fair station comes before ny penn
	// ?????
	hardcodedTrainStations := []string{
		"BOS",
		"BBY",
		"NYP",
	}

	bodyString := string(body)
	for _, station := range hardcodedTrainStations {
		if strings.Contains(bodyString, station) {
			return station, nil
		}
	}

	bodyJson := gjson.ParseBytes(body)
	offset = offset % len(bodyJson.Get("autoCompleterResponse.autoCompleteList").Array())
	r.statusC(fmt.Sprintf("Fetching Nearby Station (%s) (Offset: %d)", location, offset))

	return bodyJson.Get("autoCompleterResponse.autoCompleteList." + strconv.Itoa(offset) + ".stationCode").String(), nil
}

func (r *amtrakWorker) fetchHome() error {
	r.statusC("Fetching Session")

	req, err := http.NewRequest("GET", "https://www.amtrak.com/home.html", nil)
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
		//header("accept-encoding", "gzip, deflate, br"), // -- fuck
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

	doc.Find("script[type=\"text/javascript\"]").Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("type")
		src, _ := s.Attr("src")

		// reject external scripts
		if strings.Contains(src, ".com") && !strings.Contains(src, "www.amtrak.com") {
			return
		}

		// if the script URL contains the baseurl, trim it out
		src = strings.TrimLeft(src, "www.amtrak.com")
		if len(src) < 2 || strings.Count(src, "/") < 6 || strings.Contains(src, "etc") {
			return
		}

		if val == "text/javascript" {
			if src[0:1] != "/" {
				src = "/" + src
			}

			r.akamaiURL = "https://www.amtrak.com" + src
		}
	})

	return nil
}
