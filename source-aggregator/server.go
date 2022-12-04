package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/cache"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/shared"
	"github.com/twinj/uuid"
	"log"
	"net/http"
	"strconv"
	"time"
)

// createServer creates a server that listens on ServerAddress.
// This server handles requests on /api/* endpoints.
func createServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/createQuery", createQueryHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:" + strconv.Itoa(ServerAddress),

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	funcs.W2("Server", "Started on port "+strconv.Itoa(ServerAddress))
	log.Fatal(srv.ListenAndServe())
}

// createQueryHandler handles requests on /api/createQuery endpoint.
func createQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		//w.Header().Add("Wen-Jia-Hu", "*")
		//w.Write([]byte("wen jia hu's evil kingdom"))

		return
	} else if r.Method != http.MethodPost {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var aggregatorQuery *shared.AggregatorQuery
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")

	err := json.NewDecoder(r.Body).Decode(&aggregatorQuery)
	if err != nil {
		funcs.R2("Server", "Error decoding aggregator query: "+err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	aggregatorQuery.Id = uuid.NewV4().String()
	aggregatorQuery.RespChan = make(chan []*shared.AggregatorResponse)
	if aggregatorQuery.StartLocation != nil && aggregatorQuery.StartLocation.Location.Latitude != 0.0 {
		aggregatorQuery.StartLocation.Latitude = aggregatorQuery.StartLocation.Location.Latitude
		aggregatorQuery.StartLocation.Longitude = aggregatorQuery.StartLocation.Location.Longitude
	}

	if aggregatorQuery.EndLocation != nil && aggregatorQuery.EndLocation.Location.Latitude != 0.0 {
		aggregatorQuery.EndLocation.Latitude = aggregatorQuery.EndLocation.Location.Latitude
		aggregatorQuery.EndLocation.Longitude = aggregatorQuery.EndLocation.Location.Longitude
	}

	funcs.M2("Server", fmt.Sprintf("Received Query Request (%v to %v) [%s]", aggregatorQuery.StartLocation, aggregatorQuery.EndLocation, aggregatorQuery.Id))

	body, ok := cache.Cacher.Get(hash(aggregatorQuery))

	if !ok {
		response := handleRequest(aggregatorQuery)
		body, _ = json.Marshal(&response)
		cache.Cacher.Set(hash(aggregatorQuery), body)
	} else {
		funcs.M2("Server", "Serve Cache Response")
	}

	w.Header().Add("content-type", "application/json")

	w.Write(body)
}

func hash(a *shared.AggregatorQuery) string {
	return fmt.Sprintf("%v|%v|%v|%v", a.StartLocation.Latitude, a.StartLocation.Longitude, a.EndLocation.Latitude, a.EndLocation.Longitude)
}
