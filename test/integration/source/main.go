package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
const TOKENSIZE = 21

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(postTokens))
}

func postTokens(w http.ResponseWriter, r *http.Request) {
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 1
	}
	if _, misbehave := os.LookupEnv("MISBEHAVE"); misbehave {
		// Change the query by +/-1 because we expect the caller to be able to deal with it
		size = size + rand.Intn(3) - 2
	}
	if size < 1 {
		size = 1
	}
	for i := 0; i < size; i++ {
		for j := 0; j < TOKENSIZE; j++ {
			w.Write([]byte(string(CHARSET[rand.Intn(len(CHARSET))])))
		}
		w.Write([]byte("\n"))
	}
}
