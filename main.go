package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {

	target := "127.0.0.1"
	fmt.Printf("[*] Starting active scan on %s...\n", target)

	var wg sync.WaitGroup

	for port := 1; port <= 1024; port++ {
		wg.Add(1)

		go func(p int) {
			defer wg.Done()
			scanPort(target, p)
		}(port)

	}

	wg.Wait()
	fmt.Println("[*] Scan completed.")
}

func scanPort(ip string, port int) {
	address := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)

	if err != nil {
		return
	}

	conn.Close()
	fmt.Printf("[+] PORT FOUND: %d\n", port)

}
