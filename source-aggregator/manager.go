package main

import (
	"fmt"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	"math/rand"
	"time"
)

// NewAggregatorWorkerManager creates an AggregatorWorkerManager.
func NewAggregatorWorkerManager(site string, spawnWorkersFunction func(string) (shared.AggregatorWorker, error)) *AggregatorWorkerManager {
	return &AggregatorWorkerManager{
		site:        site,
		query:       make(chan *shared.AggregatorQuery, WorkerCount),
		workers:     make(map[int]shared.AggregatorWorker),
		workerCount: make(map[int]bool),

		idealWorkerCount: WorkerCount,

		spawnWorkersFunction: spawnWorkersFunction,
	}
}

func (a *AggregatorWorkerManager) HandleQuery(query *shared.AggregatorQuery) {
	a.query <- query
}

func (a *AggregatorWorkerManager) Start() {

	for i := 0; i < a.idealWorkerCount; i++ {
		a.workers[i], _ = a.spawnWorkersFunction(fmt.Sprintf("Worker %v", i))
	}

	for i := 0; i < len(a.workers); i++ {
		go a.workers[i].Start()
	}

	go func() {
		for {
			select {
			case query := <-a.query:
				workerId := a.availableWorker()
				worker := a.workers[workerId]
				responseArr, err := worker.Query(query)
				a.workerCount[workerId] = true
				if err != nil {
					funcs.R("Error querying worker: " + err.Error())
				}

				query.RespChan <- responseArr
				a.workerCount[workerId] = false
			}
		}
	}()
}

func (a *AggregatorWorkerManager) availableWorker() int {
	// 循环赛
	for i, available := range a.workerCount {
		if available {
			return i
		}
	}

	rand.Seed(time.Now().UnixNano())

	return rand.Intn(len(a.workers))
}

// AggregatorWorkerManager manages a pool of AggregatorWorkers.
type AggregatorWorkerManager struct {
	site  string
	query chan *shared.AggregatorQuery

	workers     map[int]shared.AggregatorWorker
	workerCount map[int]bool

	idealWorkerCount     int
	spawnWorkersFunction func(string) (shared.AggregatorWorker, error)
}
