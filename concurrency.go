package main

import (
	"fmt"
	"net/http"
	"sync"
)

var urls = []string{
	"http://yahoo.com",
	"http://google.com",
	"http://zurmo.com",
	"http://www.gravity.com",
	"http://facebook.com",
	"http://www.tripadvisor.com",
	"http://www.booking.com",
	"http://aaaro02.com",
}

type HttpResponse struct {
	url      string
	response http.Response
	error    error
}

var httpResponses []HttpResponse // Slice to store responses

var wg, allUrlsAddedToChannel sync.WaitGroup

const numberOfWorkers = 5 // Max number of threads for fetching websites

func main() {
	wg.Add(numberOfWorkers)         // We need this, to ensure that main process do not end before all threads ends
	allUrlsAddedToChannel.Add(1)    // Check if all urls are added to channel
	channel := make(chan string, 6) // Create buffered channel of 10 elements

	// Add all urls to channel, concurently. Here we can eventually have more then one worker
	// In case we would do some IO job(for example reading files), we could have more then one routine, but in this case
	// we would need to care which routine will process which data.
	go addUrlsToChannel(urls, channel)

	// Create few routines for processing website.
	for i := 1; i <= numberOfWorkers; i++ {
		go processUrl(channel, i)
	}

	// Wait for all urls to be added to channel before closing it
	allUrlsAddedToChannel.Wait()
	close(channel)
	wg.Wait() // Wait for all workers to ends

	// Output data
	for _, httpResponse := range httpResponses {
		if httpResponse.error == nil {
			fmt.Println(httpResponse.url)
		} else {
			fmt.Println("Couldn't fetch site " + httpResponse.url)
		}
	}
}

// Add all urls to channel
func addUrlsToChannel(urls []string, channel chan string) {
	defer allUrlsAddedToChannel.Done()
	for _, url := range urls {
		channel <- url
	}
}

// Process single url
func processUrl(channel <-chan string, index int) {
	defer wg.Done()
	for {
		url, ok := <-channel
		if !ok {
			fmt.Printf("Ending process %d\n", index)
			return
		}
		fmt.Printf("Start fetching url %s with process %d\n", url, index)
		httpResponse := getResponse(url)
		httpResponses = append(httpResponses, httpResponse)
	}
}

// Get status of url(website)
func getResponse(url string) (httpResponse HttpResponse) {
	response, error := http.Get(url)
	httpResponse.error = error
	httpResponse.url = url
	if error == nil {
		httpResponse.response = *response
	}
	return
}

