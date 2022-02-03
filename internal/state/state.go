package state

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

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
		{Text: "exit", Description: "Exits you from the game."},
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (st *State) UpdateGuess() {
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

func handleInput(phrase string, st *State) {
	switch phrase {
	case "exit":
		os.Exit(0)
	}

	if phrase != st.chosenWord {
		st.UpdateGuess()
	}
}

func getGuessText(st *State) string {
	var sb strings.Builder

	sb.WriteString("Guess ")
	sb.WriteString(strconv.FormatInt(int64(st.guesses), 10) + "/")
	sb.WriteString(strconv.FormatInt(int64(st.maxGuesses), 10) + " > ")

	return sb.String()
}

func initGame(st *State) {
	for st.guesses <= 6 && !st.hasWon {
		guess := getGuessText(st)

		t := prompt.Input(guess, suggestions)
		handleInput(t, st)
	}

	if st.hasWon {
		fmt.Println("Nice Job!! You're a genius. You guessed the word:", st.chosenWord)
	} else {
		fmt.Println("Hm... You failed to guess the word " + "'" + st.chosenWord + "'" + ". Better luck next time!")
	}
}

func Init(store s.Store) {
	fmt.Println("------------")
	fmt.Println("Welcome to Wordle! I have a 5 character word. Your job is to guess it within 6 Guesses. Lets go!")

	st := NewState(store)
	initGame(&st)
}
