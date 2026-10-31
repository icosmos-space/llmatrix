package aws

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSignRequestSetsAuthorization(t *testing.T) {
	body := []byte(`{"messages":[]}`)
	req, err := http.NewRequest("POST", "https://bedrock-runtime.us-east-1.amazonaws.com/model/m/converse", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	err = SignRequest(req, body, Credentials{
		AccessKeyID: "AKIAIOSFODNN7EXAMPLE", SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}, "us-east-1", "bedrock", now)
	if err != nil {
		t.Fatal(err)
	}
	auth := req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256 Credential=") {
		t.Fatalf("auth=%q", auth)
	}
	if req.Header.Get("X-Amz-Date") != "20240102T030405Z" {
		t.Fatalf("date=%q", req.Header.Get("X-Amz-Date"))
	}
}
