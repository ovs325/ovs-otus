package listsystem

import (
	"testing"

	tp "bruteforce/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestSearchEngineIp(t *testing.T) {
	subnetMap := map[string]tp.IsBlack{"192.1.1.0/25": true, "10.0.0.0/8": false}

	tests := []struct {
		ipStr     string
		isFound   bool
		isBlack   tp.IsBlack
		expectErr bool
	}{
		{"192.1.1.10", true, true, false},
		{"10.0.0.10", true, false, false},
		{"192.2.2.10", false, false, false},
		{"invalid-ip", false, false, true},
	}

	for _, test := range tests {
		isFound, isBlack, err := SearchEngineIP(subnetMap, test.ipStr)

		assert.Falsef(
			t,
			(err != nil) != test.expectErr,
			"SearchEngineIp(%s) returned error: %v; expected error: %v", test.ipStr, err != nil, test.expectErr,
		)
		assert.Falsef(
			t,
			isFound != test.isFound || isBlack != test.isBlack,
			"SearchEngineIp(%s) = (%v, %v); expected (%v, %v)", test.ipStr, isFound, isBlack, test.isFound, test.isBlack,
		)
	}
}

func TestSubnetNameVerifier(t *testing.T) {
	tests := []struct {
		subnetStr string
		expected  bool
	}{
		{"192.1.1.0/25", true},
		{"10.0.0.0/8", true},
		{"invalid-subnet", false},
	}

	for _, test := range tests {
		res := SubnetNameVerifier(test.subnetStr)
		assert.Equalf(t, test.expected, res, "SubnetNameVerifier(%s) = %v; expected %v", test.subnetStr, res, test.expected)
	}
}
