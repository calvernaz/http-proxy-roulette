# Proxy Roulette

This is a toy example of an implementation of the **[roulette-wheel selection via stochastic acceptance](https://en.wikipedia.org/wiki/Stochastic_universal_sampling)** in Golang and [http.RoundTripper](https://golang.org/pkg/net/http/#RoundTripper) for proxy selection.
In other words it selects a random proxy from a set in _O(1)_ time but also adds a rank function for increase/decrease probability depending on
the proxy state.

## Example

```$go
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
``` 
