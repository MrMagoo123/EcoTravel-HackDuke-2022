package skyscanner

type Query struct {
	Market                  string        `json:"market"`
	Locale                  string        `json:"locale"`
	Currency                string        `json:"currency"`
	AlternativeOrigins      bool          `json:"alternativeOrigins"`
	AlternativeDestinations bool          `json:"alternativeDestinations"`
	Adults                  int           `json:"adults"`
	CabinClass              string        `json:"cabin_class"`
	ChildAges               []interface{} `json:"child_ages"`
	InboundDate             string        `json:"inboundDate"`
	Legs                    []Legs        `json:"legs"`
}
type Legs struct {
	Origin                     string `json:"origin"`
	Destination                string `json:"destination"`
	Date                       string `json:"date"`
	ReturnDate                 string `json:"return_date"`
	AddAlternativeOrigins      bool   `json:"add_alternative_origins"`
	AddAlternativeDestinations bool   `json:"add_alternative_destinations"`
}
