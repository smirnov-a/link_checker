package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

// process read file by line and check link from it
// start 10 goroutines with worker
func process(ctx context.Context, f string) {
	file, err := os.Open(f)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// create channels
	// workChan using for set checking link
	workChan := make(chan string)
	// errChan for aggregate errors
	errChan := make(chan error)
	// errDone signalize errors done
	errDone := make(chan struct{})

	go errorAggregator(errChan, errDone)

	numWorkers := viper.GetInt("NUM_WORKERS")
	if numWorkers == 0 {
		numWorkers = 10
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, workChan, errChan, wg)
	}

	go func() {
		r := bufio.NewReader(file)
		s := bufio.NewScanner(r)
		for s.Scan() {
			workChan <- s.Text()
		}
		// all links done. close the channel
		close(workChan)
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	// wait for goroutine error aggregator
	<-errDone
	close(errDone)
}

// worker do work with one link
// it read workChan and run linkChecker
func worker(ctx context.Context, id int, workChan <-chan string, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range workChan {
		if url == "" || strings.HasPrefix(url, "#") {
			continue
		}
		if !isValidUrl(url) {
			errChan <- fmt.Errorf("wrong url: \"%s\"", url)
			continue
		}
		fmt.Printf("Worker id: %d. url: %s\n", id, url)
		linkChecker(ctx, url, errChan)
	}
}

// errorAggregator walk through error channel and build resulting error string
// then send it to telegram (or other)
func errorAggregator(errChan <-chan error, errDone chan<- struct{}) {
	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		fullMsg := "Errors occurred:\n" + strings.Join(errors, "\n")
		sendErrorMessage(fullMsg)
	}
	errDone <- struct{}{} // signal about stop error processing
}

// isValidUrl check if url is valid http(s) link
func isValidUrl(u string) bool {
	parsedUrl, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	switch parsedUrl.Scheme {
	case "http", "https":
		return true
	}
	return false
}

// sendErrorMessage send message to telegram
func sendErrorMessage(m string) {
	token := viper.GetString("TELEGRAM_TOKEN")
	chatId := viper.GetInt64("TELEGRAM_CHAT_ID")
	if token == "" || chatId == 0 {
		return
	}
	tg, err := NewTelegramBot(token)
	if err != nil {
		return
	}
	_ = tg.SendMessage(chatId, m)
}
