package rome2rio

import (
	"fmt"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	http "github.com/saucesteals/fhttp"
	"github.com/tidwall/gjson"
	"io"
	"net/url"
	"strconv"
	"sync"
	"time"
)

func header(input ...string) []string {
	var output []string
	for _, i := range input {
		output = append(output, i)
	}
	return output
}

func (r *rome2rioWorker) fetchResults(query *shared.AggregatorQuery) ([]*shared.AggregatorResponse, error) {
	r.statusC("Fetching Results")

	req, err := http.NewRequest("GET", "https://services.rome2rio.com/api/1.5/json/search?key=jGq3Luw3&oName="+url.QueryEscape(query.StartLocation.Formatted)+"&dName="+url.QueryEscape(query.EndLocation.Formatted)+"&languageCode=en&currencyCode=USD&uid=USRal20221021002634250ufdd&aqid=USRal20221022191356694ufdd&analytics=true&debugFeatures=&debugExperiments=&groupOperators=true", nil)
	if err != nil {
		return nil, err
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("user-agent", r.browserConfig.UserAgent),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("accept", "*/*"),
		header("origin", "https://www.rome2rio.com"),
		header("sec-fetch-site", "same-site"),
		header("sec-fetch-mode", "cors"),
		header("sec-fetch-dest", "empty"),
		header("referer", "https://www.rome2rio.com/map/Boston/New-York"),
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

	return r.parseResults(body)
}

func (r *rome2rioWorker) parseResults(body []byte) ([]*shared.AggregatorResponse, error) {
	var responses []*shared.AggregatorResponse

	wg := &sync.WaitGroup{}

	jsonBody := gjson.ParseBytes(body)
	jsonBody.Get("routes").ForEach(func(key, value gjson.Result) bool {
		scheduleInfo := value.Get("scheduleInfo.data").String()
		if scheduleInfo == "" {
			return true
		}

		wg.Add(1)

		go func(scheduleInfo string) {
			defer wg.Done()

			carrier, start, end, departureTime, arrivalTime, err := r.fetchSchedule(scheduleInfo)
			if err != nil {
				return
			}

			responses = append(responses, &shared.AggregatorResponse{
				Source:      "f-1",
				Description: value.Get("name").String(),

				Price:    int(value.Get("indicativePrices.0.price").Int()),
				Duration: int(value.Get("duration").Int()),

				Carrier:       carrier,
				StartLocation: start,
				EndLocation:   end,

				StartTime: departureTime,
				EndTime:   arrivalTime,
			})
		}(scheduleInfo)

		return true
	})

	wg.Wait()

	return responses, nil
}

func (r *rome2rioWorker) fetchSchedule(slug string) (carrier string, start string, end string, departureTime int, arrivalTime int, err error) {
	r.statusC("Fetching Route Schedule")
	req, err := http.NewRequest("GET", "https://services.rome2rio.com/api/1.5/json/schedules", nil)
	if err != nil {
		return
	}

	query := req.URL.Query()

	query.Add("key", "jGq3Luw3")
	query.Add("searchData", slug)
	query.Add("languageCode", "en")
	query.Add("oDateTime", "2022-10-23T02:00")
	query.Add("useIndicativePrices", "true")

	req.URL.RawQuery = query.Encode()

	req.Header = funcs.BuildHeaders([][]string{
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("user-agent", r.browserConfig.UserAgent),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("accept", "*/*"),
		header("origin", "https://www.rome2rio.com"),
		header("sec-fetch-site", "same-site"),
		header("sec-fetch-mode", "cors"),
		header("sec-fetch-dest", "empty"),
		header("referer", "https://www.rome2rio.com/map/Boston/New-York"),
		header("accept-encoding", "gzip, deflate, br"),
		header("accept-language", r.browserConfig.AcceptLanguage),
	})

	resp, err := r.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Error Fetching Schedule (%v)", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	bodyJson := gjson.ParseBytes(body)

	carrier = bodyJson.Get("carriers.0.name").String()

	start = bodyJson.Get("hopTicketLinks.0.originCanonical").String()
	end = bodyJson.Get("hopTicketLinks." + strconv.Itoa(len(bodyJson.Get("hopTicketLinks").Array())-1) + ".destinationCanonical").String()
	if start == "" || end == "" {
		start = bodyJson.Get("places.0.name").String()
		end = bodyJson.Get("places." + strconv.Itoa(len(bodyJson.Get("places").Array())-1) + ".name").String()
	}

	// Jan 2, 2006 at 3:04pm (MST)
	departureT, _ := time.Parse("2006-01-02 15:04", bodyJson.Get("hops.0.departureDate").String()+" "+bodyJson.Get("hops.0.departureTime").String())
	departureTime = int(departureT.Unix())

	arrivalT, _ := time.Parse("2006-01-02 15:04", bodyJson.Get("hops."+strconv.Itoa(len(bodyJson.Get("hops").Array())-1)+".arrivalDate").String()+" "+bodyJson.Get("hops."+strconv.Itoa(len(bodyJson.Get("hops").Array())-1)+".arrivalTime").String())
	arrivalTime = int(arrivalT.Unix())

	if departureTime < 0 {
		fmt.Println(bodyJson.Get("hops.0.departureDate").String() + " " + bodyJson.Get("hops.0.departureTime").String())
	}
	return
}

func (r *rome2rioWorker) fetchHome() error {
	r.statusC("Fetching Session")

	req, err := http.NewRequest("GET", "https://www.rome2rio.com/", nil)
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
