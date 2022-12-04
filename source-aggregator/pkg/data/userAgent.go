package data

import (
	wr "github.com/mroth/weightedrand"
	"math/rand"
	"time"
)

var uaChooser *wr.Chooser

func init() {
	uaChooser, _ = wr.NewChooser(
		wr.Choice{Item: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36", Weight: 30},
		wr.Choice{Item: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36", Weight: 30},
		wr.Choice{Item: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36", Weight: 20},
		wr.Choice{Item: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36", Weight: 12},
	)
}

func fetchUserAgent() string {
	rand.Seed(time.Now().UnixNano())

	return uaChooser.Pick().(string)
}
