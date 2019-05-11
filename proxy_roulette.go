package roundtripper

import (
	"math/rand"
	"time"
)

type ProxyRoulette struct {
	Proxies   []Proxy
	Step      int
	MaxWeight int
}

func (pr ProxyRoulette) Select() (Proxy, error) {
	rand.Seed(time.Now().UnixNano())
	var r float64
	for {
		r = rand.Float64() * float64(len(pr.Proxies))
		rr := rand.Float64()
		rrr := float64(pr.Proxies[int(r)].Weight) / float64(pr.MaxWeight)

		if rr < rrr {
			increment := float64(pr.MaxWeight- pr.Proxies[int(r)].Weight) / float64(pr.Step)
			if pr.Proxies[int(r)].Weight + int(increment) > pr.MaxWeight {
				pr.Proxies[int(r)].Weight = pr.MaxWeight
			} else {
				pr.Proxies[int(r)].Weight += int(increment)
			}

			if pr.Proxies[int(r)].Weight > pr.MaxWeight {
				pr.MaxWeight = pr.Proxies[int(r)].Weight
			}
			break
		}
	}
	return pr.Proxies[int(r)], nil
}

