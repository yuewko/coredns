package msg

import (
	"net"

	"github.com/miekg/dns"
)

const (
	IPv4 = iota + 1
	IPv6
	Host
)

// Type returns the DNS type of what is encoded in the Service s. The values returned are:
//
// dns.TypeMX: service encodes an MX record, host field assumed to be a name.
// dns.TypeTXT: service encodes an TXT record, host field assumed to be a name.
// dns.TypeA: the service's Host field contains an A record.
// dns.TypeAAAA: the service's Host field contains an AAAA record.
// dns.TypeANY: the service's Host field contains a name.
//
// Note that we first check for MX and then for TXT. In case of dns.TypeA and dns.TypeAAAA
// the returned address is normalized with To4() or To16() respectively.
func (s *Service) Type() (what uint16, normalized net.IP) {

	if s.Mail {
		return dns.TypeMX, nil
	}

	if len(s.Text) > 0 {
		return dns.TypeTXT, nil
	}

	ip := net.ParseIP(s.Host)

	switch {
	case ip == nil:
		return dns.TypeANY, nil

	case ip.To4() != nil:
		return dns.TypeA, ip.To4()

	case ip.To4() == nil:
		return dns.TypeAAAA, ip.To16()
	}
	// This should never be reached.
	return dns.TypeNone, nil
}
