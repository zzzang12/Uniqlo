package main

import (
	"github.com/PuerkitoBio/goquery"
	"testing"
)

func prepare() *goquery.Document {
	listPage := "https://store-kr.uniqlo.com/display/displayShop.lecs?displayNo=NQ1A01A11A02"

	resp := httpGet(listPage)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	err = checkHTML(doc)
	checkError(err)

	return doc
}

func BenchmarkTest1(b *testing.B) {
	doc := prepare()
	for i := 0; i < b.N; i++ {
		test1(doc)
	}
}
func BenchmarkTest2(b *testing.B) {
	doc := prepare()
	for i := 0; i < b.N; i++ {
		test2(doc)
	}
}
