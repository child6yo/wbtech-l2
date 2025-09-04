package utils

import (
	"log"
	"regexp"
	"strings"
)

// Matcher - функция, валидирующая строку по шаблону.
type Matcher func(string) bool

// CompileMatcher - фабрика, возвращающая Matcher соответствующий параметрам.
func CompileMatcher(pattern string, caseInsensitive, fixedTemplate bool) Matcher {
	if fixedTemplate {
		if caseInsensitive {
			pattern = strings.ToLower(pattern)
			return func(line string) bool {
				return strings.Contains(strings.ToLower(line), pattern)
			}
		} else {
			return func(line string) bool {
				return strings.Contains(line, pattern)
			}
		}
	}

	var re *regexp.Regexp
	var err error
	if caseInsensitive {
		re, err = regexp.Compile("(?i)" + pattern)
	} else {
		re, err = regexp.Compile(pattern)
	}
	if err != nil {
		log.Fatalf("invalid regex pattern: %v", err)
	}

	return re.MatchString
}