package main

import (
	"fmt"
	"slices"
	"strings"
)

func findAnnagrams(strs []string) map[string][]string {
	if len(strs) < 1 {
		return nil
	}

	// основная мапа
	res := make(map[string][]string)
	
	// вспомогательное множество, позволяющее хранить информацию о ключах
	keys := make(map[string]string) 

	// еще одно вспомогательное множество, позволяющее сохранять только уникальные значения
	set := make(map[string]struct{}) 

	for _, s := range strs {
		if _, ok := set[s]; ok {
			continue
		}
		set[s] = struct{}{}

		s = strings.ToLower(s)
		r := []rune(s)
		slices.Sort(r) // сортирует по числовому значению рун, но в данном случае это значения не имеет
		if _, ok := keys[string(r)]; !ok {
			keys[string(r)] = s
			res[s] = []string{s}
			continue
		}
		res[keys[string(r)]] = append(res[keys[string(r)]], s)
	}

	// удаление ключей, где меньше двух значений
	for k, v := range res {
		if len(v) < 2 {
			delete(res, k)
		}
	}

	return res
}

func main() {
	fmt.Println(findAnnagrams([]string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол", "стол"}))
}
