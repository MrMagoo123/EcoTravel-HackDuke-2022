package shared

type AggregatorQuery struct {
	StartLocation *Location `json:"startLocation"`
	EndLocation   *Location `json:"endLocation"`

	RespChan chan []*AggregatorResponse

	Id string `json:"id"`
}

type Location struct {
	Formatted string  `json:"formatted"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`

	Location LegacyLocation `json:"location"`
}

type LegacyLocation struct {
	Formatted string  `json:"formatted"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type AggregatorWorker interface {
	Start() error
	Stop() error

	Query(query *AggregatorQuery) ([]*AggregatorResponse, error)
}

type AggregatorResponse struct {
	Price    int `json:"price"`
	Duration int `json:"duration"`

	Description string `json:"description"`
	Source      string `json:"source"`
	Carrier     string `json:"carrier"`
	Type        string `json:"type"`

	StartTime int `json:"startTime"`
	EndTime   int `json:"endTime"`

	StartLocation string `json:"startLocation"`
	EndLocation   string `json:"endLocation"`
}
