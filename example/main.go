package main

import (
	"io/ioutil"
	"net/http"
	"time"

	roundtripper "github.com/calvernaz/http-proxy-roulette"
)

var proxies = []roundtripper.Proxy{
	{
		Id:       0,
		Username: "justwatch",
		Password: "ohthighaimeiY9",
		Scheme:   "http",
		Host:     "38.64.52.182",
		Port:     60000,
		Weight:   3,
	},
	{
		Id:       1,
		Username: "justwatch",
		Password: "ohthighaimeiY9",
		Scheme:   "http",
		Host:     "154.49.200.7",
		Port:     60000,
		Weight:   10,
	},
}

func main() {
	roulette := roundtripper.ProxyRoulette{
		Proxies:   proxies,
		Step:      5,
		MaxWeight: 13,
	}

	rt := &roundtripper.ProxyRoundTripper{
		ProxySelector: roulette,
	}

	var client = &http.Client{
		Transport: rt,
		Timeout:   30 * time.Second,
	}

	for i := 0; i < 50; i++ {

		req, _ := http.NewRequest("GET", "https://google.com", nil)
		resp, _ := client.Do(req)

		if _, err := ioutil.ReadAll(resp.Body); err != nil {
			//fmt.Println(string(b))
			panic(err)
		}
	}
}
