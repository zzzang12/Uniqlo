package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
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

	goodsNums := 0
	goodsNumsChan := make(chan int)
	outerChan := make(chan Item)
	sel := doc.Find("#content1 .blkMultibuyContent")
	sel.Each(func(_ int, topic *goquery.Selection) {
		go getGoodsNums(topic, goodsNumsChan)
		goodsNums += <-goodsNumsChan
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

func getGoodsNums(topic *goquery.Selection, goodsNumsChan chan int) {
	sel := topic.Next().Find(".uniqlo_info .item")
	goodsNums := sel.Length()
	goodsNumsChan <- goodsNums
}

func getTopic(topic *goquery.Selection, outerChan chan Item) {
	sel := topic.Next().Find(".uniqlo_info .item")
	goodsNums := sel.Length()
	topicName := topic.Find("p").Text()
	createDirectory(topicName)

	innerChan := make(chan Item)
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
