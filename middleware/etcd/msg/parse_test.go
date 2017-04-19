package msg

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		serv         Service
		expectedWhat int
	}{
		{serv: Service{Host: "example.org"}, expectedWhat: Host},
		{serv: Service{Host: "127.0.0.1"}, expectedWhat: IPv4},
		{serv: Service{Host: "2000::3"}, expectedWhat: IPv6},
		{serv: Service{Host: "2000..3"}, expectedWhat: Host},
		{serv: Service{Host: "127.0.0.257"}, expectedWhat: Host},
	}

	for i, tc := range tests {
		_, what := tc.serv.ParseHost()
		if what != tc.expectedWhat {
			t.Errorf("Test %d: Expected what %v, but got %v", i, tc.expectedWhat, what)
		}
	}

}
