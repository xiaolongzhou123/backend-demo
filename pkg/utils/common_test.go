package utils

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	// ...
	addr := "http://18.18.2.2:9090/api/v1/query"
	query := "increase(ifInOctets{instance='18.18.1.1',ifName='GigabitEthernet0/0/6'}[24h])[7d:1d]"
	rs, ok := Query(addr, query)
	fmt.Println(rs, ok)
}
