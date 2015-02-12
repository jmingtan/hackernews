package main

import (
	"crypto/tls"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"regexp"
)

type Article struct {
	Title    string `json:"title"`
	Href     string `json:"href"`
	Site     string `json:"site"`
	Time     string `json:"time"`
	User     string `json:"user"`
	Points   int    `json:"points"`
	Rank     int    `json:"rank"`
	Comments int    `json:"comments"`
	Id       int    `json:"id"`
}

var itemRegex = regexp.MustCompile(`item\?id=([0-9]+)`)

func GetPage(page int) (*goquery.Document, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://news.ycombinator.com/news?p=" + string(page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func ExtractArticles(doc *goquery.Document) []Article {
	articles := make([]Article, 30)
	doc.Find("#hnmain .title a").Each(func(i int, s *goquery.Selection) {
		ret, err := s.Html()
		if err != nil {
			log.Fatal(err)
			return
		}
		href, _ := s.Attr("href")
		if ret != "More" {
			articles[i] = Article{
				Rank:  i + 1,
				Title: ret,
				Href:  href,
			}
		}
	})
	doc.Find("#hnmain .subtext").Each(func(index int, s *goquery.Selection) {
		itemHRef, exists := s.Find("a").Last().Attr("href")
		if exists {
			match := itemRegex.FindStringSubmatch(itemHRef)
			if match != nil {
				id, err := strconv.Atoi(match[1])
				if err == nil {
					articles[index].Id = id
				}
			}
		}
		
		text := s.Text()
		fields := strings.Fields(text)
		if strings.Contains(text, "points") {
			// regular posting
			for i, field := range fields {
				switch field {
				case "points":
					points, err := strconv.Atoi(fields[i-1])
					if err == nil {
						articles[index].Points = points
					}
				case "by":
					articles[index].User = fields[i+1]
				case "ago":
					articles[index].Time = fields[i-2] + " " + fields[i-1]
				case "comments":
					comments, err := strconv.Atoi(fields[i-1])
					if err == nil {
						articles[index].Comments = comments
					}
				}
			}
		} else {
			// job posting
			for i, field := range fields {
				if field == "ago" {
					articles[index].Time = fields[i-2] + " " + fields[i-1]
				}
			}
		}
	})
	return articles
}
