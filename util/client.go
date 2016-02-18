package util

import (
	"net/http"
	"time"
)

func WebClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
