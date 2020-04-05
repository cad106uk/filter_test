package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type JokeDATA struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

type jokeCache struct {
	cache map[int]JokeDATA
	lock  sync.RWMutex
}

func (js *jokeCache) AddOne(val JokeDATA) {
	js.lock.Lock()
	defer js.lock.Unlock()
	js.cache[val.ID] = val
}

func (js *jokeCache) AddMany(vals []JokeDATA) {
	js.lock.Lock()
	defer js.lock.Unlock()
	for _, val := range vals {
		js.cache[val.ID] = val
	}
}

func (js *jokeCache) GetJoke(key int) (JokeDATA, bool) {
	js.lock.RLock()
	defer js.lock.RUnlock()
	val, ok := js.cache[key]
	return val, ok
}

func (js *jokeCache) RandomJoke() JokeDATA {
	js.lock.RLock()
	defer js.lock.RUnlock()

	keys := reflect.ValueOf(js.cache).MapKeys()
	// Get random key
	r := rand.Intn(len(keys))
	// Cast from ReflectValue to int64,
	j := keys[r].Int()
	// then cast the int64 to an int
	return js.cache[int(j)]
}

var knownJokesCache = jokeCache{
	cache: make(map[int]JokeDATA),
}

func updateJokeCache() {
	response, err := http.Get("https://official-joke-api.appspot.com/jokes/programming/ten")

	if err != nil {
		log.Fatal(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var rawRead []JokeDATA
	err = json.Unmarshal(responseData, &rawRead)
	if err != nil {
		log.Fatal(err)
	}

	knownJokesCache.AddMany(rawRead)
}

func init() {
	updateJokeCache()

	ticker := time.Tick(5 * time.Minute)
	go func() {
		for _ = range ticker {
			updateJokeCache()
		}
	}()
}

func randomHandle(w http.ResponseWriter, _ *http.Request) {
	joke := knownJokesCache.RandomJoke()
	j, err := json.Marshal(joke)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func specificHandle(w http.ResponseWriter, req *http.Request) {
	id_str := strings.TrimPrefix(req.URL.Path, "/get_joke/")
	key, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	joke, ok := knownJokesCache.GetJoke(key)
	if !ok {
		http.Error(w, "No joke! With this id", http.StatusNotFound)
		return
	}

	j, err := json.Marshal(joke)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func main() {
	http.HandleFunc("/get_joke/", specificHandle)
	http.HandleFunc("/random_joke/", randomHandle)
	http.ListenAndServe(":8080", nil)
	fmt.Printf("%+v", knownJokesCache)
}
