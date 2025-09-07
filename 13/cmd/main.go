package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
	"l2.13/pkg/utils"
)

var (
	// flag -f "fields" — указание номеров полей (колонок), которые нужно вывести. Номера через запятую, можно диапазоны.
	fields string

	// flag -d "delimiter" — использовать другой разделитель (символ). По умолчанию разделитель — табуляция ('\t').
	delimiter string

	// flag -s – (separated) только строки, содержащие разделитель. Если флаг указан, то строки без разделителя игнорируются (не выводятся).
	separated bool
)

func init() {
	pflag.StringVarP(&fields, "fields", "f", "", "print only n (or from n to m) columns. usage «-f 1,3-5»")
	pflag.StringVarP(&delimiter, "delimiter", "d", "\t", `column delimiter (default '\t')`)
	pflag.BoolVarP(&separated, "separated", "s", false, "print only strings with delimiter")
	pflag.CommandLine.SetNormalizeFunc(pflag.CommandLine.GetNormalizeFunc())
}

func main() {
	pflag.Parse()
	args := pflag.Args()

	if len(args) > 1 {
		log.Fatal("usage: cut [OPTIONS] [FILE]")
	}

	var filePath string
	if len(args) >= 1 {
		filePath = args[0]
	}

	input, err := openInput(filePath)
	if err != nil {
		log.Fatalf("failed to open input: %v", err)
	}
	defer input.Close()

	var idxs map[int]struct{}
	if fields != "" {
		idxs, err = utils.ParseFields(fields)
		if err != nil {
			log.Fatalf("failed to parse fields: %v", err)
		}
	}

	lines, err := readLines(input, delimiter, idxs, separated)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	writeOutput(writer, lines)
}

func openInput(filePath string) (io.ReadCloser, error) {
	if filePath == "" {
		return io.NopCloser(os.Stdin), nil
	}
	return os.Open(filePath)
}

func readLines(r io.Reader, delimiter string, idxs map[int]struct{}, separated bool) ([]string, error) {
	var result []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		res, ok := utils.HandleString(line, delimiter, idxs, separated)
		if !ok {
			continue
		}
		result = append(result, res)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return result, nil
}

func writeOutput(w *bufio.Writer, output []string) {
	for _, line := range output {
		fmt.Fprintln(w, line)
	}
}
