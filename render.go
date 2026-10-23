package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type navItem struct {
	Name    string
	URL     string
	Class   string
	Current string
}

// Card is the data needed to render a post in a section overview.
type Card struct {
	Title     string
	URL       string
	Summary   string
	DateISO   string
	DateHuman string
	ImageURL  string
	ImageAlt  string
	HasImage  bool
	Tags      []Link
}

type view struct {
	Site     SiteConfig
	Title    string
	Nav      []navItem
	Year     int
	Page     *Page
	Body     template.HTML
	Cards    []Card
	Heading  string
	Links    []Link
	TagLinks []Link
}

type renderer struct {
	dir   string
	cache map[string]*template.Template
}

func newRenderer(dir string) *renderer {
	return &renderer{dir: dir, cache: map[string]*template.Template{}}
}

func (r *renderer) render(kind string, data view) ([]byte, error) {
	tpl := r.cache[kind]
	if tpl == nil {
		var err error
		tpl, err = template.New("base.html").ParseFiles(
			filepath.Join(r.dir, "base.html"),
			filepath.Join(r.dir, kind+".html"),
		)
		if err != nil {
			return nil, err
		}
		r.cache[kind] = tpl
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "base.html", data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type siteBuilder struct {
	cfg     SiteConfig
	content *Content
	out     string
	tpl     *renderer
	now     time.Time
	written int
}

func newSiteBuilder(cfg SiteConfig, content *Content, out string) *siteBuilder {
	return &siteBuilder{
		cfg:     cfg,
		content: content,
		out:     out,
		tpl:     newRenderer("templates"),
		now:     time.Now(),
	}
}

func (b *siteBuilder) Run() (int, error) {
	if b.out == "" || b.out == "." || b.out == "/" {
		return 0, fmt.Errorf("refusing to clean unsafe output directory %q", b.out)
	}
	if err := os.RemoveAll(b.out); err != nil {
		return 0, err
	}
	if err := os.MkdirAll(b.out, 0o755); err != nil {
		return 0, err
	}

	if err := copyTree("static", b.out, false); err != nil {
		return 0, err
	}
	if err := copyTree("content", b.out, true); err != nil {
		return 0, err
	}

	home := b.content.Home
	if err := b.writePage("/", "home", view{
		Title: b.cfg.Title,
		Page:  home,
		Body:  home.Body,
	}); err != nil {
		return 0, err
	}

	for _, section := range b.content.Sections {
		if err := b.writePage(section.Route, "section", view{
			Title:    b.titleFor(section),
			Page:     section,
			Body:     section.Body,
			Cards:    b.cardsFor(section.Route),
			TagLinks: section.TagLinks,
		}); err != nil {
			return 0, err
		}
	}

	for _, page := range b.content.Pages {
		if err := b.writePage(page.Route, "page", view{
			Title:    b.titleFor(page),
			Page:     page,
			Body:     page.Body,
			TagLinks: page.TagLinks,
		}); err != nil {
			return 0, err
		}
	}

	for _, taxonomy := range b.content.Taxonomies {
		terms := taxonomy.SortedTerms()
		links := make([]Link, 0, len(terms))
		for _, term := range terms {
			links = append(links, Link{Name: term.Name, URL: term.URL})
		}
		if err := b.writePage(taxonomy.Route, "list", view{
			Title:   taxonomy.Name + " | " + b.cfg.Title,
			Heading: taxonomy.Name,
			Links:   links,
		}); err != nil {
			return 0, err
		}

		for _, term := range terms {
			if err := b.writePage(term.URL, "list", view{
				Title:   term.Name + " | " + b.cfg.Title,
				Heading: term.Name,
				Links:   pageLinks(term.Pages),
			}); err != nil {
				return 0, err
			}
		}
	}

	if err := b.writeSitemap(); err != nil {
		return 0, err
	}
	if err := b.writeFeeds(); err != nil {
		return 0, err
	}

	return b.written, nil
}

func (b *siteBuilder) writePage(urlPath, kind string, data view) error {
	data.Site = b.cfg
	data.Nav = b.nav(urlPath)
	data.Year = b.now.Year()

	rendered, err := b.tpl.render(kind, data)
	if err != nil {
		return fmt.Errorf("render %s: %w", urlPath, err)
	}
	if err := writeFile(outPath(b.out, urlPath), rendered); err != nil {
		return err
	}
	b.written++
	return nil
}

func (b *siteBuilder) titleFor(page *Page) string {
	if page.IsHome {
		return b.cfg.Title
	}
	return page.Title + " | " + b.cfg.Title
}

func (b *siteBuilder) nav(route string) []navItem {
	items := make([]navItem, 0, len(b.cfg.Menu))
	for _, item := range b.cfg.Menu {
		nav := navItem{Name: item.Name, URL: item.URL}
		switch {
		case route == item.URL:
			nav.Class = "active"
			nav.Current = "page"
		case item.URL != "/" && strings.HasPrefix(route, item.URL):
			nav.Class = "ancestor"
			nav.Current = "true"
		}
		items = append(items, nav)
	}
	return items
}

func (b *siteBuilder) cardsFor(route string) []Card {
	pages := b.pagesUnder(route)
	sortPages(pages)

	cards := make([]Card, 0, len(pages))
	for _, page := range pages {
		card := Card{
			Title:     page.Title,
			URL:       page.Route,
			Summary:   page.Summary,
			DateISO:   page.DateISO,
			DateHuman: page.DateHuman,
			Tags:      page.TagLinks,
		}
		if page.Image != nil && page.Image.Src != "" {
			card.HasImage = true
			imageSource := page.Image.Src
			if page.Image.Thumb != "" {
				imageSource = page.Image.Thumb
			}
			card.ImageURL = resolveAsset(page.Route, imageSource)
			card.ImageAlt = page.Image.Alt
			if card.ImageAlt == "" {
				card.ImageAlt = page.Title
			}
		} else {
			card.ImageURL = "/blog-post-no-image.jpg"
			card.ImageAlt = page.Title + " has no image"
		}
		cards = append(cards, card)
	}
	return cards
}

func (b *siteBuilder) pagesUnder(route string) []*Page {
	var pages []*Page
	for _, page := range b.content.Pages {
		if strings.HasPrefix(page.Route, route) {
			pages = append(pages, page)
		}
	}
	return pages
}

func resolveAsset(route, source string) string {
	if strings.HasPrefix(source, "/") || strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		return source
	}
	return strings.TrimSuffix(route, "/") + "/" + source
}

func pageLinks(pages []*Page) []Link {
	links := make([]Link, 0, len(pages))
	for _, page := range pages {
		links = append(links, Link{Name: page.Title, URL: page.Route})
	}
	return links
}

func (b *siteBuilder) absURL(route string) string {
	return strings.TrimRight(b.cfg.BaseURL, "/") + route
}

func (b *siteBuilder) writeSitemap() error {
	urls := make([]sitemapURL, 0)

	home := sitemapURL{Loc: b.absURL("/")}
	if b.content.Home.HasDate {
		home.LastMod = b.content.Home.DateISO
	}
	urls = append(urls, home)

	for _, section := range b.content.Sections {
		urls = append(urls, sitemapURL{Loc: b.absURL(section.Route), LastMod: section.DateISO})
	}
	for _, page := range b.content.Pages {
		urls = append(urls, sitemapURL{Loc: b.absURL(page.Route), LastMod: page.DateISO})
	}
	for _, taxonomy := range b.content.Taxonomies {
		urls = append(urls, sitemapURL{Loc: b.absURL(taxonomy.Route)})
		for _, term := range taxonomy.SortedTerms() {
			lastMod := ""
			if latest := latestDate(term.Pages); !latest.IsZero() {
				lastMod = latest.Format(time.RFC3339)
			}
			urls = append(urls, sitemapURL{Loc: b.absURL(term.URL), LastMod: lastMod})
		}
	}

	sort.Slice(urls, func(i, j int) bool {
		return urls[i].Loc < urls[j].Loc
	})
	return b.writeXML("/sitemap.xml", sitemapDoc{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	})
}

func (b *siteBuilder) writeFeeds() error {
	if err := b.writeFeed("/index.xml", b.cfg.Title, b.content.Pages); err != nil {
		return err
	}
	for _, section := range b.content.Sections {
		title := section.Title + " on " + b.cfg.Title
		if err := b.writeFeed(feedPath(section.Route), title, b.pagesUnder(section.Route)); err != nil {
			return err
		}
	}
	for _, taxonomy := range b.content.Taxonomies {
		for _, term := range taxonomy.SortedTerms() {
			title := term.Name + " on " + b.cfg.Title
			if err := b.writeFeed(feedPath(term.URL), title, term.Pages); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *siteBuilder) writeFeed(urlPath, title string, pages []*Page) error {
	sortPages(pages)

	link := strings.TrimSuffix(urlPath, "index.xml")
	if link == "" {
		link = "/"
	}

	items := make([]rssItem, 0, len(pages))
	for _, page := range pages {
		item := rssItem{
			Title:       page.Title,
			Link:        b.absURL(page.Route),
			GUID:        b.absURL(page.Route),
			Description: page.Summary,
		}
		if page.HasDate {
			item.PubDate = page.Date.Format(time.RFC1123Z)
		}
		items = append(items, item)
	}

	lastBuild := b.now
	if latest := latestDate(pages); !latest.IsZero() {
		lastBuild = latest
	}

	return b.writeXML(urlPath, rssDoc{
		Version: "2.0",
		Channel: rssChannel{
			Title:         title,
			Link:          b.absURL(link),
			Description:   "Recent content in " + title,
			Language:      b.cfg.LanguageCode,
			LastBuildDate: lastBuild.Format(time.RFC1123Z),
			Items:         items,
		},
	})
}

func (b *siteBuilder) writeXML(urlPath string, value any) error {
	data, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append([]byte(xml.Header), data...)
	return writeFile(outPath(b.out, urlPath), data)
}

func feedPath(route string) string {
	if route == "/" {
		return "/index.xml"
	}
	return strings.TrimSuffix(route, "/") + "/index.xml"
}

func latestDate(pages []*Page) time.Time {
	var latest time.Time
	for _, page := range pages {
		if page.HasDate && page.Date.After(latest) {
			latest = page.Date
		}
	}
	return latest
}

func outPath(root, urlPath string) string {
	target := strings.TrimPrefix(urlPath, "/")
	if target == "" || strings.HasSuffix(target, "/") {
		target += "index.html"
	}
	return filepath.Join(root, filepath.FromSlash(target))
}

func writeFile(file string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o644)
}

func copyTree(source, destination string, skipMarkdown bool) error {
	if _, err := os.Stat(source); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return filepath.WalkDir(source, func(file string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, file)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if skipMarkdown && strings.EqualFold(filepath.Ext(file), ".md") {
			return nil
		}
		return copyFile(file, target)
	})
}

func copyFile(source, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}

	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

type sitemapDoc struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type rssDoc struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language,omitempty"`
	LastBuildDate string    `xml:"lastBuildDate,omitempty"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate,omitempty"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
}
