package state

import (
	"fmt"
	"math/rand"
	"os"

	s "github.com/brodeynewman/wordle/internal/storage"

	"github.com/c-bata/go-prompt"
)

type State struct {
	guesses    int
	maxGuesses int
	chosenWord string
	hasWon     bool
}

type StateMachine interface {
	UpdateGuess()
}

func suggestions(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		// {Text: "users", Description: "Store the username and age"},
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (st *State) UpdateGuess() {
	fmt.Println("INC")

	(*st).guesses += 1
}

func NewState(store s.Store) State {
	words := store.Get()
	rn := rand.Intn(len(words))

	chosen := words[rn]

	s := State{
		guesses:    1,
		maxGuesses: 6,
		chosenWord: chosen,
	}

	return s
}

func handleInput(phrase string, st State) {
	switch phrase {
	case "exit":
		os.Exit(0)
	}

	if phrase != st.chosenWord {
		fmt.Println("INC?")
		st.UpdateGuess()
	}
}

func initGame(st State) {
	for st.guesses < 6 || !st.hasWon {
		fmt.Println("HMM", st)

		t := prompt.Input("> ", suggestions)
		handleInput(t, st)
	}

	if st.hasWon {
		fmt.Println("Nice Job!! You're a genius. You guessed the word:", st.chosenWord)
	} else {
		fmt.Println("Hm... You failed to guess the word:", st.chosenWord)
		fmt.Println("Better luck next time.")
	}
}

func Init(store s.Store) {
	fmt.Println("------------")
	fmt.Println("Welcome to Wordle! I have a 5 character word. Your job is to guess it within 6 Guesses. Lets go!")

	st := NewState(store)
	initGame(st)
}
