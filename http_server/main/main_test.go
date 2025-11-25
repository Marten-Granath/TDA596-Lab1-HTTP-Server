package main

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// Assumptions:
// * Supplied Port is valid
// * FileName is larger than ContentType
// * No posts with existing resources

// --- Goroutine Tests --- //
// func TestMaxGoroutines(t *testing.T) {
// 	initializeTests("1000")

// 	for i := 0; i < 10; i++ {
// 		go func() {}()
// 	}
// }

// func TestAboveMaxGoroutines(t *testing.T) {
// 	initializeTests("1001")
// }

// --- GET Tests -- //

func TestGetRequestForExistingResource(t *testing.T) {
	initializeTests("2001")
	resp, _ := http.Get("http://localhost:2001/test.txt")
	if !(resp.Status == "200 OK") {
		t.Errorf("Expected status 200 OK, received status %s \n", resp.Status)
	}
}

func TestGetRequestWithMissingResource(t *testing.T) {
	initializeTests("2002")
	resp, _ := http.Get("http://localhost:2002/test2.jpg")
	if !(resp.Status == "404 Not Found") {
		t.Errorf("Expected status 404 Not Found, received status %s \n", resp.Status)
	}
}

func TestGetRequestWithInvalidExtension(t *testing.T) {
	initializeTests("2003")
	resp, _ := http.Get("http://localhost:2003/test.invalidext")
	if !(resp.Status == "400 Bad Request") {
		t.Errorf("Expected status 400 Bad Request, received status %s \n", resp.Status)
	}
}

func TestGetRequestWithMissingResourceAndInvalidExtension(t *testing.T) {
	initializeTests("2004")
	resp, _ := http.Get("http://localhost:2004/missingResource.invalidext")
	if !(resp.Status == "400 Bad Request") {
		t.Errorf("Expected status 400 Bad Request, received status %s \n", resp.Status)
	}
}

func TestGetRequestReponseTxt(t *testing.T) {
	initializeTests("2005")
	resp, _ := http.Get("http://localhost:2005/test.txt")
	if !(resp.Header.Get("Content-Type") == "text/plain") {
		t.Errorf("Expected Content-Type text/plain, received Content-Type %s \n", resp.Header.Get("Content-Type"))
	}

}
func TestGetRequestReponseCss(t *testing.T) {
	initializeTests("2006")
	resp, _ := http.Get("http://localhost:2006/test.css")
	if !(resp.Header.Get("Content-Type") == "text/css") {
		t.Errorf("Expected Content-Type text/css, received Content-Type %s \n", resp.Header.Get("Content-Type"))
	}
}

func TestGetRequestReponseGif(t *testing.T) {
	initializeTests("2007")
	resp, _ := http.Get("http://localhost:2007/test.gif")
	if !(resp.Header.Get("Content-Type") == "image/gif") {
		t.Errorf("Expected Content-Type image/img, received Content-Type %s \n", resp.Header.Get("Content-Type"))
	}

}

// -- POST Tests -- //

func TestPostRequestWithValidSignature(t *testing.T) {
	initializeTests("3000")
	resp, _ := http.Post("http://localhost:3000/test.jpg", "jpg", strings.NewReader("body"))
	if !(resp.Status == "200 OK") {
		t.Errorf("Expected status 200 OK, received status %s \n", resp.Status)

	}
}

func TestPostRequestWithInvalidExtension(t *testing.T) {
	initializeTests("3001")
	resp, _ := http.Post("http://localhost:3001/test.invalidext", "invalidext", strings.NewReader("body"))
	if !(resp.Status == "400 Bad Request") {
		t.Errorf("Expected status 400 Bad Request, received status %s \n", resp.Status)
	}
}

func TestPostRequestWithValidExtensionAndInvalidContentType(t *testing.T) {
	initializeTests("3002")
	resp, _ := http.Post("http://localhost:3002/test.img", "invalid", strings.NewReader("body"))
	if !(resp.Status == "400 Bad Request") {
		t.Errorf("Expected status 400 Bad Request, received status %s \n", resp.Status)
	}
}

