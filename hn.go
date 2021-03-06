package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	resultsCount int
)

func init() {
	const (
		defaultResultsCount = 30
		usage               = ""
	)
	flag.IntVar(&resultsCount, "count", defaultResultsCount, usage)
	flag.IntVar(&resultsCount, "c", defaultResultsCount, usage)
}

type Item struct {
	Score int    `json:"score"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func main() {
	flag.Parse()
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var topIDs []int
	json.Unmarshal(body, &topIDs)
	ids := make(chan int)
	done := make(chan bool)
	go func(ids <-chan int) {
		for id := range ids {
			url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json?print=pretty", id)
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			var item Item
			json.Unmarshal(body, &item)
			fmt.Printf("%s (%d)\n%s\n", item.Title, item.Score, item.URL)
			time.Sleep(time.Millisecond * 10)
		}
		close(done)
	}(ids)
	for _, id := range topIDs[:resultsCount] {
		ids <- id
	}
	close(ids)
	<-done
}
