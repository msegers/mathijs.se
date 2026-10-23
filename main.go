package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// SiteConfig is loaded from site.json.
type SiteConfig struct {
	Title        string     `json:"title"`
	BaseURL      string     `json:"baseURL"`
	LanguageCode string     `json:"languageCode"`
	Menu         []MenuItem `json:"menu"`
}

type MenuItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	out := flag.String("out", "public", "directory to write the generated site to")
	serve := flag.Bool("serve", false, "serve the generated site after building")
	addr := flag.String("addr", ":1313", "address to listen on with --serve")
	flag.Parse()

	cfg, err := loadSite("site.json")
	if err != nil {
		log.Fatal(err)
	}

	content, err := loadContent("content")
	if err != nil {
		log.Fatal(err)
	}

	builder := newSiteBuilder(cfg, content, *out)
	pages, err := builder.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Built %d pages into %s\n", pages, *out)

	if *serve {
		host := *addr
		if strings.HasPrefix(host, ":") {
			host = "localhost" + host
		}
		fmt.Printf("Serving %s at http://%s/\n", *out, host)
		server := &http.Server{
			Addr:              *addr,
			Handler:           http.FileServer(http.Dir(*out)),
			ReadHeaderTimeout: 5 * time.Second,
		}
		log.Fatal(server.ListenAndServe())
	}
}

func loadSite(filename string) (SiteConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return SiteConfig{}, err
	}

	var cfg SiteConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return SiteConfig{}, fmt.Errorf("parse %s: %w", filename, err)
	}
	if cfg.Title == "" {
		return SiteConfig{}, fmt.Errorf("%s: title is required", filename)
	}
	if cfg.BaseURL == "" {
		return SiteConfig{}, fmt.Errorf("%s: baseURL is required", filename)
	}
	if cfg.LanguageCode == "" {
		cfg.LanguageCode = "en-us"
	}
	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL += "/"
	}
	return cfg, nil
}
