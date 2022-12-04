package data

import (
	wr "github.com/mroth/weightedrand"
	"math/rand"
	"time"
)

var acceptLanguageChooser *wr.Chooser

func init() {
	acceptLanguageChooser, _ = wr.NewChooser(
		wr.Choice{Item: "en-US,en;q=0.9", Weight: 100},
		wr.Choice{Item: "en-GB,en;q=0.9", Weight: 10},
		wr.Choice{Item: "ja,en-US;q=0.9,en;q=0.8", Weight: 5},
		wr.Choice{Item: "zh-CN,zh;q=0.9", Weight: 15},
	)
}

func fetchAcceptLanguage() string {
	rand.Seed(time.Now().UnixNano())

	return acceptLanguageChooser.Pick().(string)
}
