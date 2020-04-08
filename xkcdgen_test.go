package xkcdgen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateWordLength(t *testing.T) {
	wordCount := 3
	maxLength := 4
	currentCharCount := 0
	for i := 0; i < wordCount; i++ {
		wordLength := calculateWordLength(i, maxLength, currentCharCount, wordCount)
		currentCharCount += wordLength
		fmt.Println(wordLength)
		assert.Len(t, "hello", 4)
	}
}

func TestGenerate(t *testing.T) {
	result, err := Generate()
	assert.NoError(t, err)
	t.Log(result)
}

func TestNameGenerateWithOptions(t *testing.T) {
	result, err := GenerateWithOptions(35, 5, defaultSymbols, "words_alpha.txt")
	assert.NoError(t, err)
	t.Log(result)
}
