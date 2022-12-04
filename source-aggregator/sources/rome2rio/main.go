package rome2rio

import (
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/data"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	http "github.com/saucesteals/fhttp"
)

func CreateWorker(id string) (shared.AggregatorWorker, error) {
	return &rome2rioWorker{
		id:            "Rome2Rio " + id,
		browserConfig: data.NewConfig(),
	}, nil
}

type rome2rioWorker struct {
	id            string
	client        *http.Client
	browserConfig *data.BrowserConfig
}

func (r *rome2rioWorker) statusC(status string) {
	funcs.C2(r.id, status)
}

func (r *rome2rioWorker) Start() (err error) {
	funcs.C2(r.id, "Started")

	r.client, err = funcs.NewClient("")
	if err != nil {
		return
	}

	r.fetchHome()

	return
}

func (r *rome2rioWorker) Stop() error {
	return nil
}

func (r *rome2rioWorker) Query(query *shared.AggregatorQuery) ([]*shared.AggregatorResponse, error) {
	funcs.C2(r.id, "Processing Query "+query.Id)

	aggregatorResponse, err := r.fetchResults(query)
	if err != nil {
		return nil, err
	}

	funcs.G2(r.id, "Processed Query "+query.Id)

	return aggregatorResponse, nil
}
