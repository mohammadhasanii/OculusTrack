package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var totalTime int64 = 0

type filteredWriter struct {
	writer io.Writer
}

func (fw *filteredWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	if strings.Contains(msg, "TLS handshake error") {
		return len(p), nil
	}
	return fw.writer.Write(p)
}

func generateCertificate() error {
	if _, err := os.Stat("localhost.crt"); err == nil {
		if _, err := os.Stat("localhost.key"); err == nil {
			fmt.Println("Certificate files already exist.")
			return nil
		}
	}

	fmt.Println("Generating SSL certificate...")

	var ipAddresses []net.IP
	ipAddresses = append(ipAddresses, net.IPv4(127, 0, 0, 1))

	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("could not get network interfaces: %v", err)
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && !ip.IsLoopback() {
				if ip.To4() != nil {
					ipAddresses = append(ipAddresses, ip)
				}
			}
		}
	}
	fmt.Println("Certificate will be valid for the following IPs:", ipAddresses)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("error generating private key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"Local Dev Server"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"localhost"},
		IPAddresses: ipAddresses,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("error creating certificate: %v", err)
	}

	certOut, err := os.Create("localhost.crt")
	if err != nil {
		return fmt.Errorf("error creating cert file: %v", err)
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	keyOut, err := os.Create("localhost.key")
	if err != nil {
		return fmt.Errorf("error creating key file: %v", err)
	}
	defer keyOut.Close()
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	fmt.Println("SSL certificate generated successfully!")
	return nil
}

func main() {
	if err := generateCertificate(); err != nil {
		log.Fatalf("Failed to generate certificate: %v", err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/update-time", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed. Use POST to update.")
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		t := r.FormValue("time")
		if _, err := fmt.Sscanf(t, "%d", &totalTime); err != nil {
			http.Error(w, "Invalid time format", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Total time updated to %d", totalTime)
	})

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	customLogger := log.New(&filteredWriter{writer: os.Stderr}, "", log.LstdFlags)

	server := &http.Server{
		Addr:         ":8443",
		TLSConfig:    tlsConfig,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		ErrorLog:     customLogger,
	}

	fmt.Println("🚀 HTTPS Server starting on https://localhost:8443")
	if err := server.ListenAndServeTLS("localhost.crt", "localhost.key"); err != nil {
		log.Fatalf("❌ HTTPS Server failed to start: %v", err)
	}
}
