package fallback

import (
	"context"
	"testing"

	"github.com/miekg/dns"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/pkg/healthcheck"
	"github.com/coredns/coredns/plugin/proxy"
	"github.com/coredns/coredns/plugin/test"
	"github.com/coredns/coredns/request"
)

// stubNextHandler is used to simulate a proxy plugin.
// It returns a stub Handler that returns the rcode and err specified when invoked.
func stubNextHandler(rcode int, err error) test.Handler {
	return test.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
		w.WriteMsg(&dns.Msg{})
		return rcode, err
	})
}

// testProxyCreator implements the proxyCreator interface.
// Used by unit test to verify that the proxyCreator.Create() method is called as expected
type testProxyCreator struct {
	expectedUpstream proxy.Upstream
	called           int
	t                *testing.T
}

func (c *testProxyCreator) Create(trace plugin.Handler, upstream proxy.Upstream) plugin.Handler {
	c.called++

	// Ensure that it is called with the expected Upstream
	if c.expectedUpstream != upstream {
		c.t.Fatal("Unexpected upstream passed to proxyCreator")
	}

	return &proxy.Proxy{Trace: trace, Upstreams: &[]proxy.Upstream{upstream}}
}

// dummyUpstream implements the proxy.Upstream interface
// It is used to fake an upstream used for creating a proxy.
type dummyUpstream struct {
	rcode int
}

func (u dummyUpstream) From() string {
	return ""
}

func (u dummyUpstream) Select() *healthcheck.UpstreamHost {
	return nil
}

func (u dummyUpstream) IsAllowedDomain(s string) bool {
	return false
}

func (u dummyUpstream) Exchanger() proxy.Exchanger {
	return dummyExchanger{}
}

func (u dummyUpstream) Stop() error {
	return nil
}

// dummyExchanger implements the proxy.Exchanger interface
// It is used solely to implement the dummyUpstream above
type dummyExchanger struct{}

func (e dummyExchanger) Exchange(ctx context.Context, addr string, state request.Request) (*dns.Msg, error) {
	return &dns.Msg{}, nil
}

func (e dummyExchanger) Protocol() string {
	return ""
}

func (e dummyExchanger) Transport() string {
	return ""
}

func (e dummyExchanger) OnStartup(p *proxy.Proxy) error {
	return nil
}

func (e dummyExchanger) OnShutdown(p *proxy.Proxy) error {
	return nil
}

// Test case for fallback
type fallbackTestCase struct {
	rcode            int            // rcode to be returned by the stub Handler
	expectedUpstream proxy.Upstream // this upstream is expected when testProxyCreator is called
}

func TestFallback(t *testing.T) {
	// dummy Upstreams for servicing a specific rcode
	dummyRefusedUpstream := &dummyUpstream{rcode: dns.RcodeRefused}
	dummyServeFailureUpstream := &dummyUpstream{rcode: dns.RcodeServerFailure}
	dummyNxDomainUpstream := &dummyUpstream{rcode: dns.RcodeNameError}
	dummyUpstreams := []*dummyUpstream{
		dummyRefusedUpstream,
		dummyServeFailureUpstream,
		dummyNxDomainUpstream,
	}

	testCases := []fallbackTestCase{
		{
			rcode:            dns.RcodeRefused,
			expectedUpstream: dummyRefusedUpstream,
		},
		{
			rcode:            dns.RcodeServerFailure,
			expectedUpstream: dummyServeFailureUpstream,
		},
		{
			rcode:            dns.RcodeNameError,
			expectedUpstream: dummyNxDomainUpstream,
		},
	}

	for _, tc := range testCases {
		handler := newFallback(nil)
		handler.Next = stubNextHandler(tc.rcode, nil)
		// add dummyUpstreams to upstream map according to the rcode field
		for _, u := range dummyUpstreams {
			handler.mapper.Add(u.rcode, u)
		}
		proxyCreator := &testProxyCreator{t: t, expectedUpstream: tc.expectedUpstream}
		handler.proxy = proxyCreator

		ctx := context.TODO()
		req := &dns.Msg{
			Question: []dns.Question{{
				Name:   "abc.com",
				Qclass: dns.ClassINET,
				Qtype:  dns.TypeA,
			}},
		}

		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		_, _ = handler.ServeDNS(ctx, rec, req)

		// Ensure that the proxyCreator is called once
		if proxyCreator.called != 1 {
			t.Fatalf("Expect proxy creator to be called once, but got '%d", proxyCreator.called)
		}
	}
}
