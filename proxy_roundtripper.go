package roundtripper

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

var _ http.RoundTripper = (*ProxyRoundTripper)(nil)
var _ ProxySelector = (*ProxyRoundTripper)(nil)

type Proxy struct {
	Id int
	Username string
	Password string
	Scheme string
	Host string
	Port int
	Weight int
}

type ProxySelector interface {
	// TODO select could be func(*Request) (*url.URL, error)
	Select() (Proxy, error)
}

type ProxyRoundTripper struct {
	ProxySelector
}

// implements the http.RoundTripper interface
func (p *ProxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	proxy, err := p.Select()
	if err != nil {
		return nil, err
	}

	rt := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 1,
		Proxy: http.ProxyURL(&url.URL{
			Scheme: proxy.Scheme,
			User:   url.UserPassword(proxy.Username, proxy.Password),
			Host:   fmt.Sprintf("%s:%d", proxy.Host, proxy.Port),
		}),
	}

	response, err := rt.RoundTrip(req)
	if err != nil {
		return response, err
	}

	return response, err
}
