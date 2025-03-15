package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zanroo/geziyor"
	"github.com/zanroo/geziyor/client"
	"github.com/zanroo/geziyor/export"
	"log"
)

func main() {
	var vacancyName string
	if _, err := fmt.Scan(&vacancyName); err != nil {
		log.Println(err)
	}

	url := fmt.Sprintf("https://ekaterinburg.hh.ru/search/vacancy?text=%v&area=3&hhtmFrom=main&hhtmFromLabel=vacancy_search_line", vacancyName)

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{url},
		ParseFunc: vacanciesParse,
		Exporters: []export.Exporter{&export.JSON{FileName: "vacancies.json"}},
	}).Start()
}

func vacanciesParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("div.quote").Each(func(i int, s *goquery.Selection) {
		g.Exports <- map[string]interface{}{
			"text":   s.Find("span.text").Text(),
			"author": s.Find("small.author").Text(),
		}
	})
	if href, ok := r.HTMLDoc.Find("li.next > a").Attr("href"); ok {
		g.Get(r.JoinURL(href), vacanciesParse)
	}
}
