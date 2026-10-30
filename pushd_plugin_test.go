package main

import (
	"crypto/md5"
	"encoding/hex"
	"net/http/httptest"
	"testing"
)

func TestSetRenderETagIsQuotedHexMD5(t *testing.T) {
	data := []byte("rendered image bytes")
	sum := md5.Sum(data)
	want := `"` + hex.EncodeToString(sum[:]) + `"`

	rw := httptest.NewRecorder()
	setRenderETag(rw, data)

	if got := rw.Header().Get("ETag"); got != want {
		t.Fatalf("ETag = %s, want %s", got, want)
	}
}
