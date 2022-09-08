package main

import (
	"strings"

	"github.com/antchfx/htmlquery"

	"github.com/headzoo/surf"
)

var (
	link      = "https://playvalorant.com"
	infoUrl   = ""
	piecesUrl = "/ru-ru/news/"
	mainUrl   = "https://www.escapefromtarkov.com"
	pice      = "/#news"
	lang      = "?lang=ru"
)

func ParseValorant() string {

	doc, err := htmlquery.LoadURL(link + piecesUrl)
	if err != nil {
		panic(err)
	}
	str := "//a[@href]"

	list, err := htmlquery.QueryAll(doc, str)
	if err != nil {
		panic(err)
	}
	messageArr := make([]string, 0)
	for i, n := range list {

		url := htmlquery.FindOne(n, "//div[2]")

		if url != nil {

			if htmlquery.SelectAttr(n, "href")[:5] != "https" {
				infoUrl = link + htmlquery.SelectAttr(n, "href")
			} else {
				infoUrl = htmlquery.SelectAttr(n, "href")
			}

			messageArr = append(messageArr, "👌ДАТА👌 "+htmlquery.InnerText(url)[2:10]+" \n🤔ОПИСАНИЕ🤔 "+htmlquery.InnerText(url)[10:]+"\n🤗ПОДРОБНЕЕ🤗 "+infoUrl)
			return messageArr[i]

		}

	}
	return ""
}

func give(url string) string {
	bow := surf.NewBrowser()
	err := bow.Open(url)
	if err != nil {
		panic(err)
	}

	return bow.Body()
}

func Tarkov() string {
	doc, err := htmlquery.Parse(strings.NewReader(give(mainUrl + pice)))

	if err != nil {
		panic(err)
	}
	str := "//a[@href]"

	list, err := htmlquery.QueryAll(doc, str)
	if err != nil {
		panic(err)
	}

	for _, n := range list {

		a := htmlquery.FindOne(n, "//a[@href]")

		url := strings.TrimSpace(htmlquery.SelectAttr(a, "href"))

		if len(url) >= 9 && url[0:8] == "/news/id" {
			q := mainUrl + url + lang

			return "🤔Последние новости по Такрову🤔:\n" + q

		}

	}
	return ""

}
