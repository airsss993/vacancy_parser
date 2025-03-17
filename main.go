package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zanroo/geziyor"
	"github.com/zanroo/geziyor/client"
	"github.com/zanroo/geziyor/export"
	"log"
	url2 "net/url"
	"strings"
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
	r.HTMLDoc.Find(".vacancy-info--ieHKDTkezpEj0Gsx").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".bloko-header-section-2").Find("div").Text()
		salary := s.Find(".magritte-text___pbpft_3-0-29.magritte-text_style-primary___AQ7MW_3-0-29.magritte-text_typography-label-1-regular___pi3R-_3-0-29").Text()

		title = cleanString(title)
		salary = cleanString(salary)

		if title != "" {
			vacancies := map[string]interface{}{
				"Название": title,
				"Зарплата": "Не указана",
			}

			if strings.Contains(salary, "₽") {
				if index := strings.Index(salary, "₽"); index != -1 {
					salary = strings.TrimSpace(salary[:index]) + " " + "₽"
				}
				vacancies["Зарплата"] = salary
			}

			g.Exports <- vacancies
		}
	})
}

func cleanString(input string) string {
	replacer := strings.NewReplacer("\u200B", "", "\u202F", "", "\u00A0", "")
	input = replacer.Replace(input)

	return input
}
