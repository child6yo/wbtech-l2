package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"l2.10/pkg/arguments"
)

var (
	// flagF (-f/--flag) - указывает путь к файлу, из которого нужно читать (напр. -f=/example/test.txt)
	flagF string

	// flagK (-k) - сортировать по колонке строки (разделение - символ табуляции (\t)), необходимо указать колонку (напр. -k=3)
	flagK int64

	// flagN (-n) - сортировать по числовому значению
	flagN bool

	// flagR (-r) - сортировать в обратном порядке
	flagR bool

	// flagU (-u) - не выводить повторяющиеся строки
	flagU bool
)

func openInput() (io.ReadCloser, error) {
	// чтение из stdin
	if flagF == "" {
		return io.NopCloser(os.Stdin), nil
	}
	// чтение из файла
	file, err := os.Open(flagF)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return file, nil
}

func scanInput(rc io.ReadCloser) (*[]string, error) {
	defer rc.Close()

	reader := bufio.NewReader(rc)
	var lines []string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), "\\t", "\t")
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read input: %w", err)
	}

	return &lines, nil
}

func writeOutput(w *bufio.Writer, output []string) {
	if flagU {
		arguments.PrintUniqueOnly(w, output)
		return
	}

	for _, line := range output {
		fmt.Fprintln(w, line)
	}
}

func main() {
	pflag.StringVarP(&flagF, "file", "f", "", "defines file path")
	pflag.Int64VarP(&flagK, "k", "k", 0, "sort by column")
	pflag.BoolVarP(&flagN, "n", "n", false, "sort by number value")
	pflag.BoolVarP(&flagR, "r", "r", false, "reverse sort")
	pflag.BoolVarP(&flagU, "u", "u", false, "only unique strings")
	pflag.Parse()

	reader, err := openInput()
	if err != nil {
		log.Fatal(err)
	}

	lines, err := scanInput(reader)
	if err != nil {
		log.Fatal(err)
	}

	if flagK != 0 {
		arguments.ColSort(lines, int(flagK)-1, flagN)
	} else if flagN {
		arguments.NumSort(lines)
	} else {
		arguments.BasicSort(lines)
	}

	if flagR {
		arguments.Reverse(*lines)
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	writeOutput(writer, *lines)
}
