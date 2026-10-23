package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	ghtml "github.com/yuin/goldmark/renderer/html"
)

// Image is the optional image object in a page's front matter.
type Image struct {
	Src   string `toml:"src"`
	Alt   string `toml:"alt"`
	Thumb string `toml:"thumb"`
}

// dateValue accepts both quoted dates and TOML datetimes.
type dateValue string

func (d *dateValue) UnmarshalTOML(value any) error {
	switch typed := value.(type) {
	case time.Time:
		*d = dateValue(typed.Format(time.RFC3339))
	case string:
		*d = dateValue(typed)
	default:
		return fmt.Errorf("date must be a string or TOML datetime")
	}
	return nil
}

type frontMatter struct {
	Title      string    `toml:"title"`
	Date       dateValue `toml:"date"`
	Draft      bool      `toml:"draft"`
	Tags       []string  `toml:"tags"`
	Categories []string  `toml:"categories"`
	Summary    string    `toml:"summary"`
	Layout     string    `toml:"layout"`
	Image      *Image    `toml:"image"`
}

// Link is a small rendered link used by the templates.
type Link struct {
	Name string
	URL  string
}

// Page is a content file that will be turned into one HTML page.
type Page struct {
	Title      string
	Route      string
	Date       time.Time
	HasDate    bool
	Draft      bool
	Tags       []string
	Categories []string
	Summary    string
	Layout     string
	Image      *Image
	IsSection  bool
	IsHome     bool
	Source     string
	Body       template.HTML
	DateISO    string
	DateHuman  string
	TagLinks   []Link
}

// Term is one value of a taxonomy, such as the "linux" tag.
type Term struct {
	Slug  string
	Name  string
	URL   string
	Pages []*Page
}

// Taxonomy is a group of terms, such as all tags.
type Taxonomy struct {
	Name  string
	Route string
	Terms map[string]*Term
}

func (t *Taxonomy) SortedTerms() []*Term {
	terms := make([]*Term, 0, len(t.Terms))
	for _, term := range t.Terms {
		terms = append(terms, term)
	}
	sort.Slice(terms, func(i, j int) bool {
		return terms[i].Name < terms[j].Name
	})
	return terms
}

// Content holds the parsed site content.
type Content struct {
	Home       *Page
	Sections   []*Page
	Pages      []*Page
	Taxonomies []*Taxonomy
}

var (
	externalLinkPattern = regexp.MustCompile(`<a href="(https?://[^"]*)"`)
	htmlTagPattern      = regexp.MustCompile(`<[^>]+>`)
	whitespacePattern   = regexp.MustCompile(`\s+`)
)

func loadContent(root string) (*Content, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(ghtml.WithUnsafe()),
	)

	content := &Content{}
	err := filepath.WalkDir(root, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(file), ".md") {
			return nil
		}

		page, err := parsePage(root, file, md)
		if err != nil {
			return err
		}
		if page.Draft {
			return nil
		}

		switch {
		case page.IsHome:
			content.Home = page
		case page.IsSection:
			content.Sections = append(content.Sections, page)
		default:
			content.Pages = append(content.Pages, page)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if content.Home == nil {
		return nil, fmt.Errorf("content root is missing _index.md")
	}

	sortPages(content.Pages)
	sort.Slice(content.Sections, func(i, j int) bool {
		return content.Sections[i].Route < content.Sections[j].Route
	})
	content.buildTaxonomies()
	return content, nil
}

