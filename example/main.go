package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	roundtripper "github.com/calvernaz/http-proxy-roulette"
)

var proxies = []*roundtripper.Proxy{
	{
		Id:       0,
		Scheme:   "http",
		Host:     "167.86.79.154",
		Port:     3128,
		Weight:   0.25,
	},
	{
		Id:       1,
		Scheme:   "http",
		Host:     "35.199.45.8",
		Port:     3128,
		Weight:   0.25,
	},
	{
		Id:       2,
		Scheme:   "http",
		Host:     "157.230.33.168",
		Port:     8080,
		Weight:   0.25,
	},
	{
		Id:       3,
		Scheme:   "http",
		Host:     "205.235.57.12",
		Port:     3128,
		Weight:   0.25,
	},
}

func main() {
	roulette := &roundtripper.ProxyRoulette{
		Proxies: proxies,
		Step:    15,
	}

	rt := &roundtripper.ProxyRoundTripper{
		ProxySelector: roulette,
		Tr: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: 1,
		},
	}

	var client = &http.Client{
		Transport: rt,
		Timeout:   30 * time.Second,
	}

	for i := 0; i < 50; i++ {

		req, _ := http.NewRequest("GET", "https://lwn.net", nil)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if _, err := ioutil.ReadAll(resp.Body); err != nil {
			panic(err)
		}
	}

	fmt.Printf("%+v \n", roundtripper.Result)
}
