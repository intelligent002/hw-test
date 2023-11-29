package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordsCount struct {
	word  string
	count int
}

func makeWordsValid(input string) []string {
	input = strings.ToLower(input)
	wordsRaw := strings.Fields(input)
	wordsValid := make([]string, 0, len(wordsRaw))
	for _, word := range wordsRaw {
		word = strings.Trim(word, ",.!?\"'")
		if len(word) > 0 && word != "-" {
			wordsValid = append(wordsValid, word)
		}
	}
	return wordsValid
}

func makeWordsCounted(wordsRaw []string) map[string]int {
	wordsCounted := make(map[string]int)
	for _, word := range wordsRaw {
		wordsCounted[word]++
	}
	return wordsCounted
}

func makeWordsSorted(wordsCounted map[string]int) []wordsCount {
	wordsSorted := make([]wordsCount, 0, len(wordsCounted))
	for word, count := range wordsCounted {
		wordsSorted = append(wordsSorted, wordsCount{word: word, count: count})
	}
	sort.Slice(wordsSorted, func(i, j int) bool {
		if wordsSorted[i].count == wordsSorted[j].count {
			return wordsSorted[i].word < wordsSorted[j].word
		}
		return wordsSorted[i].count > wordsSorted[j].count
	})
	return wordsSorted
}

func makeWordsTop(wordsSorted []wordsCount, amount int) []string {
	wordsTop := make([]string, 0, amount)
	for index, word := range wordsSorted {
		if index >= amount {
			break
		}
		wordsTop = append(wordsTop, word.word)
	}

	return wordsTop
}

func Top10(input string) []string {
	wordsValid := makeWordsValid(input)
	wordsCounted := makeWordsCounted(wordsValid)
	wordsSorted := makeWordsSorted(wordsCounted)
	wordsTop10 := makeWordsTop(wordsSorted, 10)

	return wordsTop10
}
