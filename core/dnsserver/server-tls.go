package dnsserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/miekg/dns"
)

// ServerTLS represents an instance of a TLS-over-DNS-server.
type ServerTLS struct {
	*Server
	tlsConfig *tls.Config
}

// NewServerTLS returns a new CoreDNS TLS server and compiles all plugin in to it.
func NewServerTLS(addr string, group []*Config) (*ServerTLS, error) {
	s, err := NewServer(addr, group)
	if err != nil {
		return nil, err
	}

	return &ServerTLS{Server: s}, nil
}

func (s *ServerTLS) createTLSListener(innerListener net.Listener) net.Listener {
	return tls.NewListener(innerListener, s.tlsConfig)
}

// Serve implements caddy.TCPServer interface.
func (s *ServerTLS) Serve(l net.Listener) error {
	s.m.Lock()

	if s.tlsConfig != nil {
		l = s.createTLSListener(l)
	}
	// Only fill out the TCP server for this one.
	s.server[tcp] = &dns.Server{Listener: l, Net: "tcp-tls", Handler: dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		ctx := context.Background()
		s.ServeDNS(ctx, w, r)
	})}
	s.m.Unlock()

	return s.server[tcp].ActivateAndServe()
}

// ServePacket implements caddy.UDPServer interface.
func (s *ServerTLS) ServePacket(p net.PacketConn) error { return nil }

// Listen implements caddy.TCPServer interface.
func (s *ServerTLS) Listen() (net.Listener, error) {
	// The *tls* plugin must make sure that multiple conflicting
	// TLS configuration return an error: it can only be specified once.
	for _, conf := range s.zones {
		// Should we error if some configs *don't* have TLS?
		s.tlsConfig = conf.TLSConfig
	}

	l, err := net.Listen("tcp", s.Addr[len(TransportTLS+"://"):])
	if err != nil {
		return nil, err
	}
	return l, nil
}

// ListenPacket implements caddy.UDPServer interface.
func (s *ServerTLS) ListenPacket() (net.PacketConn, error) { return nil, nil }

// OnStartupComplete lists the sites served by this server
// and any relevant information, assuming Quiet is false.
func (s *ServerTLS) OnStartupComplete() {
	if Quiet {
		return
	}

	for zone, config := range s.zones {
		fmt.Println(TransportTLS + "://" + zone + ":" + config.Port)
	}
}
