package amtrak

import (
	"github.com/maxwellmlin/duke-hackathon-2022/source-aggregator/pkg/funcs"
	http "github.com/saucesteals/fhttp"
	"net/url"
	"strings"
)

func (r *amtrakWorker) solveAkamai() error {
	// 依赖不同的事件，生成不同的策略数据
	// - gt

	for i := 0; i < SensorDataRetries; i++ {
		sd, err := r.fetchAkamaiSd()
		if err != nil {
			return err
		}

		err = r.postSd(sd)
		if err != nil {
			return err
		}
	}

	if strings.Contains(grabCookie(r.client, "https://www.amtrak.com/", "_abck"), "~0~") {
		funcs.G("Generated Valid Akamai Cookie! (Contains ~0~)")
	} else {
		r.statusC("Fail generate cookie!")
	}

	return nil
}

func (r *amtrakWorker) postSd(sd []byte) error {
	// --- private API keys

	return nil
}

func (r *amtrakWorker) fetchAkamaiSd() ([]byte, error) {
	// --- private API keys

	return nil, nil
}

func decodeSensorData(input []byte) []byte {
	// --- private API keys

	return input
}

func grabCookie(client *http.Client, targetUrl string, targetCookie string) string {
	s, _ := url.Parse(targetUrl)
	cookies := client.Jar.Cookies(s)
	for _, c := range cookies {
		if c.Name == targetCookie {
			return c.Value
		}
	}
	return ""
}
