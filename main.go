package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	envHost   = "ECHO_SERVER_HOST"
	envPort   = "ECHO_SERVER_PORT"
	envPrefix = "ECHO_SERVER_PREFIX"
)

type config struct {
	host   string
	port   string
	prefix string
}

func getConfig() config {
	host := os.Getenv(envHost)
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv(envPort)
	if port == "" {
		port = "9002"
	}

	prefix := os.Getenv(envPrefix)
	if prefix == "" {
		prefix = "ECHO: "
	}

	return config{
		host:   host,
		port:   port,
		prefix: prefix,
	}
}

func getServerInfo(conn net.Conn) (string, error) {
	// Get server hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Get all network interfaces to show all available IPs
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}
	}

	// Get the local address that the client connected to
	localAddr := conn.LocalAddr().(*net.TCPAddr)

	// Build the info string
	var info strings.Builder
	info.WriteString("\n")
	info.WriteString("╔══════════════════════════════════════╗\n")
	info.WriteString("║        TCP Echo Server Started        ║\n")
	info.WriteString("╠══════════════════════════════════════╣\n")
	info.WriteString(fmt.Sprintf("║ %-20s: %-15s ║\n", "Hostname", hostname))
	info.WriteString(fmt.Sprintf("║ %-20s: %-15s ║\n", "Listening on", localAddr.IP))
	info.WriteString(fmt.Sprintf("║ %-20s: %-15d ║\n", "Port", localAddr.Port))
	if len(ips) > 0 {
		info.WriteString(fmt.Sprintf("║ %-20s: %-15s ║\n", "Available IPs", ips[0]))
		for _, ip := range ips[1:] {
			info.WriteString(fmt.Sprintf("║ %-20s  %-15s ║\n", "", ip))
		}
	}
	info.WriteString(fmt.Sprintf("║ %-20s: %-15s ║\n", "Current Time", time.Now().Format("2006-01-02 15:04:05")))
	info.WriteString("╚══════════════════════════════════════╝\n\n")

	return info.String(), nil
}

func handleConnection(conn net.Conn, cfg config) {
	defer func() {
		log.Printf("Closing connection from %s", conn.RemoteAddr().String())
		conn.Close()
	}()

	// Get client address
	clientAddr := conn.RemoteAddr().String()
	localAddr := conn.LocalAddr().String()
	
	log.Printf("New connection from %s to %s", clientAddr, localAddr)

	// Send welcome message and server info
	welcome := "\n=== Welcome to TCP Echo Server! ===\n"
	welcome += "Type anything and it will be echoed back.\n"
	welcome += "Type 'quit' to exit.\n"
	welcome += "================================\n\n"
	
	log.Printf("Sending welcome message to %s", clientAddr)
	
	// Get and send server info
	info, err := getServerInfo(conn)
	if err != nil {
		log.Printf("Error getting server info: %v", err)
		info = "Error getting server info\n"
	}
	
	// Send welcome message and server info
	_, err = conn.Write([]byte(welcome + info + "\n> "))
	if err != nil {
		log.Printf("Error writing to %s: %v", clientAddr, err)
		return
	}
	log.Printf("Welcome message sent to %s", clientAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if strings.ToLower(message) == "quit" {
			conn.Write([]byte("Goodbye!\n"))
			break
		}

		response := fmt.Sprintf("%s%s\n> ", cfg.prefix, message)
		_, err := conn.Write([]byte("\n" + response))
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from connection: %v", err)
	}

	log.Printf("Client disconnected: %s", clientAddr)
}

func main() {
	cfg := getConfig()
	address := fmt.Sprintf("%s:%s", cfg.host, cfg.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	// Get and display server info on startup
	hostname, _ := os.Hostname()
	addrs, _ := net.InterfaceAddrs()
	var ips []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║        TCP Echo Server Started        ║")
	fmt.Println("╠══════════════════════════════════════╣")
	fmt.Printf("║ %-20s: %-15s ║\n", "Hostname", hostname)
	fmt.Printf("║ %-20s: %-15s ║\n", "Listening on", strings.Split(address, ":")[0])
	fmt.Printf("║ %-20s: %-15s ║\n", "Port", strings.Split(address, ":")[1])
	if len(ips) > 0 {
		fmt.Printf("║ %-20s: %-15s ║\n", "Available IPs", ips[0])
		for _, ip := range ips[1:] {
			fmt.Printf("║ %-20s  %-15s ║\n", "", ip)
		}
	}
	fmt.Printf("║ %-20s: %-15s ║\n", "Start Time", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("║ %-20s: %-15s ║\n", "Prefix", cfg.prefix)
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()

	log.Printf("TCP Echo Server is running on %s", address)
	log.Printf("Using prefix: '%s'", cfg.prefix)
	log.Println("Ready to accept connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, cfg)
	}
}
