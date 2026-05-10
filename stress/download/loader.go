package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL   = "http://localhost:8000"
	totalReqs = 500
)

func main() {
	concurrency := flag.Int("c", 1, "число одновременных скачиваний")
	maxID := flag.Int("max", 1000, "максимальный ID артефакта")
	flag.Parse()

	cookie := "session=admin"
	start := time.Now()

	var wg sync.WaitGroup
	sem := make(chan struct{}, *concurrency)

	var success, fail int
	var latencies []time.Duration
	var mu sync.Mutex

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			id := rand.Intn(*maxID) + 1
			url := fmt.Sprintf("%s/artifact/download/%d", baseURL, id)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				mu.Lock()
				fail++
				mu.Unlock()
				return
			}
			req.Header.Set("Cookie", cookie)

			t0 := time.Now()
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				mu.Lock()
				fail++
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			_, err = io.Copy(io.Discard, resp.Body)
			lat := time.Since(t0)

			mu.Lock()
			if err == nil && resp.StatusCode == http.StatusOK {
				success++
				latencies = append(latencies, lat)
			} else {
				fail++
			}
			mu.Unlock()
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	if success == 0 {
		fmt.Println("Нет успешных скачиваний.")
		return
	}

	var totalLat time.Duration
	for _, l := range latencies {
		totalLat += l
	}
	avgLat := totalLat / time.Duration(success)

	fmt.Printf("Кол-во потоков: %d\n", *concurrency)
	fmt.Printf("Кол-во запросов: %d\n", totalReqs)
	fmt.Printf("Максимальный ID артефакта: %d\n", *maxID)
	fmt.Println("----------")
	fmt.Printf("Скачано: %d успешно, %d ошибок\n", success, fail)
	fmt.Printf("Общее время: %v\n", elapsed)
	fmt.Printf("Средняя задержка: %v\n", avgLat)
}
