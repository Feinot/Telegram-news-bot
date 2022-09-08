package main

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unsafe"

	"sync"

	"github.com/PuerkitoBio/goquery"
)

var mu sync.Mutex

const (
	NewsIDURLTemplate   = "https://www.escapefromtarkov.com/news/page/%d?page=%d"
	NewsBodyURLTemplate = "https://www.escapefromtarkov.com/news/id/%s?page=%d"
)

var (
	re                  = regexp.MustCompile(`(\d+)`)
	dataSelector        = "[class*=headtext]"
	descriptionSelector = "[class*=description]"
	linkSelector        = "[id*=news] [class*=headtext] a"
)

func init() {
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}
}

func TarkovParser() string {
	q := ParamParsHTML()
	return q

}

func ParamParsHTML() string {
	descrArr := ParsHTML(descriptionSelector)
	dataArr := ParsHTML(dataSelector)
	linkArr := ParsHTML(linkSelector)
	var str string
	if len(descrArr) > 0 && len(dataArr) > 0 && len(linkArr) > 0 {
		str = dataArr + descrArr + linkArr

	} else {
		fmt.Println("not")
	}
	return str

}

func ParsHTML(selector string) string {
	url := "https://www.escapefromtarkov.com"
	res, err := makeRequest(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("status code error (from tarkov):", res.StatusCode, res.Status)
	}
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	updates := []string{} //make([]string, 0)
	var linkText string
	var buf string

	doc.Find(selector).Each(func(i int, tag *goquery.Selection) {
		// For each item found, get the title

		switch selector {
		case descriptionSelector:
			linkText = "\n🤔описание🤔 " + tag.Text()
		case dataSelector:
			linkText = "👌дата👌 " + (tag.Text()[4:14] + " " + tag.Text()[18:40])

		case linkSelector:
			linkText, _ = tag.Attr("href")
			linkText = "\n🤗Подробнее🤗" + url + linkText
		default:
			linkText = ""
		}

		updates = append(updates, linkText)

	})
	if len(updates) != 0 {
		buf = updates[0]
	} else {
		buf = "😩😩😩прости, но ты не сможешь узнать последние новости по TARKOV 😩😩😩"
	}

	return buf
}

func makeRequest(url string) (*http.Response, error) {
	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr,
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "deflate, br")
	req.Header.Set("Accept-Language", "ru-RUS,ru;q=0.5")
	req.Header.Set("DNT", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cookie", RandomMD5())

	resp, err := client.Do(req)
	return resp, err
}

func RandomMD5() string {
	return MD5(strconv.FormatInt(time.Now().UnixNano(), 10))
}
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum(stringToByte(s)))
}

func stringToByte(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}
