package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var Status = make(map[string]string)

type Website struct {
	URL string `json:"URL"`
}

var Websites []Website

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Assignment!")
}

func handleRequests() {
	fmt.Println("Server is starting on port 8080...")

	http.HandleFunc("/", homePage)

	http.HandleFunc("/postSites", returnAllWebsites)
	http.HandleFunc("/getSites", postAllWebsites)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func main() {
	handleRequests()
}

func returnAllWebsites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var websites []Website
	json.Unmarshal(reqBody, &websites)
	c := make(chan string)
	for _, link := range websites {
		go checkStatus(link.URL, c)
	}
	for l := range c {
		go func(link string) {
			time.Sleep(60 * time.Second)
			checkStatus(link, c)
		}(l)

	}
	json.NewEncoder(w).Encode(websites)
}

func postAllWebsites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := r.URL.Query().Get("name")
	c := make(chan string)
	if len(param) > 0 {
		fmt.Printf("Status of %s: %s", param, checkStatus(param, c)[param])
		// json.NewEncoder(w).Encode(Status[param])
	}
	json.NewEncoder(w).Encode(Status)

}

func checkStatus(link string, c chan string) map[string]string {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "Website is Down!")
		Status[link] = "Website is Down!"
		c <- link
		return Status
	}
	fmt.Println(link, "Website is Running!")
	Status[link] = "Website is Running!"

	c <- link
	return Status
}
