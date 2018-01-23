package fallback

import (
	"fmt"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/proxy"
	"strings"

	"github.com/mholt/caddy"
	"github.com/miekg/dns"
)

func init() {
	caddy.RegisterPlugin("fallback", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	t := dnsserver.GetConfig(c).Handler("trace")
	f := &Fallback{trace: t, rules: make(map[int]proxy.Upstream)}

	for c.Next() {
		var rcode string
		if !c.Dispenser.Args(&rcode) {
			return c.ArgErr()
		}

		fmt.Printf("Found rcode %s\n", rcode)
		rc, ok := dns.StringToRcode[strings.ToUpper(rcode)]
		if !ok {
			return fmt.Errorf("%s is not a valid rcode", rcode)
		}

		fmt.Printf("Maps to rcode %d\n", rc)

		u, err := proxy.NewStaticUpstream(&c.Dispenser)
		if err != nil {
			return plugin.Error("fallback", err)
		}

		f.rules[rc] = u
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		f.Next = next
		return f
	})

	return nil
}
