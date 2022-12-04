package funcs

import (
	"github.com/saucesteals/fhttp"
	"github.com/saucesteals/fhttp/cookiejar"
	"net/url"
	"strings"
	"time"
)

const DefineOrder = "DefineOrder"

func NewClient(proxyURL string) (*http.Client, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	jar, _ := cookiejar.New(nil)

	transport := &http.Transport{
		Proxy: http.ProxyURL(u),
	}

	return &http.Client{
		Transport: transport,
		Jar:       jar,
		Timeout:   10 * time.Second,
	}, nil
}

func Header(input ...string) []string {
	var output []string
	for _, i := range input {
		output = append(output, i)
	}
	return output
}

func BuildHeaders(headers [][]string) map[string][]string {
	var headerOrder []string
	headerMap := make(map[string][]string)

	for _, header := range headers {
		if len(header) > 1 {
			var name string
			name, header = header[0], header[1:] // pop header name

			for _, val := range header { // build header values
				if val != DefineOrder && strings.ToLower(name) != "content-length" { // check to see if we are just defining order, for things such as cookie
					headerMap[name] = append(headerMap[name], val)
				}
			}
			// add to header order

			if strings.ToLower(name) == "content-length" {
				headerOrder = append(headerOrder, name)
			}
			headerOrder = append(headerOrder, name)
		}
	}

	headerMap["Header-Order:"] = headerOrder
	headerMap["PHeader-Order:"] = []string{
		":method",
		":authority",
		":scheme",
		":path",
	}
	return headerMap
}
