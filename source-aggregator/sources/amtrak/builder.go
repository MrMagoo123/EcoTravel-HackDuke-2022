package amtrak

func newJourney(origin, destination string) *Journey {
	return &Journey{
		InitialJourneyLegOnly: false,
		JourneyRequest: JourneyRequest{
			Fare:               Fare{PricingUnit: "DOLLARS"},
			AlternateDayOption: true,
			Customer: Customer{
				TierStatus: "MEMBER",
			},
			Type: "OW",
			JourneyLegRequests: []JourneyLegRequests{JourneyLegRequests{
				Origin: Origin{
					Code: origin,
					Schedule: Schedule{
						DepartureDateTime: "2022-11-21T00:00:00",
					},
				},

				Destination: Destination{
					Code: destination,
				},

				Passengers: []Passengers{
					Passengers{
						ID:          "P1",
						Type:        "F",
						InitialType: "adult",
					},
				},
			}},
		},
		ReservableAccomodationOptions: "ALL",
	}
}

type Journey struct {
	JourneyRequest                JourneyRequest `json:"journeyRequest"`
	InitialJourneyLegOnly         bool           `json:"initialJourneyLegOnly"`
	ReservableAccomodationOptions string         `json:"reservableAccomodationOptions"`
}
type Fare struct {
	PricingUnit string `json:"pricingUnit"`
}
type Schedule struct {
	DepartureDateTime string `json:"departureDateTime"`
}
type Origin struct {
	Code     string   `json:"code"`
	Schedule Schedule `json:"schedule"`
}
type Destination struct {
	Code string `json:"code"`
}
type Passengers struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	InitialType string `json:"initialType"`
}
type JourneyLegRequests struct {
	Origin      Origin       `json:"origin"`
	Destination Destination  `json:"destination"`
	Passengers  []Passengers `json:"passengers"`
}
type Customer struct {
	TierStatus string `json:"tierStatus"`
}
type JourneyRequest struct {
	Fare               Fare                 `json:"fare"`
	AlternateDayOption bool                 `json:"alternateDayOption"`
	Type               string               `json:"type"`
	JourneyLegRequests []JourneyLegRequests `json:"journeyLegRequests"`
	Customer           Customer             `json:"customer"`
	IsPassRider        bool                 `json:"isPassRider"`
}