func parsePage(root, file string, md goldmark.Markdown) (*Page, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	front, body := splitFrontMatter(string(data))
	var meta frontMatter
	if strings.TrimSpace(front) != "" {
		if _, err := toml.Decode(front, &meta); err != nil {
			return nil, fmt.Errorf("parse front matter: %w", err)
		}
	}

	rel, err := filepath.Rel(root, file)
	if err != nil {
		return nil, err
	}
	rel = filepath.ToSlash(rel)
	dir := path.Dir(rel)
	if dir == "." {
		dir = ""
	}
	base := path.Base(rel)

	page := &Page{
		Title:      meta.Title,
		Draft:      meta.Draft,
		Tags:       meta.Tags,
		Categories: meta.Categories,
		Summary:    meta.Summary,
		Layout:     meta.Layout,
		Image:      meta.Image,
		Source:     file,
	}

	switch {
	case base == "_index.md" && dir == "":
		page.IsHome = true
		page.IsSection = true
		page.Route = "/"
	case base == "_index.md":
		page.IsSection = true
		page.Route = "/" + dir + "/"
	case base == "index.md":
		page.Route = "/" + dir + "/"
	default:
		page.Route = "/" + strings.TrimSuffix(rel, ".md") + "/"
	}

	if page.Title == "" {
		page.Title = humanTitle(base)
	}
	if meta.Date != "" {
		page.Date, page.HasDate = parseDate(string(meta.Date))
	}
	if page.HasDate {
		page.DateISO = page.Date.Format(time.RFC3339)
		page.DateHuman = page.Date.Format("January 2, 2006")
	}

	rendered, err := renderMarkdown(md, body)
	if err != nil {
		return nil, err
	}
	page.Body = rendered
	if page.Summary == "" {
		page.Summary = fallbackSummary(string(rendered))
	}
	return page, nil
}

func splitFrontMatter(source string) (string, string) {
	lines := strings.Split(source, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) != "+++" {
			continue
		}
		if start == -1 {
			start = i
			continue
		}
		return strings.Join(lines[start+1:i], "\n"), strings.Join(lines[i+1:], "\n")
	}
	return "", source
}

func parseDate(value string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05-0700",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, strings.TrimSpace(value)); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func renderMarkdown(md goldmark.Markdown, source string) (template.HTML, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	rendered := externalLinkPattern.ReplaceAllString(buf.String(), `<a href="$1" target="_blank"`)
	return template.HTML(rendered), nil
}

func fallbackSummary(rendered string) string {
	text := html.UnescapeString(htmlTagPattern.ReplaceAllString(rendered, " "))
	text = strings.TrimSpace(whitespacePattern.ReplaceAllString(text, " "))
	words := strings.Fields(text)
	if len(words) > 40 {
		return strings.Join(words[:40], " ") + "..."
	}
	return strings.Join(words, " ")
}

func humanTitle(name string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.TrimPrefix(name, "_")
	if name == "" || name == "index" {
		return ""
	}
	words := strings.Fields(strings.ReplaceAll(name, "-", " "))
	for i, word := range words {
		words[i] = titleFirst(word)
	}
	return strings.Join(words, " ")
}

func titleFirst(value string) string {
	if value == "" {
		return value
	}
	runes := []rune(value)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func slugify(value string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case r == ' ' || r == '-' || r == '_' || r == '/':
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func sortPages(pages []*Page) {
	sort.SliceStable(pages, func(i, j int) bool {
		a, b := pages[i], pages[j]
		if a.HasDate != b.HasDate {
			return a.HasDate
		}
		if !a.Date.Equal(b.Date) {
			return a.Date.After(b.Date)
		}
		return a.Title < b.Title
	})
}

func (c *Content) buildTaxonomies() {
	c.Taxonomies = []*Taxonomy{
		{Name: "Tags", Route: "/tags/", Terms: map[string]*Term{}},
		{Name: "Categories", Route: "/categories/", Terms: map[string]*Term{}},
	}

	tags := c.Taxonomies[0]
	categories := c.Taxonomies[1]

	for _, page := range c.Pages {
		addTerms(tags, page.Tags, page)
		addTerms(categories, page.Categories, page)
	}

	for _, taxonomy := range c.Taxonomies {
		for _, term := range taxonomy.Terms {
			sortPages(term.Pages)
		}
	}

	for _, page := range c.Pages {
		for _, tag := range page.Tags {
			term := tags.Terms[slugify(tag)]
			if term != nil {
				page.TagLinks = append(page.TagLinks, Link{Name: term.Name, URL: term.URL})
			}
		}
	}
}

func addTerms(taxonomy *Taxonomy, values []string, page *Page) {
	for _, value := range values {
		slug := slugify(value)
		if slug == "" {
			continue
		}
		term := taxonomy.Terms[slug]
		if term == nil {
			term = &Term{
				Slug: slug,
				Name: titleFirst(value),
				URL:  taxonomy.Route + slug + "/",
			}
			taxonomy.Terms[slug] = term
		}
		term.Pages = append(term.Pages, page)
	}
}
