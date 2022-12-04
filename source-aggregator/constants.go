package main

import "time"

// TimeoutDuration is the timeout duration for the aggregator sources.
// If sources do not respond within this duration, the aggregator will return the responses it has received.
const TimeoutDuration = time.Duration(30) * time.Second

// ServerAddress is the local IP address of the server.
const ServerAddress = 80

// WorkerCount is the number of workers to spawn for each aggregator source.
const WorkerCount = 20
