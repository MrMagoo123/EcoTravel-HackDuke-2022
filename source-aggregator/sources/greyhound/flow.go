package greyhound

import (
	"fmt"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	http "github.com/saucesteals/fhttp"
	"github.com/tidwall/gjson"
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

func (r *greyhoundWorker) fetchNearbyLocation(locationName string) (string, string, error) {
	r.statusC(fmt.Sprintf("Fetching Nearby Location (%v)", locationName))

	req, err := http.NewRequest("GET", fmt.Sprintf("https://eapi.greyhound.com/commercial-locations/locationsPartialName?partialName=%s", url.QueryEscape(locationName)), nil)
	if err != nil {
		return "", "", err
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("Host", "eapi.greyhound.com"),
		header("Connection", "keep-alive"),
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("Accept", "*/*"),
		header("Content-Type", "application/x-www-form-urlencoded"),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("User-Agent", r.browserConfig.UserAgent),
		header("Ocp-Apim-Subscription-Key", "B4C0755B-228A-48C9-892D-ED258055D763"),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("Origin", "https://www.greyhound.com"),
		header("Sec-Fetch-Site", "same-site"),
		header("Sec-Fetch-Mode", "cors"),
		header("Sec-Fetch-Dest", "empty"),
		header("Referer", "https://www.greyhound.com/"),
		header("Accept-Encoding", "gzip, deflate, br"),
		header("Accept-Language", r.browserConfig.AcceptLanguage),
	})

	resp, err := r.client.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("Error Fetching Nearby Location (%v)", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	jsonBody := gjson.ParseBytes(body)
	stationLocation := jsonBody.Get("Data.0.StationLocationUrl").String()
	stationLocationSplit := strings.Split(stationLocation, "bus-station-")
	if len(stationLocationSplit) < 2 {
		return "", "", fmt.Errorf("Error Fetching Nearby Location (%v)", "Invalid Station Location")
	}

	return jsonBody.Get("Data.0.City").String(), strings.Split(stationLocationSplit[1], "#")[0], nil
}

func (r *greyhoundWorker) fetchResults(query *shared.AggregatorQuery) ([]*shared.AggregatorResponse, error) {
	r.statusC("Fetching Results")

	destCity, destCode, _ := r.fetchNearbyLocation(query.StartLocation.Formatted)
	originCity, originCode, _ := r.fetchNearbyLocation(query.EndLocation.Formatted)

	req, err := http.NewRequest("GET", "https://eapi.greyhound.com/commercial/schedules/schedule?originCity="+url.QueryEscape(originCity)+"&originCode="+originCode+"&destinationCity="+url.QueryEscape(destCity)+"&destinationCode="+destCode+"&departOn=2022-11-21&returnOn=&totalPassengers=1+Passengers&adults=1&adultWheelchairs=0&children=0&childWheelchairs=0&seniors=0&seniorWheelchairs=0&discountCode=&version=4&GetTotalPassengers=1&HasWheelchairTraveler=false&HasReturnFare=false&webOnly=false&originStateAbbreviation=MA&originStateCode=4&destinationStateAbbreviation=NY&destinationStateCode=15&calStart=2022-11-18&calEnd=2022-11-24", nil)
	if err != nil {
		return nil, err
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("Host", "eapi.greyhound.com"),
		header("Connection", "keep-alive"),
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("Accept", "*/*"),
		header("Content-Type", "application/x-www-form-urlencoded"),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("User-Agent", r.browserConfig.UserAgent),
		header("Ocp-Apim-Subscription-Key", "B4C0755B-228A-48C9-892D-ED258055D763"),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("Origin", "https://www.greyhound.com"),
		header("Sec-Fetch-Site", "same-site"),
		header("Sec-Fetch-Mode", "cors"),
		header("Sec-Fetch-Dest", "empty"),
		header("Referer", "https://www.greyhound.com/"),
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
	hSplit := strings.Split(duration, ":")
	if len(hSplit) < 2 {
		return -1
	}

	hours, _ := strconv.Atoi(hSplit[0])
	minutes, _ := strconv.Atoi(hSplit[1])

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

	jsonBody.Get("Data.Schedules").ForEach(func(key, value gjson.Result) bool {
		// Jan 2, 2006 at 3:04pm (MST)
		departureTime, _ := time.Parse("2006-01-02T15:04:00", value.Get("DepartureDateTime").String())
		arrivalTime, _ := time.Parse("2006-01-02T15:04:00", value.Get("ArrivalDateTime").String())

		responses = append(responses, &shared.AggregatorResponse{
			Source:      "g-1",
			Carrier:     "Greyhound",
			Type:        "Bus",
			Description: "Greyhound Route " + value.Get("ScheduleNumber").String(),

			StartLocation: stations[int(value.Get("DepartureLocationCode").Int())],
			EndLocation:   stations[int(value.Get("ArrivalLocationCode").Int())],

			StartTime: int(departureTime.Unix()),
			EndTime:   int(arrivalTime.Unix()),

			Price:    int(value.Get("OnlineFares.LowestFare").Int()),
			Duration: parseDuration(value.Get("Duration").String()),
		})

		return true
	})

	return responses, nil
}

func (r *greyhoundWorker) fetchHome() error {
	r.statusC("Fetching Session")

	req, err := http.NewRequest("GET", "https://www.greyhound.com/", nil)
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
