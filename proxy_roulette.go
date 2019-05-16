package roundtripper

import (
	"math/rand"
	"sync"
	"time"
)

// A proxy representation
type Proxy struct {
	Id int
	Username string
	Password string
	Scheme string
	Host string
	Port int
	Weight float32
}

type ProxyRoulette struct {
	once      sync.Once
	maxWeight float32

	Step    float32
	Proxies []*Proxy
}

// floor clipping
func (pr *ProxyRoulette) WeightDown(proxy *Proxy) {
	proxy.Weight = proxy.Weight - (proxy.Weight / (pr.maxWeight * pr.Step))
	if proxy.Weight < 0 {
		proxy.Weight = 0
	}
}

// ceil clipping
func (pr *ProxyRoulette) WeightUp(proxy *Proxy) {
	proxy.Weight = proxy.Weight + (proxy.Weight / (pr.maxWeight * pr.Step))
	if proxy.Weight > pr.maxWeight {
		proxy.Weight = pr.maxWeight
	}
}

// Roulette-wheel selection via stochastic acceptance
func (pr *ProxyRoulette) Select() (*Proxy, error) {
	pr.once.Do(func() {
		rand.Seed(time.Now().UnixNano())
		// linear normalization
		val := 1 / float32(len(pr.Proxies))
		for _, proxy := range pr.Proxies {
			proxy.Weight = val
		}

		pr.maxWeight = val
	})

	var index float32
	for {
		index = rand.Float32() * float32(len(pr.Proxies))
		rr := rand.Float32()
		rrr := pr.Proxies[int(index)].Weight / pr.maxWeight
		if rr < rrr {
			break
		}
	}

	return pr.Proxies[int(index)], nil
}
