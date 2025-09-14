package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
	"l2.12/pkg/utils"
)

var (
	// flag -A N — после каждой найденной строки дополнительно вывести N строк после неё (контекст).
	afterContext int

	// flag -B N — вывести N строк до каждой найденной строки.
	beforeContext int

	// flag -C N — вывести N строк контекста вокруг найденной строки (включает и до, и после; эквивалентно -A N -B N).
	context int

	// flag -c — выводить только то количество строк, что совпадающих с шаблоном (т.е. вместо самих строк — число).
	onlyNumLines bool

	// flag -i — игнорировать регистр.
	caseInsensitive bool

	// flag -v — инвертировать фильтр: выводить строки, не содержащие шаблон.
	invertTemplate bool

	// -F — воспринимать шаблон как фиксированную строку, а не регулярное выражение (т.е. выполнять точное совпадение подстроки).
	fixedTemplate bool

	// flag -n — выводить номер строки перед каждой найденной строкой.
	showLineNumbers bool
)

func init() {
	pflag.IntVarP(&afterContext, "after-context", "A", 0, "print N lines of trailing context after matching lines")
	pflag.IntVarP(&beforeContext, "before-context", "B", 0, "print N lines of leading context before matching lines")
	pflag.IntVarP(&context, "context", "C", 0, "print N lines of output context (equivalent to -A N -B N)")
	pflag.BoolVarP(&onlyNumLines, "num-lines", "c", false, "print only amount of lines")
	pflag.BoolVarP(&caseInsensitive, "ignore-case", "i", false, "ignore case distinctions in patterns")
	pflag.BoolVarP(&invertTemplate, "invert", "v", false, "print lines that not match template")
	pflag.BoolVarP(&fixedTemplate, "fixed-template", "F", false, "template like fixed string (not regexp)")
	pflag.BoolVarP(&showLineNumbers, "line-number", "n", false, "print line numbers")

	pflag.CommandLine.SetNormalizeFunc(pflag.CommandLine.GetNormalizeFunc())
}

func main() {
	pflag.Parse()

	args := pflag.Args()

	if len(args) == 0 {
		log.Fatal("usage: grep [OPTIONS] PATTERN [FILE]")
	}

	pattern := args[0]
	var filePath string
	if len(args) > 1 {
		filePath = args[1]
	}

	if context > 0 {
		if afterContext == 0 {
			afterContext = context
		}
		if beforeContext == 0 {
			beforeContext = context
		}
	}

	matcher := utils.CompileMatcher(pattern, caseInsensitive, fixedTemplate)

	input, err := openInput(filePath)
	if err != nil {
		log.Fatalf("failed to open input: %v", err)
	}
	defer input.Close()

	lines, err := readLines(input)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	matches := utils.FindMatchesIndexes(lines, matcher, invertTemplate)

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if onlyNumLines {
		utils.PrintOnlyNumLines(writer, matches)
	} else {
		utils.PrintResults(writer, lines, matches, beforeContext, afterContext, showLineNumbers)
	}
}

func openInput(filePath string) (io.ReadCloser, error) {
	if filePath == "" {
		return io.NopCloser(os.Stdin), nil
	}
	return os.Open(filePath)
}

func readLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
