package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type JokeJSON struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

var knownJokesCache = make(map[int]JokeJSON)

func main() {
	response, err := http.Get("https://official-joke-api.appspot.com/jokes/programming/ten")

	if err != nil {
		log.Fatal(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rawRead []JokeJSON
	err = json.Unmarshal(responseData, &rawRead)
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range rawRead {
		knownJokesCache[val.ID] = val
	}
	fmt.Printf("%+v", knownJokesCache)
}
