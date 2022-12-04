package skyscanner

import (
	"fmt"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/data"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	http "github.com/saucesteals/fhttp"
)

func CreateWorker(id string) (shared.AggregatorWorker, error) {
	return &skyscannerWorker{
		id:            "skyscanner " + id,
		browserConfig: data.NewConfig(),
	}, nil
}

type skyscannerWorker struct {
	id            string
	client        *http.Client
	browserConfig *data.BrowserConfig
}

func (r *skyscannerWorker) statusC(status string) {
	funcs.C2(r.id, status)
}

func (r *skyscannerWorker) Start() (err error) {
	funcs.C2(r.id, "Started")

	// 是故百战百胜，非善之善也﹔不战而屈人之兵，善之善者也。
	r.browserConfig.UserAgent = fmt.Sprintf("%s (Compatible %s)", r.browserConfig.UserAgent, "Qwantify/1.0")

	r.client, err = funcs.NewClient("")
	if err != nil {
		return
	}

	r.fetchHome()

	return
}

func (r *skyscannerWorker) Stop() error {
	return nil
}

func (r *skyscannerWorker) Query(query *shared.AggregatorQuery) ([]*shared.AggregatorResponse, error) {
	funcs.C2(r.id, "Processing Query "+query.Id)

	aggregatorResponse, err := r.createSlug(query)
	if err != nil {
		return nil, err
	}

	funcs.G2(r.id, "Processed Query "+query.Id)
	for {
		aggregatorResponse, err = r.fetchResults(aggregatorResponse)
		if err != nil {
			return nil, err
		}

		if aggregatorResponse.Status == "Completed" {
			break
		}
	}

	return aggregatorResponse, nil
}
