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
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
 * :::::::::::::
 * :: Structs ::
 * :::::::::::::
 */

 type FileData struct {
    ContentType string
    Body        []byte
}


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

func initializeServer() (string, *errgroup.Group, *map[string]FileData, *map[string]FileData, []string) {
	var port string = ":2000"
	goroutines := new(errgroup.Group)
	goroutines.SetLimit(10)
	var textMap = make(map[string]FileData)
	var imageMap = make(map[string]FileData)
	var textMapPointer = &textMap
	var imageMapPointer = &imageMap
	allowedExtensions := []string{}
	allowedExtensions = append(allowedExtensions, "html", "txt", "css", "jpg", "jpeg", "gif")

	fmt.Printf("Server initialized\n")
	return port, goroutines, textMapPointer, imageMapPointer, allowedExtensions
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
//, contentLength string, body string
func buildResponse(version string, status string, contentType string) []byte {

	header := version + " " + status + "\r\n" +
		"Content-Type: " + contentType + "\r\n" + 
		"\r\n"
	return []byte(header)
}

func buildGetResponse(version string, status string, contentType string, body []byte) []byte {
	header := version + " " + status + "\r\n" +
		"Content-Type: " + contentType + "\r\n" +
		 "\r\n"
	return append([]byte(header), body...)
}

func buildErrorResponse(version string, status string) []byte {
	header := version + " " + status + "\r\n" + "\r\n"
	return []byte(header)
}

func handleClientRequest(conn net.Conn, number int, textMapPointer *map[string]FileData, imageMapPointer *map[string]FileData, allowedExtensions []string) error {
	defer conn.Close()

	fmt.Printf("Server handles client request %d from %s\n",
		number, conn.RemoteAddr().String())

	// Parse the HTTP request (allowed by assignment)
	reader := bufio.NewReader(conn)
	request, err := http.ReadRequest(reader)
	if err != nil {
		log.Println("Failed to parse HTTP request:", err)
		return err
	}
	var response []byte

	fmt.Printf("Received %s request: %s\n", request.Method, request.Body)

	switch httpMethod := request.Method; httpMethod {
	case "GET":
		fmt.Println("Server: 200 OK")
		u, err := url.Parse(request.RequestURI)
		if err != nil {
			fmt.Println("parse error:", err)
			return err
		}
		var fileName = fmt.Sprint(strings.TrimPrefix(u.Path, "/"))
		var body []byte
		var contentType string
		fmt.Println((*textMapPointer)["dog.jpg"])
		for k := range *textMapPointer { 
    		if fileName == k {
				if (*textMapPointer)[fileName].ContentType == "txt" {
					body = (*textMapPointer)[fileName].Body
					contentType = "text/" + "plain"
				} else {
					body = (*textMapPointer)[fileName].Body
					contentType = "text/" + (*textMapPointer)[fileName].ContentType
				}
			}
		}
		for k := range *imageMapPointer { 
    		if fileName == k {
				body = (*imageMapPointer)[fileName].Body
				contentType = "image/" + (*imageMapPointer)[fileName].ContentType
			}
		}
		if  body != nil {
			response = buildGetResponse("HTTP/1.1", "200 OK", contentType, body)
		} else {
			fmt.Println("400 Bad Request")
			response = buildErrorResponse("HTTP/1.1", "400 Bad Request")
		}

	case "POST":
		fmt.Println("Server: 200 OK")
		u, err := url.Parse(request.RequestURI)
		if err != nil {
			fmt.Println("parse error:", err)
			return err
		}
		//fmt.Sprintln(strings.TrimPrefix(u.Path, "/"))
		var fileName = strings.TrimPrefix(u.Path, "/")

		if slices.Contains(allowedExtensions[:3], request.Header.Get("Content-Type")) {
			fmt.Println("text file detected")
			bodyBytes, err := io.ReadAll(request.Body)
			if err != nil {
				fmt.Println("error reading body:", err)
			}
			var fd FileData
			fd.Body = bodyBytes
			fd.ContentType = request.Header.Get("Content-Type")
			(*textMapPointer)[fileName] = fd
			response = buildResponse("HTTP/1.1", "200 OK", "text/" + request.Header.Get("Content-Type"))
											
		} else if slices.Contains(allowedExtensions[3:6], request.Header.Get("Content-Type")) {
			fmt.Println("image file detected")
			bodyBytes, err := io.ReadAll(request.Body)
			if err != nil {
				fmt.Println("error reading body:", err)
			}
			var fd FileData
			fd.Body = bodyBytes
			fd.ContentType = request.Header.Get("Content-Type")
			(*imageMapPointer)[fileName] = fd
			response = buildResponse("HTTP/1.1", "200 OK", "image/" + request.Header.Get("Content-Type"))

		} else {
			fmt.Println("400 Bad Request")
			response = buildErrorResponse("HTTP/1.1", "400 Bad Request")
		}

	case "HEAD":
		fmt.Println("Server: 501 Not Implemented")
	case "PUT":
		fmt.Println("Server: 501 Not Implemented")
	case "DELETE":
		fmt.Println("Server: 501 Not Implemented")
	case "CONNECT":
		fmt.Println("Server: 501 Not Implemented")
	case "OPTIONS":
		fmt.Println("Server: 501 Not Implemented")
	case "TRACE":
		fmt.Println("Server: 501 Not Implemented")
	case "PATCH":
		fmt.Println("Server: 501 Not Implemented")
	default:
		fmt.Println("400 Bad Request")
		
	}

	// Write response to connection
	_, err = conn.Write([]byte(response))
	if err != nil {
		return err
	}

	return nil
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
		port, goroutines, textMapPointer, imageMapPointer, allowedExtensions := initializeServer()

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
				return handleClientRequest(connection, number, textMapPointer, imageMapPointer, allowedExtensions)
			})
		}

	}()

	// TEMPORARY: Infinite loop to simulate a constant stream of requests
	for {
		//connection ,err := net.Dial("tcp", GetLocalIP().String()+":2000")
		//if err != nil {
		//	log.Fatal(err)
		//}

		resp, err := http.Post("http://localhost:2000/germanenginering.txt", "txt", strings.NewReader("Boo"))
		if err != nil {
			log.Fatalln(err)
		}

		// Print status
		fmt.Println("Status:", resp.Status)
		fmt.Println("Status Code:", resp.StatusCode)

		// Print headers
		fmt.Println("Headers:")
		for k, v := range resp.Header {
			fmt.Printf("%s: %v\n", k, v)
		}

		//We Read the response body on the line below.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		
		resp.Body.Close()
		//Convert the body to type string
		sb := string(body)
		log.Print(sb)

		time.Sleep(5 * time.Second)
		
		// Test 2
		fmt.Println("GET REQUEST IS FROM HERE ON OUT BSCHES!!!")
		resp2, err := http.Get("http://localhost:2000/germanenginering.txt")
		if err != nil {
			log.Fatalln(err)
		}

		// Print status
		fmt.Println("Status:", resp2.Status)
		fmt.Println("Status Code:", resp2.StatusCode)

		// Print headers
		fmt.Println("Headers:")
		for k, v := range resp2.Header {
			fmt.Printf("%s: %v\n", k, v)
		}

		//We Read the response body on the line below.
		body2, err := io.ReadAll(resp2.Body)
		if err != nil {
			log.Fatalln(err)
		}
		
		resp2.Body.Close()
		//Convert the body to type string
		sb2 := string(body2)
		log.Print(sb2)

		//fmt.Printf("Client sends request from %s \n", GetLocalIP().String())
		time.Sleep(15 * time.Second)

	}

}
