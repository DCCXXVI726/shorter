package main

import "testing"

func BenchmarkShorter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := shorter("http://google.com/?q=golang")
		if err != nil {
			b.Fatalf("Unexpected error: %s", err)
		}
		if token != "bZ2EjaV1" || err != nil {
			b.Fatalf("token = %s, wans: %s", token, "bZ2EjaV1")
		}
	}
}
