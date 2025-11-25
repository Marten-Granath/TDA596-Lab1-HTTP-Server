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
	"os"
	"slices"
	"strings"

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

func initializeServer(port string) (string, *errgroup.Group, *map[string]FileData, *map[string]FileData, []string) {
	port = ":" + port
	goroutines := new(errgroup.Group)
	goroutines.SetLimit(10)
	var textMap = make(map[string]FileData)
	var imageMap = make(map[string]FileData)
	var textMapPointer = &textMap
	var imageMapPointer = &imageMap
	allowedExtensions := []string{}
	allowedExtensions = append(allowedExtensions, "html", "txt", "css", "jpg", "jpeg", "gif")

	//fmt.Printf("Server initialized\n")
	return port, goroutines, textMapPointer, imageMapPointer, allowedExtensions
}

func startServer(port string) net.Listener {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Server begins to listen on address: %s \n", listener.Addr().String())
	return listener
}

func awaitConnection(listener net.Listener) net.Conn {
	connection, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Server accepts client request on %s from %s \n", connection.LocalAddr().String(), connection.RemoteAddr().String())
	return connection
}

// Note that these lengths only work for the files specified in the question details
// * To accomodate all possible extensions a for loop could be utilized to find out min and max extension sizes
func isValidExt(fileName string, allowedExtensions []string) bool {
	var templist []string
	for _, element := range allowedExtensions {
		templist = append(templist, ("." + element))
	}
	return (slices.Contains(templist, fileName[(len(fileName)-4):]) || slices.Contains(templist, fileName[(len(fileName)-3):]))
}

func isValid(contentType string, allowedExtensions []string) bool {
	return (slices.Contains(allowedExtensions, contentType))
}

func isMatching(fileName string, contentType string) bool {
	var size int = len(contentType)
	var size2 int = len(fileName)
	return (fileName[(size2-size-1):] == "."+contentType)
}

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

func handleGetRequest(request *http.Request, allowedExtensions []string, textMapPointer *map[string]FileData, imageMapPointer *map[string]FileData) []byte {
	u, err := url.Parse(request.RequestURI)
	if err != nil {
		return buildErrorResponse("HTTP/1.1", "400 Bad Request")
	}

	var response []byte
	var fileName = fmt.Sprint(strings.TrimPrefix(u.Path, "/"))
	var body []byte
	var contentType string

	if !isValidExt(fileName, allowedExtensions) {
		response = buildErrorResponse("HTTP/1.1", "400 Bad Request")
	} else if (isMatching(fileName, (*textMapPointer)[fileName].ContentType) || isMatching(fileName, (*imageMapPointer)[fileName].ContentType)) == false {
		response = buildErrorResponse("HTTP/1.1", "404 Not Found")
	} else {
		//fmt.Println((*textMapPointer)["dog.jpg"])
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
		// if isValid(request.Header.Get("Content-Type"), allowedExtensions) {
		// 	response = buildErrorResponse("HTTP/1.1", "404 Not Found")
		// } else if body != nil {
		response = buildGetResponse("HTTP/1.1", "200 OK", contentType, body)
		// } else {
		// 	response = buildErrorResponse("HTTP/1.1", "400 Bad Request")
		// }
	}
	return response
}

func handlePostRequest(request *http.Request, allowedExtensions []string, textMapPointer *map[string]FileData, imageMapPointer *map[string]FileData) []byte {
	u, err := url.Parse(request.RequestURI)
	if err != nil {
		return buildErrorResponse("HTTP/1.1", "400 Bad Request")
	}

	var response []byte
	var fileName = strings.TrimPrefix(u.Path, "/")

	if !isMatching(fileName, request.Header.Get("Content-Type")) {
		response = buildErrorResponse("HTTP/1.1", "400 Bad Request")
	} else {
		// Post for text file
		if slices.Contains(allowedExtensions[:3], request.Header.Get("Content-Type")) {
			bodyBytes, err := io.ReadAll(request.Body)
			if err != nil {
				return buildErrorResponse("HTTP/1.1", "400 Bad Request")
			}
			var fd FileData
			fd.Body = bodyBytes
			fd.ContentType = request.Header.Get("Content-Type")

			(*textMapPointer)[fileName] = fd
			response = buildResponse("HTTP/1.1", "200 OK", "text/"+request.Header.Get("Content-Type"))
			// Post for image file
		} else if slices.Contains(allowedExtensions[3:6], request.Header.Get("Content-Type")) {
			bodyBytes, err := io.ReadAll(request.Body)
			if err != nil {
				return buildErrorResponse("HTTP/1.1", "400 Bad Request")
			}
			var fd FileData
			fd.Body = bodyBytes
			fd.ContentType = request.Header.Get("Content-Type")

			(*imageMapPointer)[fileName] = fd
			response = buildResponse("HTTP/1.1", "200 OK", "image/"+request.Header.Get("Content-Type"))

		} else {
			response = buildErrorResponse("HTTP/1.1", "400 Bad Request")
		}
	}

	return response
}

func handleClientRequest(conn net.Conn, textMapPointer *map[string]FileData, imageMapPointer *map[string]FileData, allowedExtensions []string) error {
	defer conn.Close()
	//fmt.Printf("Server handles client request from %s\n", conn.RemoteAddr().String())
	var response []byte
	// Parse the HTTP request (allowed by assignment)
	reader := bufio.NewReader(conn)
	request, err := http.ReadRequest(reader)

	//fmt.Printf("Server received %s request: %s\n", request.Method, request.Body)

	// Handle different http methods
	switch httpMethod := request.Method; httpMethod {
	case "GET":
		response = handleGetRequest(request, allowedExtensions, textMapPointer, imageMapPointer)
	case "POST":
		response = handlePostRequest(request, allowedExtensions, textMapPointer, imageMapPointer)
	case "HEAD":
		response = buildErrorResponse("HTTP/1.1", "501 Not Implemented")
	case "PUT":
		response = buildErrorResponse("HTTP/1.1", "501 Not Implemented")
	case "DELETE":
		response = buildErrorResponse("HTTP/1.1", "501 Not Implemented")
	case "CONNECT":
		response = buildErrorResponse("HTTP/1.1", "501 Not Implemented")
	case "OPTIONS":
		response = buildErrorResponse("HTTP/1.1", "501 Not Implemented")
	case "TRACE":
		response = buildErrorResponse("HTTP/1.1", "501 Not Implemented")
	case "PATCH":
		response = buildErrorResponse("HTTP/1.1", "501 Not Implemented")
	default:
		response = buildErrorResponse("HTTP/1.1", "400 Bad Request")
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
	// First command line argument is set as port
	var port string = os.Args[1]

	//fmt.Printf("Server starts...\n")

	// Server initialization
	// * Variable definitions and assignments
	port, goroutines, textMapPointer, imageMapPointer, allowedExtensions := initializeServer(port)

	// Server start
	// * Listen on given port and filter connection type to tcp
	listener := startServer(port)
	defer listener.Close()

	// Main Server Loop
	for {
		// Wait for incomming connection & accept it
		// * Establishes socket connection on client request
		connection := awaitConnection(listener)

		// Start goroutine to concurrently handle client request
		// * Limited to 10 goroutines at any given time, see initializeServer()
		goroutines.Go(func() error {
			return handleClientRequest(connection, textMapPointer, imageMapPointer, allowedExtensions)
		})
	}
}
