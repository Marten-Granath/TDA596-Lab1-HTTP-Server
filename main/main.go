package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
 * ::::::::::::::::::::::
 * :: Helper Functions ::
 * ::::::::::::::::::::::
 */

func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Bananas are for monkeys, drink gainomax instead, no sponsor")
}

/*
 * :::::::::::::::::::
 * :: Main Function ::
 * :::::::::::::::::::
 */

func main() {
	http.HandleFunc("/", httpHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
