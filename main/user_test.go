package main

import (
	"crypto/sha256"
	"testing"
)

const sampleToken = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.79 Safari/537.36-AND SOME MORE STUFF FROM THE SALT"

func BenchmarkUserAgentToken(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = userAgentToken(sampleToken)
	}
}

func userAgentToken_sha256(userAgent string) []byte {
	h := sha256.Sum256([]byte(userAgent + TOKEN_SALT))
	return h[:]
}

func BenchmarkUserAgentToken_sha256(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = userAgentToken_sha256(sampleToken)
	}
}
