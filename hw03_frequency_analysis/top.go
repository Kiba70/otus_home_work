package hw03frequencyanalysis

import (
	"cmp"
	"slices"
	"strings"
	"unicode"
)

const simpleSort = false

type sliceWordsType struct {
	count int
	word  string
}

func Top10(text string) []string {
	// Граничные значения
	if len(text) == 0 {
		return nil
	}

	// Считаем количество повторений слов по всему массиву
	words := make(map[string]int)

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '-'
	}

	for _, word := range strings.Fields(text) {
		if word == "-" {
			continue // "-" не слово
		}
		if s := strings.ToLower(strings.TrimFunc(word, f)); s != "" { // После truncate бывает ""
			words[s]++
		}
	}

	// Переносим в слайс и сортируем
	sliceWords := make([]sliceWordsType, 0, len(words))

	for s, n := range words {
		sliceWords = append(sliceWords, sliceWordsType{
			count: n,
			word:  s,
		})
	}

	lenResult := min(10, len(sliceWords)) // вдруг меньше 10?
	result := make([]string, lenResult)

	if simpleSort {
		// Вариант простой, но не очень эффективный т.к. производится лексикографическая сортировка всего слайса
		// Всё зависит от объёма текста
		slices.SortFunc(sliceWords, func(a, b sliceWordsType) int {
			//nolint: lll
			if n := cmp.Compare(b.count, a.count); n != 0 { // Сначала сравниваем количество совпадений. Задом на перёд - от большего к меньшему
				return n
			}
			return strings.Compare(a.word, b.word) // Затем лексикографически
		})
	} else {
		// Более замороченный вариант
		// На самом деле тесты показывают одинаковое время выполнения
		slices.SortFunc(sliceWords, func(a, b sliceWordsType) int {
			return cmp.Compare(b.count, a.count)
		})
		idx := 0
		for i := max(lenResult, 1); i < len(sliceWords); i++ {
			if sliceWords[i-1].count != sliceWords[i].count {
				idx = i
				break
			}
		}
		slices.SortFunc(sliceWords[:idx], func(a, b sliceWordsType) int {
			//nolint: lll
			if n := cmp.Compare(b.count, a.count); n != 0 { // Сначала сравниваем количество совпадений. Задом на перёд - от большего к меньшему
				return n
			}
			return strings.Compare(a.word, b.word) // Затем лексикографически
		})
	}

	for i := 0; i < lenResult; i++ {
		result[i] = sliceWords[i].word
	}

	return result
}
