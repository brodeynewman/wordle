package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const DEFAULT_WORD_LEN = 5
const CACHE_PATH = "./.cache/words.txt"
const DICTIONARY_URL = "https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt"

type Words []string

type Store struct {
	usingCache bool
	words      Words
}

func (s *Store) Get() Words {
	// can do some manipulation in here at some point
	return s.words
}

func format(w Words) Words {
	words := Words{}

	for i := 0; i < len(w); i++ {
		if len(w[i]) > 6 || len(w[i]) < 1 {
			continue
		}

		t := w[i][1:len(w[i])]

		// only use words that have the same length as our default
		if len(t) == DEFAULT_WORD_LEN {
			words = append(words, t)
		}
	}

	return words
}

func formatForStorage(w Words) (Words, Words) {
	words := Words{}
	cache := Words{}

	for i := 0; i < len(w); i++ {
		// only use words that have the same length as our default
		if len(w[i]) == DEFAULT_WORD_LEN {
			words = append(words, w[i])
			cache = append(cache, w[i]+",")
		}
	}

	return words, cache
}

func (s *Store) loadFromCache() {
	fmt.Println("Loading dictionary from cache...")

	buffer := &bytes.Buffer{}
	abs, _ := filepath.Abs(CACHE_PATH)
	data, _ := os.ReadFile(abs)

	gob.NewDecoder(buffer).Decode(&s.words)

	f := strings.Split(string(data), ",")

	s.words = format(f)
	s.usingCache = true
}

func checkForCache() bool {
	abs, err := filepath.Abs(CACHE_PATH)

	if err != nil {
		fmt.Println("An error occurred when fetching cache path:", err)
	}

	if _, err := os.Stat(abs); err == nil {
		return true
	}

	return false
}

func fetchAndPrepareWords() (Words, Words) {
	fmt.Println("Caching dictionary...")

	resp, err := http.Get(DICTIONARY_URL)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Cannot read body:", body)
	}

	words := strings.Split(string(body), "\r\n")

	return formatForStorage(words)
}

func storeWords(w Words) {
	abs, _ := filepath.Abs(CACHE_PATH)

	buffer := &bytes.Buffer{}

	gob.NewEncoder(buffer).Encode(w)
	byteSlice := buffer.Bytes()

	derr := os.Mkdir(".cache", 0755)
	err := os.WriteFile(abs, byteSlice, 0644)

	if derr != nil {
		panic(err)
	}

	if err != nil {
		fmt.Println("An error occurred when writing cache file", err)
	}
}

func newStore() Store {
	c := checkForCache()

	if !c {
		w, stored := fetchAndPrepareWords()
		storeWords(stored)
		s := Store{usingCache: false, words: w}

		return s
	} else {
		fmt.Println("Using cached dictionary...")

		s := Store{}
		s.loadFromCache()

		return s
	}
}

func Init() Store {
	m := newStore()

	return m
}
