package msg

import (
	"net"
)

const (
	IPv4 = iota + 1
	IPv6
	Host
)

// ParseHost parses the Host field in s. It returns an address when the address resembles an IP
// address (done with net.ParseIP) or nil when the Host is a name. What differentiates between
// a hostname (Host), an IPv4 address (IPv4) or an IPv6 address (IPv6).
func (s *Service) ParseHost() (addr net.IP, what int) {

	ip := net.ParseIP(s.Host)

	switch {
	case ip == nil:
		// Host
		return nil, Host

	case ip.To4() != nil:
		// IPv4 address
		return ip.To4(), IPv4

	case ip.To4() == nil:
		// IPv6 address
		return ip.To16(), IPv6
	}
	// This should never be reached.
	return
}
