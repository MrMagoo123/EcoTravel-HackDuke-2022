package main

import (
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/sources/amtrak"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/sources/checkmybus"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/sources/greyhound"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/sources/megabus"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/sources/rome2rio"
)

var pools []*AggregatorWorkerManager

func main() {
	initialize()

	pools = []*AggregatorWorkerManager{
		NewAggregatorWorkerManager("rome2rio", rome2rio.CreateWorker),
		NewAggregatorWorkerManager("megabus", megabus.CreateWorker),
		NewAggregatorWorkerManager("greyhound", greyhound.CreateWorker),
		NewAggregatorWorkerManager("checkmybus", checkmybus.CreateWorker),
		NewAggregatorWorkerManager("amtrak", amtrak.CreateWorker),
	}

	for _, manager := range pools {
		manager.Start()

		// Synchronize the worker start
		//time.Sleep(100 * time.Millisecond)
	}

	createServer()
}
