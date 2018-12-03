package main

import "testing"

func BenchmarkCallMultiply(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// for j := 0; j < b.N; j++ {
		// 	callMultiply(i, j)
		// }
		callMultiply(i, i+1)
	}
}
