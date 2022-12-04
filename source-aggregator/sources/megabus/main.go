package megabus

import (
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/data"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	http "github.com/saucesteals/fhttp"
)

func CreateWorker(id string) (shared.AggregatorWorker, error) {
	return &megabusWorker{
		id:            "Megabus " + id,
		browserConfig: data.NewConfig(),
	}, nil
}

type megabusWorker struct {
	id            string
	client        *http.Client
	browserConfig *data.BrowserConfig
	nodes         *MegabusNodes
}

func (r *megabusWorker) statusC(status string) {
	funcs.C2(r.id, status)
}

func (r *megabusWorker) Start() (err error) {
	funcs.C2(r.id, "Started")

	r.client, err = funcs.NewClient("")
	if err != nil {
		return
	}

	r.fetchHome()
	r.fetchNodes()

	return
}

func (r *megabusWorker) Stop() error {
	return nil
}

func (r *megabusWorker) Query(query *shared.AggregatorQuery) ([]*shared.AggregatorResponse, error) {
	funcs.C2(r.id, "Processing Query "+query.Id)

	aggregatorResponse, err := r.fetchResults(query)
	if err != nil {
		return nil, err
	}

	funcs.G2(r.id, "Processed Query "+query.Id)

	return aggregatorResponse, nil
}
