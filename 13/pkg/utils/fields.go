package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseFields парсит строку вида "1,3,5-7,10" и возвращает map с номерами полей.
// Результат — это множество индексов (1-based), которые нужно обработать.
func ParseFields(fields string) (map[int]struct{}, error) {
	result := make(map[int]struct{})
	parts := strings.Split(strings.TrimSpace(fields), ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			startStr, endStr := strings.TrimSpace(rangeParts[0]), strings.TrimSpace(rangeParts[1])
			start, err1 := strconv.Atoi(startStr)
			end, err2 := strconv.Atoi(endStr)

			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("parse range '%s': both parts must be valid integers", part)
			}

			if start <= 0 || end <= 0 {
				return nil, fmt.Errorf("field numbers must be positive, got range: %d-%d", start, end)
			}

			if start > end {
				return nil, fmt.Errorf("invalid range: %d-%d (start greater than end)", start, end)
			}

			for i := start; i <= end; i++ {
				result[i] = struct{}{}
			}
		} else {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("parse field number: '%s' is not a valid integer", part)
			}
			if num <= 0 {
				return nil, fmt.Errorf("number must be positive: %d", num)
			}
			result[num] = struct{}{}
		}
	}

	return result, nil
}
