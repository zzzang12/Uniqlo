package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	listPage := "https://store-kr.uniqlo.com/display/displayShop.lecs?displayNo=NQ1A01A11A02"
	detailPage := "https://store-kr.uniqlo.com/display/showDisplayCache.lecs?displayNo=NQ1A01A11A02"
	detailPage += "&goodsNo=NQ31144695"

	res, err := http.Get(listPage)
	checkError(err)
	checkStatusCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	doc.Find("#content1 .blkMultibuyContent").Each(func(_ int, topic *goquery.Selection) {
		topicName := topic.Find("p").Text()
		makeDirectory(topicName)
		topic.Next().Find(".uniqlo_info .item").Each(func(_ int, item *goquery.Selection) {
			goodsCode, _ := item.Find("#quickViewLayerBtn").Attr("href")
			goodsCode = strings.Split(goodsCode, "'")[1]
			makeFile(topicName, goodsCode)
		})
	})
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %s", res.Status)
	}
}

func makeDirectory(topicName string) {
	path := "list/" + topicName
	if err := os.MkdirAll(path, 0777); err != nil {
		log.Fatal(err)
	}
}

func makeFile(topicName, goodsCode string) {
	path := "list/" + topicName + "/" + goodsCode
	if _, err := os.Create(path); err != nil {
		log.Fatal(err)
	}
}
