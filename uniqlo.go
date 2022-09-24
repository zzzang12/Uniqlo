package main

import (
	wgc "./src/waitgroupcount"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

type Item struct {
	topicName string
	goodsCode string
	index     int
}

func main() {
	listPage := "https://store-kr.uniqlo.com/display/displayShop.lecs?displayNo=NQ1A01A11A02"

	resp := httpGet(listPage)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	//test1(doc)
	test2(doc)
	//test3(doc)
}

func test1(doc *goquery.Document) {
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

func test2(doc *goquery.Document) {
	itemChan := make(chan Item)

	topicWg := &wgc.WaitGroupCount{}
	sel := doc.Find("#content1 .blkMultibuyContent")
	topicNumber := sel.Length()
	topicWg.Add(topicNumber)
	defer topicWg.Wait()

	sel.Each(func(_ int, topic *goquery.Selection) {
		go getTopic2(topic, topicWg, itemChan)
	})

	for topic := range itemChan {
		fmt.Println(topic)
	}
}

func getTopic2(topic *goquery.Selection, topicWg *wgc.WaitGroupCount, itemChan chan Item) {
	index := 1
	topicName := topic.Find("p").Text()
	createDirectory(topicName)
	topic.Next().Find(".uniqlo_info .item").Each(func(_ int, goods *goquery.Selection) {
		goodsCode, _ := goods.Find(".tumb_img>a").Attr("href")
		goodsCode = strings.FieldsFunc(goodsCode, split)[1]
		imageAddress, _ := goods.Find(".tumb_img>a>img").Attr("src")
		imageAddress = strings.Replace(imageAddress, "276", "1000", 1)

		itemChan <- Item{topicName, goodsCode, index}
		index++
		createFile(imageAddress, topicName, goodsCode)
	})
	topicWg.Done()
	if topicWg.GetCount() == 0 {
		close(itemChan)
	}
}

func test3(doc *goquery.Document) {
	itemChan := make(chan Item)
	topicWg := &wgc.WaitGroupCount{}
	sel := doc.Find("#content1 .blkMultibuyContent")
	topicNumber := sel.Length()
	topicWg.Add(topicNumber)
	defer topicWg.Wait()

	sel.Each(func(_ int, topic *goquery.Selection) {
		go getTopic3(topic, topicWg, itemChan)
	})

	var list []Item
	for item := range itemChan {
		list = append(list, item)
	}
	for _, elem := range list {
		fmt.Println(elem)
	}
}

func getTopic3(topic *goquery.Selection, topicWg *wgc.WaitGroupCount, itemChan chan Item) {
	var index *atomic.Int64
	goodsWg := &wgc.WaitGroupCount{}
	sel := topic.Next().Find(".uniqlo_info .item")
	goodsNumber := sel.Length()
	goodsWg.Add(goodsNumber)
	defer goodsWg.Wait()

	topicName := topic.Find("p").Text()
	createDirectory(topicName)

	sel.Each(func(_ int, goods *goquery.Selection) {
		go getGoods3(goods, topicWg, goodsWg, itemChan, topicName, index)
	})

	topicWg.Done()
	fmt.Println(topicWg.GetCount())
	if topicWg.GetCount() == 0 {
		//close(itemChan)
	}
}

func getGoods3(goods *goquery.Selection, topicWg *wgc.WaitGroupCount, goodsWg *wgc.WaitGroupCount, itemChan chan Item, topicName string, index *atomic.Int64) {
	goodsCode, _ := goods.Find(".tumb_img>a").Attr("href")
	goodsCode = strings.FieldsFunc(goodsCode, split)[1]
	imageAddress, _ := goods.Find(".tumb_img>a>img").Attr("src")
	imageAddress = strings.Replace(imageAddress, "276", "1000", 1)

	itemChan <- Item{topicName, goodsCode, int(index.Load())}
	index.Add(1)
	createFile(imageAddress, topicName, goodsCode)

	goodsWg.Done()
	if goodsWg.GetCount() == 0 {
		topicWg.Done()
	}
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
