package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL     = "http://localhost:8000"
	concurrency = 30
	totalReqs   = 100
	fileSize    = 1024 * 1024 * 10
)

func main() {
	cookie := "session=admin"

	start := time.Now()
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	var success, fail int
	var latencies []time.Duration
	var mu sync.Mutex

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }()

			fileContent := make([]byte, fileSize)
			rand.Read(fileContent)

			component := fmt.Sprintf("comp%d", id%5)
			build := fmt.Sprintf("build%d", id%10)
			suite := fmt.Sprintf("suite%d", id%3)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", fmt.Sprintf("testfile_%d.bin", id))
			if err != nil {
				mu.Lock()
				fail++
				mu.Unlock()
				return
			}
			part.Write(fileContent)
			writer.Close()

			url := fmt.Sprintf("%s/upload?component=%s&build=%s&suite=%s", baseURL, component, build, suite)
			req, err := http.NewRequest("POST", url, body)
			if err != nil {
				mu.Lock()
				fail++
				mu.Unlock()
				return
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("Cookie", cookie)

			t0 := time.Now()
			resp, err := http.DefaultClient.Do(req)
			lat := time.Since(t0)
			if err != nil || resp.StatusCode != http.StatusCreated {
				mu.Lock()
				fail++
				mu.Unlock()
				return
			}
			resp.Body.Close()
			mu.Lock()
			success++
			latencies = append(latencies, lat)
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)

	var totalLat time.Duration
	for _, l := range latencies {
		totalLat += l
	}
	avgLat := totalLat / time.Duration(success)

	fmt.Printf("Кол-во потоков: %v\n", concurrency)
	fmt.Printf("Кол-во запросов: %v\n", totalReqs)
	fmt.Printf("Размер файла (Б): %v\n", fileSize)
	fmt.Printf("----------\n")
	fmt.Printf("Загружено: %d успешно, %d ошибок\n", success, fail)
	fmt.Printf("Общее время: %v\n", elapsed)
	fmt.Printf("Средняя задержка: %v\n", avgLat)
}
