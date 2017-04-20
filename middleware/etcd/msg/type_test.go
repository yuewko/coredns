package msg

import (
	"testing"

	"github.com/miekg/dns"
)

func TestType(t *testing.T) {
	tests := []struct {
		serv         Service
		expectedType uint16
	}{
		{Service{Host: "example.org"}, dns.TypeANY},
		{Service{Host: "127.0.0.1"}, dns.TypeA},
		{Service{Host: "2000::3"}, dns.TypeAAAA},
		{Service{Host: "2000..3"}, dns.TypeANY},
		{Service{Host: "127.0.0.257"}, dns.TypeANY},
		{Service{Host: "127.0.0.257", Mail: true}, dns.TypeMX},
		{Service{Host: "127.0.0.257", Mail: true, Text: "a"}, dns.TypeMX},
		{Service{Host: "127.0.0.257", Mail: false, Text: "a"}, dns.TypeTXT},
	}

	for i, tc := range tests {
		what, _ := tc.serv.Type()
		if what != tc.expectedType {
			t.Errorf("Test %d: Expected what %v, but got %v", i, tc.expectedType, what)
		}
	}

}
