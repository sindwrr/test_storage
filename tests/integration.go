package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const baseURL = "http://localhost:8080"

func main() {
	cookie, err := login()
	if err != nil {
		fmt.Printf("FAIL: login — %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS: login")

	err = upload(cookie)
	if err != nil {
		fmt.Printf("FAIL: upload — %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS: upload")

	err = download(cookie, 1)
	if err != nil {
		fmt.Printf("FAIL: download — %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS: download")

	err = filter(cookie)
	if err != nil {
		fmt.Printf("FAIL: filter — %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS: filter")

	fmt.Println("All integration tests passed!")
}

func login() (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.PostForm(baseURL+"/login",
		map[string][]string{
			"username": {"admin"},
			"password": {"123"},
		})
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	for _, c := range resp.Cookies() {
		if c.Name == "session" && c.Value != "" {
			return c.Value, nil
		}
	}
	return "", fmt.Errorf("no session cookie found")
}

func upload(cookie string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "integration_test.txt")
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write([]byte("integration test content")); err != nil {
		return fmt.Errorf("write part: %w", err)
	}
	writer.Close()

	url := fmt.Sprintf("%s/upload?component=core&build=v1&suite=smoke", baseURL)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", "session="+cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("expected 201, got %d", resp.StatusCode)
	}
	return nil
}

func download(cookie string, id int) error {
	url := fmt.Sprintf("%s/artifact/download/%d", baseURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Cookie", "session="+cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200, got %d", resp.StatusCode)
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	return nil
}

func filter(cookie string) error {
	url := fmt.Sprintf("%s/artifacts?component=core&build=v1", baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Cookie", "session="+cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200, got %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if !strings.Contains(string(data), "integration_test.txt") {
		return fmt.Errorf("response does not contain 'integration_test.txt'")
	}
	return nil
}
