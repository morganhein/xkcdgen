package xkcdgen

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

//TODO: this barely works. Needs to be completely reworked.

var (
	defaultSymbols = "!@#$%^:;&*_-.+"
	r1             *rand.Rand
)

type Dictionary struct {
	fileLocation string
	lines        int
}

func Generate() (string, error) {
	return GenerateWithOptions(-1, 3, defaultSymbols, "dictionary.txt")
}

func GenerateWithOptions(maxLength, wordCount int, symbols, dictionaryLoc string) (string, error) {
	//defer func() {
	//	err := dictionaryLoc.Close()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()
	if wordCount == 0 || wordCount < -1 {
		return "", errors.New("wordCount must be greater than 0")
	}
	if maxLength == 0 {
		return "", errors.New("maxLength must be greater than 0 or -1 for infinite")
	}
	minChars := determineMinimumLength(wordCount)
	if maxLength < minChars && maxLength > 0 {
		return "", fmt.Errorf("a wordcount of %v requires a maxLength of at least %v", wordCount, minChars+1)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)

	wrapperSymbol := getSymbol(symbols)
	bufferSymbol := ""
	for bufferSymbol == "" || bufferSymbol == wrapperSymbol {
		bufferSymbol = getSymbol(symbols)
	}
	bufferNumber := r1.Intn(100)
	//todo: the bufferNumber should be in form 00 in all cases
	currentCharCount := 2 + 2 + 2 + len(fmt.Sprintf("%v", bufferNumber)) + wordCount
	if maxLength < currentCharCount && maxLength > 0 {
		return "", errors.New("maxLength is not long enough for the required words and various symbols")
	}
	words := make([]string, wordCount)
	var err error
	lines, err := getLength(dictionaryLoc)
	if err != nil {
		return "", err
	}
	for i := 0; i < wordCount; i++ {
		wordLength := calculateWordLength(i, maxLength, currentCharCount, wordCount)
		if wordLength <= 0 {
			return "", errors.New("unable to create a password with the desired wordLength and number of words")
		}
		words[i], err = getWord(wordLength, Dictionary{
			fileLocation: dictionaryLoc,
			lines:        lines,
		})
		if err != nil {
			return "", err
		}
		currentCharCount += len(words[i])
	}
	result := fmt.Sprintf("%s%s%v", wrapperSymbol, wrapperSymbol, bufferNumber)

	capital := false
	for _, v := range words {
		if capital {
			v = strings.ToUpper(v)
		}
		result += fmt.Sprintf("%s%s", bufferSymbol, v)
		capital = !capital
	}
	result += fmt.Sprintf("%s%v%s%s", bufferSymbol, bufferNumber, wrapperSymbol, wrapperSymbol)
	return result, nil
}

func determineMinimumLength(wordCount int) int {
	//wrapper + wrapper + number(2) + (bufferSymbol + 4) + bufferSymbol + number(2) + wrapper + wrapper
	min := 1 + 1 + 2 + wordCount*(1+4) + 1 + 2 + 1 + 1
	return min
}

func getSymbol(symbols string) string {
	return string(symbols[r1.Intn(len(symbols))])
}

func calculateWordLength(i, maxLength, currentCharCount, wordCount int) int {
	if maxLength == -1 {
		return 25
	}
	return maxLength - currentCharCount - (wordCount - (i + 1))
}

func getLength(original string) (int, error) {
	r, err := os.Open(original)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := r.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func getWord(maxLength int, dictionary Dictionary) (string, error) {
	if dictionary.lines == 0 {
		return "", errors.New("dictionary must have words in it")
	}
	getWordHelper := func(start int) (string, error) {
		r, err := os.Open(dictionary.fileLocation)
		if err != nil {
			return "", err
		}
		defer func() {
			err := r.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		scanner := bufio.NewScanner(r)
		current := -1
		for scanner.Scan() {
			current++
			if current < start {
				continue
			}
			text := scanner.Text()
			if len(text) <= maxLength && len(text) >= 3 {
				return text, nil
			}
		}
		return "", fmt.Errorf("could not grab a random word, maxLength was %v and random start was %v", maxLength, start)
	}
	result := ""
	lines := dictionary.lines
	var err error
	for lines > 100 && len(result) == 0 {
		start := r1.Intn(lines)
		result, err = getWordHelper(start)
		lines = lines / 2
	}
	if err != nil {
		return "", err
	}
	return result, nil
}
