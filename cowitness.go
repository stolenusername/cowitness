package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/miekg/dns"
)

const (
	HTTPPort   = 80
	HTTPSPort  = 443
	DNSPort    = 53
	DefaultTTL = 3600
)

var (
	DNSResponseIP   string // User-defined DNS response IP
	DNSResponseName string // User-defined DNS response name
	QuietMode       bool   // Flag to enable quiet mode
)

func main() {
	// Check if the program should run in quiet mode
	if len(os.Args) > 1 && os.Args[1] == "-q" {
		QuietMode = true
	}

	// Display the ASCII art banner unless running in quiet mode
	if !QuietMode {
		displayBanner()
	}

	// Get the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Ask the user for DNSResponseIP and store it as a variable
	fmt.Print("Enter the DNS response IP: ")
	fmt.Scanln(&DNSResponseIP)

	// Ask the user for DNSResponseName and store it as a variable
	fmt.Print("Enter the DNS response name: ")
	fmt.Scanln(&DNSResponseName)

	// Ask the user for DefaultTTL and store it as a variable
	fmt.Print("Enter the Default TTL: ")
	var DefaultTTL int
	fmt.Scanln(&DefaultTTL)

	// Create log files
	httpLogFile, err := os.OpenFile("./http.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer httpLogFile.Close()

	dnsLogFile, err := os.OpenFile("./dns.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer dnsLogFile.Close()

	// Create HTTP request logger
	httpLogger := log.New(httpLogFile, "", log.LstdFlags)

	// Start HTTP server on port 80
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log IP address, HTTP resource, and user agent to http.log
		ipAddress := strings.Split(r.RemoteAddr, ":")[0]
		requestResource := r.URL.Path
		userAgent := r.UserAgent()
		logMessage := fmt.Sprintf("IP address: %s, Resource: %s, User agent: %s\n", ipAddress, requestResource, userAgent)
		httpLogger.Println(logMessage)

		// Serve the requested file
		http.FileServer(http.Dir(rootDir)).ServeHTTP(w, r)
	})

	go func() {
		log.Printf("Starting HTTP server on port %d\n", HTTPPort)
		err := http.ListenAndServe(fmt.Sprintf(":%d", HTTPPort), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Start HTTP server on port 443
	go func() {
		log.Printf("Starting HTTPS server on port %d\n", HTTPSPort)
		err := http.ListenAndServe(fmt.Sprintf(":%d", HTTPSPort), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Start DNS server
	go func() {
		addr := fmt.Sprintf(":%d", DNSPort)
		server := &dns.Server{Addr: addr, Net: "udp"}

		dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
			ipAddress := w.RemoteAddr().(*net.UDPAddr).IP
			logMessage := fmt.Sprintf("IP address: %s, DNS request: %s\n", ipAddress, r.Question[0].Name)
			if _, err := dnsLogFile.WriteString(logMessage); err != nil {
				log.Println(err)
			}

			// Prepare the response
			response := new(dns.Msg)
			response.SetReply(r)
			response.Authoritative = true
			response.RecursionAvailable = true

			// Extract the subdomain from the DNS request
			domain := r.Question[0].Name
			subdomain := strings.TrimSuffix(domain, "."+DNSResponseName)

			// Check if the request is for NS records
			if r.Question[0].Qtype == dns.TypeNS {
				response.Answer = append(response.Answer,
					&dns.NS{
						Hdr: dns.RR_Header{Name: DNSResponseName, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: uint32(DefaultTTL)},
						Ns:  "ns1.domain.com.", //Change to the desired name servers
					},
					&dns.NS{
						Hdr: dns.RR_Header{Name: DNSResponseName, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: uint32(DefaultTTL)},
						Ns:  "ns2.domain.com.", //Change to the desired name servers
					})
			} else if r.Question[0].Qtype == dns.TypeA {
				// Check if the request is for A records
				if domain == DNSResponseName {
					// Request for the main domain
					response.Answer = append(response.Answer,
						&dns.A{
							Hdr: dns.RR_Header{Name: DNSResponseName, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: uint32(DefaultTTL)},
							A:   net.ParseIP(DNSResponseIP),
						})
				} else {
					// Request for a subdomain
					response.Answer = append(response.Answer,
						&dns.A{
							Hdr: dns.RR_Header{Name: subdomain + "." + DNSResponseName, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: uint32(DefaultTTL)},
							A:   net.ParseIP(DNSResponseIP),
						})
				}
			}

			// Send the response
			if err := w.WriteMsg(response); err != nil {
				log.Println(err)
			}
		})

		log.Printf("Starting DNS server on port %d\n", DNSPort)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Kill the DNS server process ID when the program is closed
	defer func() {
		pid := os.Getpid()
		cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", pid))
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	}()

	// Output a link to the URL that the user can click on
	log.Printf("Open the following URL in your browser:\n")
	log.Printf("http://localhost:%d\n", HTTPPort)

	// Wait indefinitely
	select {}
}

func displayBanner() {
	red := "\033[31m"
	reset := "\033[0m"
	banner := red + `
 	          ⢠⡄
	    	⣠⣤⣾⣷⣤⣄⡀⠀⠀⠀⠀
  @@@@@@     ⣴⡿⠋⠁⣼⡇⠈⠙⢿⣧⠀⠀⠀ @@@  @@@  @@@  @@@  @@@@@@@ @@@  @@@ @@@@@@@@  @@@@@@  @@@@@@
 !@@        ⣸⡟⠀⠀⠀⠘⠃⠀⠀⠀⢻⣇⠀⠀ @@!  @@!  @@!  @@!    @@!   @@!@!@@@ @@!      !@@     !@@    
 !@!     ⠰⠶⣿⡷⠶⠶⠀⠀⠀⠀⠶⠶⢾⣿⠶⠆  @!!  !!@  @!@  !!@    @!!   @!@@!!@! @!!!:!    !@@!!   !@@!! 
 :!!        ⢹⣧⠀⠀⠀⢠⡄⠀⠀⠀⣼⡏⠀   !:  !!:  !!   !!:    !!:   !!:  !!! !!:          !:!     !:!
  :: :: :    ⠹⣷⣆⡀⢸⡇⢀⣠⣾⠏⠀⠀⠀⠀  ::.:  :::    :       :    ::    :  : :: ::: ::.: :  ::.: :
               ⠈⠙⠛⣿⡿⠛⠋⠀⠀⠀⠀  
	           ⠘⠃⠀⠀
` + reset

	fmt.Print(banner)
	fmt.Println("             CoWitness - Tool for HTTP, HTTPS, and DNS Server")
	fmt.Println()
}
