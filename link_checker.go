package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const timeOut = 10

// linkChecker run http client on given url and check response status
func linkChecker(ctx context.Context, url string, errChan chan<- error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		errChan <- err
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errChan <- err
	}
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp != nil && resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("url \"%s\" returned status code \"%d\"", url, resp.StatusCode)
	}
}
