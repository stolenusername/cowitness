//June 7, 2023 - Added port 443 and the ability to serve files.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/miekg/dns"
)

const (
	HTTPPort        = 80
	HTTPSPort       = 443
	DNSPort         = 53
	DNSResponseIP   = "127.0.0.1"   //Change this to the desired IP address
	DNSResponseName = "domain.com." //Change to the desired domain name
	DefaultTTL      = 3600
)

func main() {
	// Get the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

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
						Hdr: dns.RR_Header{Name: DNSResponseName, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: DefaultTTL},
						Ns:  "ns1.domain.com.", //Change to the desired name servers
					},
					&dns.NS{
						Hdr: dns.RR_Header{Name: DNSResponseName, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: DefaultTTL},
						Ns:  "ns2.domain.com.", //Change to the desired name servers
					})
			} else if r.Question[0].Qtype == dns.TypeA {
				// Check if the request is for A records
				response.Answer = append(response.Answer,
					&dns.A{
						Hdr: dns.RR_Header{Name: subdomain + "." + DNSResponseName, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: DefaultTTL},
						A:   net.ParseIP(DNSResponseIP),
					})
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

	// Wait indefinitely
	select {}
}