func TestPostRequestWithInvalidExtensionAndValidContentType(t *testing.T) {
	initializeTests("3003")
	resp, _ := http.Post("http://localhost:3003/test.inv", "txt", strings.NewReader("body"))
	if !(resp.Status == "400 Bad Request") {
		t.Errorf("Expected status 400 Bad Request, received status %s \n", resp.Status)
	}
}

func TestPostRequestWithInvalidExtensionAndInvalidContentType(t *testing.T) {
	initializeTests("3004")
	resp, _ := http.Post("http://localhost:3004/test.inv", "inv", strings.NewReader("body"))
	if !(resp.Status == "400 Bad Request") {
		t.Errorf("Expected status 400 Bad Request, received status %s \n", resp.Status)
	}
}

// --- Not Implemented Tests --- //

func TestNotImplementedPutRequest(t *testing.T) {
	initializeTests("4000")
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", "http://localhost:4000/test.txt", strings.NewReader("body"))
	resp, _ := client.Do(req)
	if !(resp.Status == "501 Not Implemented") {
		t.Errorf("Expected status 501 Not Implemented, received status %s \n", resp.Status)
	}
}

func TestNotImplementedDeleteRequest(t *testing.T) {
	initializeTests("4001")
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", "http://localhost:4001/test.txt", nil)
	resp, _ := client.Do(req)
	if !(resp.Status == "501 Not Implemented") {
		t.Errorf("Expected status 501 Not Implemented, received status %s \n", resp.Status)
	}
}

func TestNotImplementedPatchRequest(t *testing.T) {
	initializeTests("4002")
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", "http://localhost:4002/test.txt", strings.NewReader("body"))
	resp, _ := client.Do(req)
	if !(resp.Status == "501 Not Implemented") {
		t.Errorf("Expected status 501 Not Implemented, received status %s \n", resp.Status)
	}
}

func TestNotImplementedHeadRequest(t *testing.T) {
	initializeTests("4003")
	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", "http://localhost:4003/test.txt", nil)
	resp, _ := client.Do(req)
	if !(resp.Status == "501 Not Implemented") {
		t.Errorf("Expected status 501 Not Implemented, received status %s \n", resp.Status)
	}
}

func TestNotImplementedOptionsRequest(t *testing.T) {
	initializeTests("4004")
	client := &http.Client{}
	req, _ := http.NewRequest("OPTIONS", "http://localhost:4004/test.txt", nil)
	resp, _ := client.Do(req)
	if !(resp.Status == "501 Not Implemented") {
		t.Errorf("Expected status 501 Not Implemented, received status %s \n", resp.Status)
	}
}
func TestNotImplementedTraceRequest(t *testing.T) {
	initializeTests("4005")
	client := &http.Client{}
	req, _ := http.NewRequest("TRACE", "http://localhost:4005/test.txt", nil)
	resp, _ := client.Do(req)
	if !(resp.Status == "501 Not Implemented") {
		t.Errorf("Expected status 501 Not Implemented, received status %s \n", resp.Status)
	}
}
func TestNotImplementedConnectRequest(t *testing.T) {
	initializeTests("4006")
	client := &http.Client{}
	req, _ := http.NewRequest("CONNECT", "http://localhost:4006/test.txt", nil)
	resp, _ := client.Do(req)
	if !(resp.Status == "501 Not Implemented") {
		t.Errorf("Expected status 501 Not Implemented, received status %s \n", resp.Status)
	}
}

// Helper Functions
func initializeTests(port string) {
	// Start Server
	go func() {
		os.Args[1] = port
		main()
	}()

	// Set Server Values
	var url string = "http://localhost:" + port
	http.Post(url+"/test.txt", "txt", strings.NewReader("test"))
	http.Post(url+"/test.gif", "gif", strings.NewReader("test"))
	http.Post(url+"/test.css", "css", strings.NewReader("test"))

	time.Sleep(time.Second / 4)
}
