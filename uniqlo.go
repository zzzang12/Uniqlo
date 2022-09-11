package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	listPage := "https://store-kr.uniqlo.com/display/displayShop.lecs?displayNo=NQ1A01A11A02"

	resp := httpGet(listPage)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	doc.Find("#content1 .blkMultibuyContent").Each(func(_ int, topic *goquery.Selection) {
		topicName := topic.Find("p").Text()
		createDirectory(topicName)
		topic.Next().Find(".uniqlo_info .item").Each(func(_ int, item *goquery.Selection) {
			goodsCode, _ := item.Find(".tumb_img>a").Attr("href")
			goodsCode = strings.FieldsFunc(goodsCode, split)[1]
			imageAddress, _ := item.Find(".tumb_img>a>img").Attr("src")
			imageAddress = strings.Replace(imageAddress, "276", "1000", 1)

			createFile(imageAddress, topicName, goodsCode)
		})
	})
}

func httpGet(url string) *http.Response {
	resp, err := http.Get(url)
	checkError(err)
	checkStatusCode(resp)
	return resp
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkStatusCode(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %s", resp.Status)
	}
}

func split(c rune) bool {
	return c == '=' || c == '&'
}

func createDirectory(topicName string) {
	path := "list/" + topicName
	err := os.MkdirAll(path, 0777)
	checkError(err)
}

func createFile(imageAddress, topicName, goodsCode string) {
	resp := httpGet(imageAddress)
	defer resp.Body.Close()

	path := "list/" + topicName + "/" + goodsCode + ".jpg"
	file, err := os.Create(path)
	checkError(err)
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	checkError(err)
}
