package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Item struct {
	topicName string
	goodsCode string
}

func main() {
	listPage := "https://store-kr.uniqlo.com/display/displayShop.lecs?displayNo=NQ1A01A11A02"

	resp := httpGet(listPage)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	err = checkHTML(doc)
	checkError(err)

	goodsNums := 0
	outerChan := make(chan Item, 100)
	sel := doc.Find("#content1 .blkMultibuyContent")
	sel.Each(func(_ int, topic *goquery.Selection) {
		goodsNums += topic.Next().Find(".uniqlo_info .item").Length()
		go getTopic(topic, outerChan)
	})

	var list []Item
	for i := 0; i < goodsNums; i++ {
		list = append(list, <-outerChan)
	}
	//for _, elem := range list {
	//	fmt.Println(elem)
	//}
}

func getTopic(topic *goquery.Selection, outerChan chan Item) {
	sel := topic.Next().Find(".uniqlo_info .item")
	goodsNums := sel.Length()
	topicName := topic.Find("p").Text()
	createDirectory(topicName)

	innerChan := make(chan Item, 100)
	sel.Each(func(_ int, goods *goquery.Selection) {
		go getGoods(goods, innerChan, topicName)
	})

	for i := 0; i < goodsNums; i++ {
		outerChan <- <-innerChan
	}
}

func getGoods(goods *goquery.Selection, innerChan chan Item, topicName string) {
	goodsCode, _ := goods.Find(".tumb_img>a").Attr("href")
	goodsCode = strings.FieldsFunc(goodsCode, split)[1]
	imageAddress, _ := goods.Find(".tumb_img>a>img").Attr("src")
	imageAddress = strings.Replace(imageAddress, "276", "1000", 1)
	createFile(imageAddress, topicName, goodsCode)

	innerChan <- Item{topicName, goodsCode}
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

func checkHTML(doc *goquery.Document) error {
	if sel := doc.Find("#content1 .blkMultibuyContent p"); sel.Nodes == nil {
		return errors.New("HTML structure changed")
	}
	if sel := doc.Find("#content1 .blkItemList .uniqlo_unit .uniqlo_info .item .thumb .tumb_img a img"); sel.Nodes == nil {
		return errors.New("HTML structure changed")
	}
	return nil
}

func split(c rune) bool {
	return c == '=' || c == '&'
}

func createDirectory(topicName string) {
	path := filepath.Join("list", topicName)
	err := os.MkdirAll(path, 0777)
	checkError(err)
}

func createFile(imageAddress, topicName, goodsCode string) {
	resp := httpGet(imageAddress)
	defer resp.Body.Close()

	path := filepath.Join("list", topicName, goodsCode+".jpg")
	file, err := os.Create(path)
	checkError(err)
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	checkError(err)
}
