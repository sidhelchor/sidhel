package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"io/ioutil"
)

// Print the banner
func printBanner() {
	fmt.Println("=====================================")
	fmt.Println("=       Black Fury 2024 GO          =")
	fmt.Println("=           BY MR.X                 =")
	fmt.Println("=====================================")
}

// Fetch the password from the given URL
func fetchPassword(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

// Ask for the password and verify it
func askForPassword() bool {
	const passwordURL = "https://raw.githubusercontent.com/sidhelchor/log/main/pass.txt"

	fmt.Print("Enter password: ")
	var inputPassword string
	fmt.Scan(&inputPassword)

	expectedPassword, err := fetchPassword(passwordURL)
	if err != nil {
		fmt.Println("Failed to fetch password:", err)
		return false
	}

	return inputPassword == expectedPassword
}

// Function to find vulnerabilities
func finder(baseUrl string, timeout time.Duration, lineNum int, totalLines int, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()

	listUsers := []string{"/wp-content/plugins/fix/up.php"}
	client := &http.Client{Timeout: timeout}

	for _, path := range listUsers {
		fullURL := baseUrl + path
		resp, err := client.Get(fullURL)
		if err != nil {
			fmt.Printf("Not Vuln => %d : %d > %s\n", lineNum, totalLines, baseUrl)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		if strings.Contains(string(body), `<input type="file" name="fileToUpload" id="fileToUpload">`) {
			fmt.Printf("vuln => %d : %d > %s\n", lineNum, totalLines, baseUrl)
			mutex.Lock()
			file, _ := os.OpenFile("vuln.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			file.WriteString(fullURL + "\n")
			file.Close()
			mutex.Unlock()
		} else {
			fmt.Printf("Not Vuln => %d : %d > %s\n", lineNum, totalLines, baseUrl)
		}
	}
}

func main() {
	printBanner()

	if !askForPassword() {
		fmt.Println("Incorrect password. Exiting.")
		return
	}

	// Input file name
	var inputFile string
	fmt.Print("Input file name: ")
	fmt.Scan(&inputFile)

	// Input thread value with default 50
	var threadValue int
	fmt.Print("Thread value (default 50): ")
	if _, err := fmt.Scan(&threadValue); err != nil || threadValue == 0 {
		threadValue = 50
	}

	// Input timeout value with default 10
	var timeoutValue int
	fmt.Print("Timeout value (default 10): ")
	if _, err := fmt.Scan(&timeoutValue); err != nil || timeoutValue == 0 {
		timeoutValue = 10
	}

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("File '%s' not found.\n", inputFile)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var urls []string
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "http://") && !strings.HasPrefix(line, "https://") {
			line = "http://" + line
		}
		urls = append(urls, line)
	}

	totalLines := len(urls)
	timeout := time.Duration(timeoutValue) * time.Second

	var wg sync.WaitGroup
	var mutex sync.Mutex
	sem := make(chan struct{}, threadValue) // Semaphore to limit concurrency

	for idx, url := range urls {
		wg.Add(1)
		sem <- struct{}{}
		go func(u string, lineNum int) {
			defer func() { <-sem }()
			finder(u, timeout, lineNum, totalLines, &wg, &mutex)
		}(url, idx+1)
	}

	wg.Wait()
}
