package main

import (
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	"sync"
	"time"
)

func handleRequest(query *shared.AggregatorQuery) []*shared.AggregatorResponse {
	wg := &sync.WaitGroup{}
	doneChan := make(chan bool)

	for _, manager := range pools {
		go manager.HandleQuery(query)
	}

	var responses []*shared.AggregatorResponse

	for _, _ = range pools {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			r := <-query.RespChan
			responses = append(responses, r...)
			wg.Done()
		}(wg)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		doneChan <- true
	}(wg)

	// implement 30s timeout
	select {
	case <-time.After(TimeoutDuration):
		return responses
	case <-doneChan:
		return responses
	}

	return responses
}
