package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordsCount struct {
	word  string
	count int
}

func Top10(input string) []string {
	input = strings.ToLower(input)
	wordsRaw := strings.Fields(input)
	wordsCounted := make(map[string]int)
	for _, word := range wordsRaw {
		word = strings.Trim(word, ",.!?\"'")
		if len(word) == 0 || word == "-" {
			continue
		}
		wordsCounted[word]++
	}

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

	wordsTop := make([]string, 0, 9)
	for index, word := range wordsSorted {
		if index > 9 {
			break
		}
		wordsTop = append(wordsTop, word.word)
	}

	return wordsTop
}
