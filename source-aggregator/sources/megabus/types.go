package megabus

import (
	"math"
)

type MegabusNodes struct {
	Cities []City `json:"cities"`
}

func (r *megabusWorker) findClosestNode(lat, long float64) int {
	var closestNode *City
	var closestDistance float64

	for _, node := range r.nodes.Cities {
		node := node // holyshit what the fuck golang moment
		distance := node.distance(lat, long)
		if closestNode == nil || distance < closestDistance {
			closestNode = &node
			closestDistance = distance
		}
	}

	// if this panics you are fucked, sorry
	return closestNode.Id
}

type City struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Id int `json:"id"`
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

func hsin(t float64) float64 {
	return math.Pow(math.Sin(t/2), 2)
}

func (c *City) distance(lat, long float64) float64 {
	lat1 := degreesToRadians(lat)
	lon1 := degreesToRadians(long)
	lat2 := degreesToRadians(c.Latitude)
	lon2 := degreesToRadians(c.Longitude)
	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := hsin(diffLat) + math.Cos(lat1)*math.Cos(lat2)*hsin(diffLon)
	c2 := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := c2 * 3958 // earth radius in miles

	return math.Round(distance*100) / 100
}
