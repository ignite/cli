package availableport_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/ignite/cli/ignite/pkg/availableport"
)

func TestFind(t *testing.T) {
	for n := 0; n < 50; n++ {
		ports, err := availableport.Find(n)
		if err != nil {
			t.Errorf("Error Find: %v", err)
		}

		if len(ports) != n {
			t.Errorf("Expected %d ports, found %d", n, len(ports))
		}

		// Verifies ports are in range 44000 - 55000
		minPort := 44000
		maxPort := 55000
		for _, port := range ports {
			if port < minPort || port > maxPort {
				t.Errorf("Port %d out of range %d-%d", port, minPort, maxPort)
			}
		}

		// Verifica que los puertos encontrados no están siendo utilizados actualmente
		for _, port := range ports {
			conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
			if err == nil {
				conn.Close()
				t.Errorf("Port %d in use", port)
			}
		}
	}

	for idx := 11001; idx < 12000; idx++ {
		_, err := availableport.Find(idx)
		// You can't request more than 11000 ports (55000-44000 = 11000)
		if err == nil {
			t.Errorf("Ports overflow. %d bigger than the limit", idx)
		}
	}

}
