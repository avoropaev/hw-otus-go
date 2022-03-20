package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCounter struct {
	Word  string
	Count int
}

func Top10(s string) []string {
	words := strings.Fields(s)
	counterMap := make(map[string]int)

	for _, word := range words {
		if _, ok := counterMap[word]; !ok {
			counterMap[word] = 0
		}

		counterMap[word]++
	}

	counterSlice := make([]wordCounter, 0, len(counterMap))

	for word, count := range counterMap {
		counterSlice = append(counterSlice, wordCounter{word, count})
	}

	sort.SliceStable(counterSlice, func(i, j int) bool {
		return counterSlice[i].Count > counterSlice[j].Count ||
			(counterSlice[i].Count == counterSlice[j].Count && counterSlice[i].Word < counterSlice[j].Word)
	})

	resultCap := 10
	if len(counterSlice) < 10 {
		resultCap = len(counterSlice)
	}

	result := make([]string, 0, resultCap)

	for _, wordCounter := range counterSlice[:resultCap] {
		result = append(result, wordCounter.Word)
	}

	return result
}
