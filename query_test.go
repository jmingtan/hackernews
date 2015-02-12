package main

import (
	"github.com/PuerkitoBio/goquery"
	"os"
	"testing"
)

func TestExtractArticles(t *testing.T) {
	file, err := os.Open("test_page.html")
	if err != nil {
		t.Error(err)
		return
	}
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Error(err)
	}
	articles := ExtractArticles(doc)
	first_article := articles[0]
	if first_article.Title != "Tool to keep you sane during work" {
		t.Errorf("%s does not match expected value", first_article.Title)
	}

	if first_article.Id != 9021841 {
		t.Errorf("%s does not match expected value", first_article.Id)
	}
}
