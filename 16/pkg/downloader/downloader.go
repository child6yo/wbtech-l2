package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Downloader — структура для управления загрузкой.
type Downloader struct {
	userAgent string
	outputDir string
	maxDepth  int
	visited   map[string]bool
	mu        sync.RWMutex
	client    *http.Client
}

// NewDownloader — создает новый экземпляр загрузчика.
func NewDownloader(userAgent, outputDir string, maxDepth int) *Downloader {
	return &Downloader{
		userAgent: userAgent,
		outputDir: outputDir,
		maxDepth:  maxDepth,
		visited:   make(map[string]bool),
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// Download — точка входа для загрузки URL.
func (d *Downloader) Download(startURL string) error {
	return d.downloadRecursive(startURL, 0)
}

// downloadRecursive — рекурсивная загрузка с ограничением по глубине.
func (d *Downloader) downloadRecursive(rawURL string, depth int) error {
	if depth > d.maxDepth {
		return nil
	}

	u, err := normalizeURL(rawURL)
	if err != nil {
		return fmt.Errorf("failed to normalize URL %s: %w", rawURL, err)
	}

	d.mu.RLock()
	if d.visited[u.String()] {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	d.visited[u.String()] = true
	d.mu.Unlock()

	resp, err := d.fetch(u.String())
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", u.String(), err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	isHTML := strings.Contains(contentType, "text/html")

	localPath, err := d.getLocalPath(u, isHTML)
	if err != nil {
		return fmt.Errorf("failed to get local path for %s: %w", u.String(), err)
	}

	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create dir for %s: %w", localPath, err)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", localPath, err)
	}

	_, err = io.Copy(file, resp.Body)
	file.Close()
	if err != nil {
		return fmt.Errorf("failed to save content for %s: %w", localPath, err)
	}

	fmt.Printf("📥 Saved: %s -> %s\n", u.String(), localPath)

	if isHTML {
		err = d.rewriteLinksInHTML(localPath, u)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to rewrite links in %s: %v\n", localPath, err)
		}
	}

	if isHTML && depth < d.maxDepth {
		links, err := extractLinks(localPath, u)
		if err != nil {
			return fmt.Errorf("failed to extract links from %s: %w", localPath, err)
		}

		var wg sync.WaitGroup
		for _, link := range links {
			wg.Add(1)
			go func(l string) {
				defer wg.Done()
				err := d.downloadRecursive(l, depth+1)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to download %s: %v\n", l, err)
				}
			}(link)
		}
		wg.Wait()
	}

	return nil
}

// fetch — делает HTTP GET запрос.
func (d *Downloader) fetch(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", d.userAgent)

	return d.client.Do(req)
}

// getLocalPath — возвращает локальный путь для сохранения файла.
func (d *Downloader) getLocalPath(u *url.URL, isHTML bool) (string, error) {
	path := u.Path
	if path == "" || path == "/" {
		path = "/index.html"
	}

	if isHTML && filepath.Ext(path) == "" {
		path += ".html"
	}

	cleanPath := filepath.Clean(path)
	if cleanPath[0] == '/' {
		cleanPath = cleanPath[1:]
	}

	return filepath.Join(d.outputDir, u.Host, cleanPath), nil
}

// normalizeURL — нормализует URL, приводя его к стандартному виду.
func normalizeURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	if u.Host == "" {
		return nil, fmt.Errorf("missing host in URL: %s", raw)
	}

	u.Fragment = ""
	u.RawQuery = ""

	return u, nil
}

// extractLinks — извлекает ссылки из HTML-файла.
func extractLinks(filePath string, baseURL *url.URL) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return nil, err
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "img" || n.Data == "link" || n.Data == "script") {
			for _, attr := range n.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					link := attr.Val
					resolved, err := resolveLink(link, baseURL)
					if err == nil {
						links = append(links, resolved)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	unique := make(map[string]bool)
	var result []string
	for _, l := range links {
		if !unique[l] {
			unique[l] = true
			result = append(result, l)
		}
	}

	return result, nil
}

// resolveLink — разрешает относительную ссылку относительно базового URL.
func resolveLink(link string, baseURL *url.URL) (string, error) {
	if strings.HasPrefix(link, "#") || strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") || strings.HasPrefix(link, "javascript:") {
		return "", nil
	}

	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return link, nil
	}

	u, err := baseURL.Parse(link)
	if err != nil {
		return "", err
	}

	if u.Host != baseURL.Host {
		return "", nil
	}

	return u.String(), nil
}

// rewriteLinksInHTML — перезаписывает ссылки в HTML-файле на локальные пути.
func (d *Downloader) rewriteLinksInHTML(filePath string, baseURL *url.URL) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "img" || n.Data == "link" || n.Data == "script" || n.Data == "form") {
			for i, attr := range n.Attr {
				key := attr.Key
				val := attr.Val

				if key == "href" || key == "src" || key == "action" {
					resolved, err := resolveLink(val, baseURL)
					if err != nil || resolved == "" {
						continue
					}

					u, _ := url.Parse(resolved)
					localPath, err := d.getLocalPath(u, false) 
					if err != nil {
						continue
					}

					if _, err := os.Stat(localPath); err == nil {
						relPath, err := filepath.Rel(filepath.Dir(filePath), localPath)
						if err != nil {
							continue
						}
						relPath = filepath.ToSlash(relPath)
						n.Attr[i].Val = relPath
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return html.Render(outputFile, doc)
}
