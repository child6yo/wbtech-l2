package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"l2.16/pkg/downloader"

	"github.com/spf13/pflag"
)

var (
	// flag -r "recursion" — глубина рекурсии (по умолчанию 0 — только корень)
	recursion int

	// flag -o "output" — путь к выходной директории
	outputDir string

	// flag -u "user-agent" — пользовательский агент (опционально)
	userAgent string
)

func init() {
	pflag.IntVarP(&recursion, "recursive", "r", 0, "recursion depth (0 = no recursion)")
	pflag.StringVarP(&outputDir, "output", "o", "./wget_output", "output directory for downloaded files")
	pflag.StringVarP(&userAgent, "user-agent", "u", "wget-go/1.0", "custom User-Agent header")
	pflag.CommandLine.SetNormalizeFunc(pflag.CommandLine.GetNormalizeFunc())
}

func main() {
	pflag.Parse()
	args := pflag.Args()

	if len(args) != 1 {
		log.Fatal("usage: wget [OPTIONS] <URL>")
	}

	urlStr := args[0]

	if _, err := url.ParseRequestURI(urlStr); err != nil {
		log.Fatalf("invalid URL: %v", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}

	dl := downloader.NewDownloader(userAgent, outputDir, recursion)

	err := dl.Download(urlStr)
	if err != nil {
		log.Fatalf("download failed: %v", err)
	}

	fmt.Printf("Download completed. Files saved in: %s\n", outputDir)
}
