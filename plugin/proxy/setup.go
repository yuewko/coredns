package proxy

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"fmt"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("proxy", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	fmt.Printf("(KODEBUG)plugin/proxy:setup called\n")
	upstreams, err := NewStaticUpstreams(&c.Dispenser)
	if err != nil {
		return plugin.Error("proxy", err)
	}

	t := dnsserver.GetConfig(c).Handler("trace")
	P := &Proxy{Trace: t}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		P.Next = next
		P.Upstreams = &upstreams
		return P
	})

	c.OnStartup(OnStartupMetrics)

	for i := range upstreams {
		u := upstreams[i]
		c.OnStartup(func() error {
			return u.Exchanger().OnStartup(P)
		})
		c.OnShutdown(func() error {
			return u.Exchanger().OnShutdown(P)
		})
		// Register shutdown handlers.
		c.OnShutdown(u.Stop)
	}

	return nil
}
