// Package fallback implements a fallback plugin for CoreDNS
package fallback

import (
	"context"

	"github.com/miekg/dns"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"
	"github.com/coredns/coredns/plugin/proxy"
)

// Fallback plugin
type Fallback struct {
	Next  plugin.Handler
	trace plugin.Handler
	rules map[int]proxy.Upstream
}

// ServeDNS implements the plugin.Handler interface.
func (f Fallback) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	nw := nonwriter.New(w)
	rcode, err := plugin.NextOrFailure(f.Name(), f.Next, ctx, nw, r)
	if err != nil {
		return rcode, err
	}
	//WTF? WHY DOES PROXY ALWAYS RETURN 0
	if rcode == 0 {
		rcode = nw.Msg.Rcode
	}
	if u, ok := f.rules[rcode]; ok {
		p := &proxy.Proxy{Trace: f.trace, Upstreams: &[]proxy.Upstream{u}}
		return p.ServeDNS(ctx, w, r)
	}
	w.WriteMsg(nw.Msg)
	return rcode, nil
}

// Name implements the Handler interface.
func (f Fallback) Name() string { return "fallback" }
