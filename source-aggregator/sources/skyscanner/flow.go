package skyscanner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	http "github.com/saucesteals/fhttp"
	"github.com/tidwall/gjson"
	"github.com/twinj/uuid"
	"io"
)

var (
	ErrNotCompleted = errors.New("not completed, go away!")
)

func header(input ...string) []string {
	var output []string
	for _, i := range input {
		output = append(output, i)
	}
	return output
}

func (r *skyscannerWorker) createSlug(query *shared.AggregatorQuery) (string, error) {
	r.statusC("Fetching Results")

	apiQuery := &Query{
		Market:                  "US",
		Locale:                  "en-US",
		Currency:                "USD",
		AlternativeOrigins:      false,
		AlternativeDestinations: false,
		Adults:                  1,
		CabinClass:              "economy",
		ChildAges:               nil,
		InboundDate:             "2022-11-03",
		Legs: []Legs{
			Legs{
				Origin:                     "JFK",
				Destination:                "BOS",
				Date:                       "2022-10-27",
				AddAlternativeDestinations: false,
				AddAlternativeOrigins:      false,
			},
		},
	}

	payload, err := json.Marshal(&apiQuery)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://www.skyscanner.com/g/conductor/v1/fps3/search/?geo_schema=skyscanner&carrier_schema=skyscanner&response_include=query%3Bdeeplink%3Bsegment%3Bstats%3Bfqs%3Bpqs", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("content-length", ""),
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("x-skyscanner-devicedetection-istablet", "false"),
		header("x-skyscanner-channelid", "website"),
		header("x-skyscanner-devicedetection-istablet", "false"),
		header("x-skyscanner-utid", uuid.NewV4().String()),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("x-skyscanner-traveller-context", uuid.NewV4().String()),
		header("content-type", "application/json"),
		header("accept", "application/json"),
		header("x-skyscanner-viewid", uuid.NewV4().String()),
		header("user-agent", r.browserConfig.UserAgent),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("origin", "https://www.skyscanner.com"),
		header("sec-fetch-site", "same-site"),
		header("sec-fetch-mode", "cors"),
		header("sec-fetch-dest", "empty"),
		header("referer", "https://www.skyscanner.com/transport/flights"),
		header("accept-encoding", "gzip, deflate, br"),
		header("accept-language", r.browserConfig.AcceptLanguage),
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

	return gjson.ParseBytes(body).Get("context.session_id").String(), nil
}

func parseResults(body []byte) ([]*shared.AggregatorResponse, error) {
	var responses []*shared.AggregatorResponse

	jsonBody := gjson.ParseBytes(body)
	jsonBody.Get("itineraries").ForEach(func(key, value gjson.Result) bool {
		pricingOptions := value.Get("pricing_options.0")

		responses = append(responses, &shared.AggregatorResponse{
			Source:      "skyscanner",
			Description: value.Get("name").String(),

			Price:    int(value.Get("indicativePrices.0.price").Int()),
			Duration: int(value.Get("duration").Int()),
		})

		return true
	})

	return responses, nil
}

func (r *skyscannerWorker) fetchHome() error {
	r.statusC("Fetching Session")

	req, err := http.NewRequest("GET", "https://www.skyscanner.com/", nil)
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

func (r *skyscannerWorker) fetchResults(slug string) ([]*shared.AggregatorResponse, error) {
	r.statusC(fmt.Sprintf("Fetching Results (%v)", slug))

	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.skyscanner.com/g/conductor/v1/fps3/search/%s?geo_schema=skyscanner&carrier_schema=skyscanner&response_include=query%3Bdeeplink%3Bsegment%3Bstats%3Bfqs%3Bpqs", slug), nil)
	if err != nil {
		return nil, nil
	}

	req.Header = funcs.BuildHeaders([][]string{
		header("sec-ch-ua", r.browserConfig.SecChUa),
		header("x-skyscanner-devicedetection-istablet", "false"),
		header("x-skyscanner-channelid", "website"),
		header("x-skyscanner-devicedetection-istablet", "false"),
		header("x-skyscanner-utid", uuid.NewV4().String()),
		header("sec-ch-ua-mobile", r.browserConfig.SecChUaMobile),
		header("x-skyscanner-traveller-context", uuid.NewV4().String()),
		header("content-type", "application/json"),
		header("accept", "application/json"),
		header("x-skyscanner-viewid", uuid.NewV4().String()),
		header("user-agent", r.browserConfig.UserAgent),
		header("sec-ch-ua-platform", r.browserConfig.SecChUaPlatform),
		header("origin", "https://www.skyscanner.com"),
		header("sec-fetch-site", "same-site"),
		header("sec-fetch-mode", "cors"),
		header("sec-fetch-dest", "empty"),
		header("referer", "https://www.skyscanner.com/transport/flights"),
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
