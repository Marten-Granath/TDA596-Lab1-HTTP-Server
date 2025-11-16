/*
 * ::::::::::::::::
 * :: To-Do-List ::
 * ::::::::::::::::

 * Database Implementation

 * HTTP Request Handling
 - Handle GET Requests
 - Handle POST Requests
 - All other requests --> HTTP Code 501 "Not Implemented" Server Error

 * Compilation should yield binary called http_server
 - First argument: port to listen from
 - Note: No assumption made to IP address of server or incomming clients

 * Socket Connection Implementation
 - listen for incomming connections
 - listen on port specified from command line
 - new client request: spawn GO routine to handle it (max 10)
   - when exceeded we wait until a child process exits
 - Once connected:
   - Read client data
   - Check for properly formatted HTTP request
   - Accept the following formats for files:
     - html
	 - txt
	 - gif
	 - jpeg
	 - jpg
	 - css
   - Respond and transmit the following:
	 - text/html
	 - text/plain
	 - image/gif
	 - image/jpeg
	 - image/jpg
	 - text/css
   - For any other extension --> HTTP Code 400 "Bad Request"
   - If requested file doesn't exist --> HTTP Code 404 "Not Found"
   - If headers are not formatted correctly or any other error occurs --> HTTP Code 400 "Bad Request"

   - You should use the net library for networking, NOT the net/http library!
     - Ex: net.Listen("tcp", address)
   - You may use the net/http for parsing and working with http request objects (i.e not the networking part)
     - Ex: Dont use the following:
	       * http.ListenAndServe(...)
		   * http.Listen(...)
		   * http.Serve(...)

   - Set up some tests!
*/

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
 * ::::::::::::::::::::::
 * :: Helper Functions ::
 * ::::::::::::::::::::::
 */

func GetLocalIP() net.IP {
	connection, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	localAddress := connection.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func initializeServer() (string, *errgroup.Group) {
	var port string = ":2000"
	goroutines := new(errgroup.Group)
	goroutines.SetLimit(10)
	fmt.Printf("Server initialized\n")
	return port, goroutines
}

func startServer(port string) net.Listener {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server begins to listen on address: %s \n", listener.Addr().String())
	return listener
}

func awaitConnection(listener net.Listener) net.Conn {
	connection, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server accepts client request on %s from %s \n", connection.LocalAddr().String(), connection.RemoteAddr().String())
	return connection
}

func handleClientRequest(connection net.Conn, number int) error {
	fmt.Printf("Server handles client request %d from %s \n", number, connection.RemoteAddr().String())
	
	connection.Close()
	return errors.New("Placeholder")
}

/*
 * :::::::::::::::::::
 * :: Main Function ::
 * :::::::::::::::::::
 */

func main() {
	go func() { // TEMPORARY: Spawn server as goroutine, allowing us to run client requests concurrently
		fmt.Printf("Server starts...\n")

		// Server initialization
		// * Variable definitions and assignments
		port, goroutines := initializeServer()

		// Server start
		// * Listen on given port and filter connection type to tcp
		listener := startServer(port)
		defer listener.Close()

		// Main Server Loop
		var number int = 0 // TEMPORARY: Check routine number serviced
		for {
			// Wait for incomming connection & accept it
			// * Establishes socket connection on client request
			connection := awaitConnection(listener)

			// Start goroutine to concurrently handle client request
			// * Limited to 10 goroutines at any given time, see initializeServer()
			number++ // TEMPORARY: Check routine number serviced
			goroutines.Go(func() error {
				return handleClientRequest(connection, number)
			})
		}
	}()

	// TEMPORARY: Infinite loop to simulate a constant stream of requests
	for {
		//connection ,err := net.Dial("tcp", GetLocalIP().String()+":2000")
		//if err != nil {
		//	log.Fatal(err)
		//}

		resp, err := http.Get("https://localhost:2000")
		if err != nil {
			log.Fatalln(err)
		}
		//We Read the response body on the line below.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		resp.Body.Close()
		//Convert the body to type string
		sb := string(body)
		log.Printf(sb)

		//fmt.Printf("Client sends request from %s \n", GetLocalIP().String())
		time.Sleep(15 * time.Second)
	}
}
