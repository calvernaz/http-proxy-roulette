package roundtripper

import (
	"fmt"
	"net/http"
	"net/url"
)


var Result = make(map[int]int)

var _ http.RoundTripper = (*ProxyRoundTripper)(nil)
var _ ProxySelector = (*ProxyRoundTripper)(nil)

// ProxySelector is the interface that should be
// implemented by the algorithm behind the round tripper
type ProxySelector interface {
	Select() (*Proxy, error)
	WeightDown(*Proxy)
	WeightUp(*Proxy)
}

// The proxy round tripper
type ProxyRoundTripper struct {
	ProxySelector

	Tr      *http.Transport
}

// implements the http.RoundTripper interface
func (p *ProxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	proxy, err := p.Select()
	if err != nil {
		return nil, err
	}
	Result[proxy.Id]++

	p.Tr.Proxy = http.ProxyURL(&url.URL{
		Scheme: proxy.Scheme,
		User:   url.UserPassword(proxy.Username, proxy.Password),
		Host:   fmt.Sprintf("%s:%d", proxy.Host, proxy.Port),
	})
	response, err := p.Tr.RoundTrip(req)

	if err != nil {
		p.WeightDown(proxy)
		return response, err
	}

	if response.StatusCode != http.StatusOK {
		p.WeightDown(proxy)
	} else {
		p.WeightUp(proxy)
	}
	return response, err
}
