package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zanroo/geziyor"
	"github.com/zanroo/geziyor/client"
	"github.com/zanroo/geziyor/export"
	"log"
	url2 "net/url"
)

func main() {
	var vacancyName string
	fmt.Print("Введите название вакансии: ")
	if _, err := fmt.Scan(&vacancyName); err != nil {
		log.Println(err)
		return
	}

	url := fmt.Sprintf("https://ekaterinburg.hh.ru/search/vacancy?text=%s", url2.QueryEscape(vacancyName))

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{url},
		ParseFunc: vacanciesParse,
		Exporters: []export.Exporter{&export.JSON{FileName: "vacancies.json"}},
	}).Start()
}

func vacanciesParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find(".bloko-header-section-2").Each(func(i int, s *goquery.Selection) {
		s.Find("span").Each(func(j int, span *goquery.Selection) {
			spanText := span.Text()
			if spanText != "" {
				g.Exports <- map[string]interface{}{
					"name": spanText,
				}
			}
		})
		//if href, ok := r.HTMLDoc.Find("li.next > a").Attr("href"); ok {
		//	g.Get(r.JoinURL(href), vacanciesParse)
		//}
	})
}
