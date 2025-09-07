package utils

import "strings"

// HandleString разбивает строку s по разделителю delimiter и оставляет только те поля,
// номера которых указаны в fields (нумерация с 1). Если separated == true, то ожидается,
// что строка содержит как минимум два поля; в противном случае возвращается ( "", false ).
// Возвращает результирующую строку и флаг успеха.
func HandleString(s string, delimiter string, fields map[int]struct{}, separated bool) (string, bool) {
	parts := strings.Split(s, delimiter)
	if separated && len(parts) < 2 {
		return "", false
	}

	res := handleParts(parts, delimiter, fields)

	return res, true
}

func handleParts(parts []string, delimiter string, fields map[int]struct{}) string {
	var selectedParts []string

	if fields == nil {
		return strings.Join(selectedParts, delimiter)
	}

	for idx, field := range parts {
		if _, ok := fields[idx+1]; ok {
			selectedParts = append(selectedParts, field)
		}
	}

	return strings.Join(selectedParts, delimiter)
}
