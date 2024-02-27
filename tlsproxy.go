package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strings"
)

// Define the address of the upstream server
const upstreamAddr = "www.cnn.com:443"	//"127.0.0.1:8443"  	//"upstream-server-address:443"
const cert_pem = "../cert/certificate.crt"    //path/to/certificate.crt
const cert_key = "../cert/private.key"	//  path/to/private.key

/*
We can genertate self signed if one - like burp does
openssl req -newkey rsa:2048 -nodes -keyout private.key -out request.csr
openssl x509 -req -days 365 -in request.csr -signkey private.key -out certificate.crt

*/

func main() {
	
	// Create a listener for incoming TLS connections
	listener, err := net.Listen("tcp", "localhost:4443")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Proxy server is listening on localhost:8443")

	// Loop indefinitely to accept incoming connections
	for {
		// Accept incoming client connections
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		// Handle client connection in a separate goroutine
		go handleClient(clientConn)
	}
}

func handleClient(clientConn net.Conn) {
        // Certificates and configuration for TLS handshake with clients
        cert, _ := tls.LoadX509KeyPair(cert_pem, cert_key)
	config := tls.Config{Certificates: []tls.Certificate{cert}}


    defer clientConn.Close()

    // Perform TLS handshake with the client
    clientTLSConn := tls.Server(clientConn, &config)
    if err := clientTLSConn.Handshake(); err != nil {
        fmt.Println("Error performing TLS handshake with client:", err.Error())
        return
    }

    // Connect to the upstream server
    upstreamConn, err := tls.Dial("tcp", upstreamAddr, &tls.Config{
        InsecureSkipVerify: true, // Skip TLS certificate verification (not recommended in production)
    })
    if err != nil {
        fmt.Println("Error connecting to upstream server:", err.Error())
        return
    }
    defer upstreamConn.Close()

    // Continuously read data from the client, analyze it, and forward it to the upstream server
    go func() {
        defer clientTLSConn.Close()

        buffer := make([]byte, 512) // Example buffer size
        for {
            n, err := clientTLSConn.Read(buffer)
            if err != nil {
                if err != io.EOF {
                    fmt.Println("Error reading from client:", err.Error())
                }
                break
            }

            // Analyze the data (replace this with your analysis logic)
            analyzedData := analyzeData(buffer[:n],n ,1)

            // Forward the analyzed data to the upstream server
            if _, err := upstreamConn.Write(analyzedData); err != nil {
                fmt.Println("Error writing to upstream:", err.Error())
                break
            }
        }
    }()

    // Continuously read data from the upstream server, analyze it, and forward it to the client
    buffer := make([]byte, 512) // Example buffer size
    for {
        n, err := upstreamConn.Read(buffer)
        if err != nil {
            if err != io.EOF {
                fmt.Println("Error reading from upstream:", err.Error())
            }
            break
        }

        // Analyze the data (replace this with your analysis logic)
        analyzedData := analyzeData(buffer[:n],n,0)

        // Forward the analyzed data to the client
        if _, err := clientTLSConn.Write(analyzedData); err != nil {
            fmt.Println("Error writing to client:", err.Error())
            break
        }
    }
}

func analyzeData(data []byte, length int, fromClient int) []byte {
    message := string(data[:length])
    newStr := strings.Replace(message, "GET", "POST", -1)
    
    // Perform analysis logic here
    if fromClient == 1 {
    	fmt.Printf("message recieved from client: %s\n", string(data[:length]))
    	}
    	
    if fromClient == 1 {
    	fmt.Printf("message recieved from server: %s\n", string(data[:length]))
    }
    
    newData := []byte(newStr)
    return newData
}
