package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"

	_ "embed"

	"github.com/fatih/color"
)

type WordPallet struct {
	Word   string
	Pallet []Letter
}
type Letter struct {
	Position int
	Guessed  bool
	Value    rune
}

var (
	attempts = 5
	//go:embed dictionary.txt
	rawWords   string
	dictionary []string
)

func main() {
	fmt.Println("Wordie...")
	dictionary = strings.Split(rawWords, "\n")
	idx := rand.Intn(len(dictionary))
	word := dictionary[idx]
	if os.Getenv("DEBUG") != "" {
		fmt.Printf("DEBUG enabled, solution is: %s\n", word)
	}

	// create new pallet.
	wp := WordPallet{
		Word: word,
	}
	wp.InitializePallet()
	// display pallet.
	wp.Display()
	// start game loop.
	wp.Play(attempts)
}

func (wp *WordPallet) Display() {
	for i, l := range wp.Pallet {
		// print different if its last letter.
		last := i == len(wp.Pallet)-1
		val := '_'
		if l.Guessed {
			val = l.Value
		}
		if last {
			fmt.Printf("%c\n", val)
		} else {
			fmt.Printf("%c ", val)
		}
	}
}

func (wp *WordPallet) InitializePallet() error {
	if wp.Word == "" {
		return errors.New("empty word for pallet init")
	}
	wp.Pallet = []Letter{}

	for i, val := range wp.Word {
		l := Letter{
			Value:    val,
			Position: i,
			Guessed:  false,
		}
		wp.Pallet = append(wp.Pallet, l)
	}
	return nil
}

// Search returns whether the rune can be found in the pallet at all
// and if it is found, returns its positions in the pallet.
func (wp *WordPallet) Search(r rune) (bool, []int) {
	found := false
	postitions := []int{}
	for i, v := range wp.Pallet {
		if v.Value == r {
			found = true
			postitions = append(postitions, i)
		}
	}
	return found, postitions
}

func (wp *WordPallet) Play(attempts int) {
	guesses := 0
	for {
		//   ask for guess.
		g := Prompt("")
		//   register guess.
		err := wp.RegisterGuess(g)
		if err != nil {
			fmt.Println("error:", err)
			// mulligan!
			continue
		} else {
			guesses++
		}

		if wp.IsSolved() {
			fmt.Println("SOLVED!")
			switch guesses {
			case 3:
				fmt.Println("GET OUTTA HERE FOURS!")
			case 4:
				fmt.Println("STUPID FOURS!")
			}
			return
		} else if guesses >= attempts {
			fmt.Printf("GAME OVER. The wordie was: %s\n", wp.Word)
			return
		}
	}
}

func (wp *WordPallet) IsSolved() bool {
	solved := true
	for _, v := range wp.Pallet {
		if !v.Guessed {
			return false
		}
	}
	return solved
}

func (wp *WordPallet) RegisterGuess(raw string) error {
	guess := sanitize(raw)
	if len(guess) != len(wp.Pallet) {
		return fmt.Errorf("bad length: %d, use length: %d", len(guess), len(wp.Pallet))
	}
	// TODO check dictionary validity of guess
	debug := os.Getenv("DEBUG")
	if debug == "" {
		if !slices.Contains(dictionary, guess) {
			return fmt.Errorf("guess not found in dictionary")
		}
	}
	m := map[rune]int{}

	for i, v := range guess {
		color.Set(color.FgHiWhite)
		found, positions := wp.Search(v)
		if found && slices.Contains(positions, i) {
			// Green
			wp.Pallet[i].Guessed = true
			color.Set(color.BgHiGreen, color.FgBlack)
		} else if found && !slices.Contains(positions, i) {
			// Yellow
			color.Set(color.BgHiYellow, color.FgBlack)

			count := len(positions)
			num := m[v]
			// fmt.Printf("num: %d, count:%d", num, count)
			if num >= count {
				color.Unset()
			}

		} else {
			// bust
			color.Unset()
		}

		m[v]++
		fmt.Printf(" %c ", v)
		color.Unset()
	}
	fmt.Println()
	return nil
}

func Prompt(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	text, _ := reader.ReadString('\n')
	return text
}

func sanitize(raw string) string {
	s := strings.ToLower(raw)
	s = strings.Trim(s, "\n")
	return s
}
