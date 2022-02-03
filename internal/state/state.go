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

func suggestions(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exits you from the game."},
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func printGreen(l byte) {
	fmt.Print("\033[42m\033[1;30m")
	fmt.Printf(" %c ", l)
	fmt.Print("\033[m\033[m")
}

func printYellow(l byte) {
	fmt.Print("\033[43m\033[1;30m")
	fmt.Printf(" %c ", l)
	fmt.Print("\033[m\033[m")
}

func printGrey(l byte) {
	fmt.Print("\033[40m\033[1;37m")
	fmt.Printf(" %c ", l)
	fmt.Print("\033[m\033[m")
}

func printToConsole(st *State, p string) {
	for i := 0; i < len(p); i++ {
		x := p[i]
		y := st.chosenWord[i]

		if string(x) == string(y) {
			printGreen(x)
		} else if strings.Contains(st.chosenWord, string(x)) {
			printYellow(x)
		} else {
			printGrey(x)
		}
	}

	fmt.Println()
}

func (st *State) updateGuess(p string) {
	if len(p) > len(st.chosenWord) {
		fmt.Println("Invalid guess. You guessed a word that is longer than 5 characters.")
		(*st).guesses += 1

		return
	}

	printToConsole(st, p)

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

func (st *State) setWin() {
	(*st).hasWon = true
}

func handleInput(phrase string, st *State) {
	switch phrase {
	case "exit":
		os.Exit(0)
	}

	if phrase != st.chosenWord {
		st.updateGuess(phrase)
	} else {
		st.setWin()
	}
}

func getGuessText(st *State) string {
	var sb strings.Builder

	sb.WriteString("Guess ")
	sb.WriteString(strconv.FormatInt(int64(st.guesses), 10) + "/")
	sb.WriteString(strconv.FormatInt(int64(st.maxGuesses), 10) + " > ")

	return sb.String()
}

func announceWin(st *State) string {
	var sb strings.Builder

	if st.guesses == 1 {
		sb.WriteString("Amazing! Either you're a genius, or you cheated! You guessed the word: ")
	} else if st.guesses > 1 && st.guesses < 4 {
		sb.WriteString("Great Job!! You're an above average player! You guessed the word: ")
	} else {
		sb.WriteString("Good Job!! You aren't a complete dummy! You guessed the word: ")
	}

	sb.WriteString(st.chosenWord)
	sb.WriteString(" in ")
	sb.WriteString(strconv.FormatInt(int64(st.guesses), 10))

	if st.guesses > 1 {
		sb.WriteString(" guesses!")
	} else {
		sb.WriteString(" guess!")
	}

	return sb.String()
}

func initGame(st *State) {
	for st.guesses <= 6 && !st.hasWon {
		guess := getGuessText(st)

		t := prompt.Input(guess, suggestions)
		handleInput(t, st)
	}

	if st.hasWon {
		fmt.Println(announceWin(st))
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
