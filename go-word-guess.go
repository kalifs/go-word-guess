package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

func main() {
	argsWithoutProg := os.Args[1:]
	var language string
	if argsWithoutProg[0] != "" {
		language = argsWithoutProg[0]
	} else {
		language = "english"
	}

	dictionary := loadDictionary(&language)
	rand_source := rand.NewSource(time.Now().UnixNano())
	term := dictionary[rand.New(rand_source).Intn(len(dictionary)+1)]
	givenWord := term[0]

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Let's begin!")
	fmt.Printf("Definition: %s\n", term[1])
	for guessNum := 1; guessNum <= 5; guessNum++ {
		guessedWord := readWord(reader, guessNum)
		printColoredWord(givenWord, guessedWord)

		if guessedWord == givenWord {
			fmt.Printf("\033[1;32m%s\033[0m\n", "Congratulations you've won!")
			return
		}
	}
	fmt.Printf("\033[1;31m%s\033[0m\n", "You've failed!")
	fmt.Printf("Correct word was: %s\n", givenWord)
}

func readWord(reader *bufio.Reader, guessNum int) string {
	for {
		fmt.Printf("Guess[%d]: ", guessNum)
		guessedWordRaw, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		guessedWord := strings.TrimSuffix(guessedWordRaw, "\n")
		if utf8.RuneCountInString(guessedWord) != 5 {
			fmt.Println("Try again, exactly 5 letters!")
		} else {
			return guessedWord
		}
	}

}

func loadDictionary(language *string) [][]string {
	file, err := os.Open(fmt.Sprintf("data/%s.tsv", *language))
	check(err)

	defer file.Close()
	scanner := bufio.NewScanner(file)
	var line string
	var results [][]string

	for scanner.Scan() {
		line = scanner.Text()
		parts := strings.Split(line, "\t")
		results = append(results, parts)
	}
	return results
}

func printColoredWord(givenWord string, guessedWord string) {
	colors := map[string]string{
		"=": "\033[1;32m%s\033[0m",
		"~": "\033[1;33m%s\033[0m",
		"x": "\033[1;31m%s\033[0m",
	}

	matches := positionalMatches(givenWord, guessedWord)

	output := make([]string, 5)
	args := make([]interface{}, 5)
	gussedLetters := strings.Split(guessedWord, "")

	for i, matchType := range matches {
		output[i] = colors[matchType]
		if matchType == "x" {
			args[i] = "_"
		} else {
			args[i] = gussedLetters[i]
		}
	}
	fmt.Printf(strings.Join(output, ""), args...)
	fmt.Println("")
}

func containsLetter(letters []string, givenLetter string) int {
	for i, letter := range letters {
		if letter == givenLetter {
			return i
		}
	}
	return -1
}

func positionalMatches(givenWord string, guessedWord string) [5]string {
	correctLetters := strings.Split(givenWord, "")
	guessedLetters := strings.Split(guessedWord, "")
	var matches [5]string

	for i, letter := range correctLetters {
		if letter == guessedLetters[i] {
			matches[i] = "="
			correctLetters[i] = ""
		} else if correctIndex := containsLetter(correctLetters, guessedLetters[i]); correctIndex > -1 {
			matches[i] = "~"
			correctLetters[correctIndex] = ""
		} else {
			matches[i] = "x"
		}
	}

	return matches
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
