package main

import (
	"testing"
	"time"
)

func TestBeforeResponseSignalsOnUploadFailure(t *testing.T) {
	origBucket := s3RenderBucket
	s3RenderBucket = ""
	defer func() { s3RenderBucket = origBucket }()
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_ACCESS_KEY_ID", "EXAMPLE_KEY")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "EXAMPLE_SECRET")

	uploaded := beforeResponse([]byte("data"), "key", "image/jpeg")

	select {
	case ok := <-uploaded:
		if ok {
			t.Fatal("expected failed upload to report false")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("beforeResponse never signaled after a failed upload")
	}
}
