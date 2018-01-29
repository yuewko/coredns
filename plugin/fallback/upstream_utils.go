// Package fallback implements a fallback plugin for CoreDNS
package fallback

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/proxy"
)

// upstreamMapper manages map of rcode to upstream
type upstreamMapper interface {
	Add(int, proxy.Upstream) proxy.Upstream
	Get(int) (proxy.Upstream, bool)
}

// fallbackUpstreamMapper mplements the upstreamMapper interface
// Used by the fallback plugin
type fallbackUpstreamMapper struct {
	upstreamMap map[int]proxy.Upstream
}

func (f *fallbackUpstreamMapper) Add(rcode int, u proxy.Upstream) (oldU proxy.Upstream) {
	var ok bool
	if oldU, ok = f.upstreamMap[rcode]; !ok {
		f.upstreamMap[rcode] = u
	}

	return
}

func (f *fallbackUpstreamMapper) Get(rcode int) (u proxy.Upstream, ok bool) {
	u, ok = f.upstreamMap[rcode]
	return
}

func newFallbackUpstreamMapper() upstreamMapper {
	return &fallbackUpstreamMapper{upstreamMap: make(map[int]proxy.Upstream)}
}

// proxyCreator creates a proxy with the specified upstream
type proxyCreator interface {
	Create(trace plugin.Handler, upstream proxy.Upstream) plugin.Handler
}

// fallbackProxyCreator implements the proxyCreator interface
// Used by the fallback plugin to create proxy using specified for upstream
type fallbackProxyCreator struct{}

func (f fallbackProxyCreator) Create(trace plugin.Handler, upstream proxy.Upstream) plugin.Handler {
	return &proxy.Proxy{Trace: trace, Upstreams: &[]proxy.Upstream{upstream}}
}
