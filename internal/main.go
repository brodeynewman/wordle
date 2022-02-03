package main

import (
	st "github.com/brodeynewman/wordle/internal/state"
	s "github.com/brodeynewman/wordle/internal/storage"
)

func main() {
	store := s.Init()
	st.Init(store)
}
