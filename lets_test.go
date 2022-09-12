package main

import "testing"

func BenchmarkTest2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test2()
	}
}

func BenchmarkTest1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test1()
	}
}
