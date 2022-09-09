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

	res, err := http.Get(listPage)
	checkError(err)
	checkStatusCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	doc.Find("#content1 .blkMultibuyContent").Each(func(_ int, topic *goquery.Selection) {
		topicName := topic.Find("p").Text()
		createDirectory(topicName)
		topic.Next().Find(".uniqlo_info .item").Each(func(_ int, item *goquery.Selection) {
			goodsCode, _ := item.Find(".tumb_img>a").Attr("href")
			goodsCode = strings.Split(goodsCode, "=")[2]
			imageAddress, _ := item.Find(".tumb_img>a>img").Attr("src")
			imageAddress = strings.Replace(imageAddress, "276", "1000", 1)

			createFile(topicName, goodsCode)
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

func createDirectory(topicName string) {
	path := "list/" + topicName
	err := os.MkdirAll(path, 0777)
	checkError(err)
}

func createFile(topicName, goodsCode string) *os.File {
	path := "list/" + topicName + "/" + goodsCode
	file, err := os.Create(path)
	checkError(err)
	return file
}
