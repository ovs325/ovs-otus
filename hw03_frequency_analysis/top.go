package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	type Word struct {
		Word      string
		Frequency int
	}
	textMap := make(map[string]int, 0)
	for _, word := range strings.Fields(text) {
		if _, ok := textMap[word]; !ok {
			textMap[word] = 1
		} else {
			textMap[word]++
		}
	}
	words := make([]Word, 0)
	for key, val := range textMap {
		words = append(words, Word{key, val})
	}
	sort.Slice(words, func(i, j int) bool {
		if words[i].Frequency == words[j].Frequency {
			return words[i].Word < words[j].Word
		}
		return words[i].Frequency > words[j].Frequency
	})
	res := make([]string, 0)
	for _, word := range words {
		res = append(res, word.Word)
	}
	if len(res) > 10 {
		return res[:10]
	}
	return res
}
