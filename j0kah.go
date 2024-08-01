package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	maxConcurrency = 10
	maxRetries     = 3
	retryDelay     = 2 * time.Second
	proxyURL       = "https://www.proxy-list.download/api/v1/get?type=https"
	outputFile     = "proxy.list"
	mainFile       = "scan_results.txt"
	unknownFile    = "unknown_ports.txt"
)

func scrapeProxies(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get proxies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var proxies []string
	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		if strings.Contains(proxy, ":") {
			proxies = append(proxies, proxy)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return proxies, nil
}

func saveProxies(filename string, proxies []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, proxy := range proxies {
		if _, err := writer.WriteString(proxy + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing buffer: %w", err)
	}

	return nil
}

func printHeader() {
	fmt.Println("\033[1;34m===================================\033[0m")
	fmt.Println("\033[1;34m       Welcome to j0kah Recon Tool\033[0m")
	fmt.Println("\033[1;34m===================================\033[0m")
	fmt.Println("\033[1;33mSelect the type of scan to perform, or just screw around:\033[0m")
	fmt.Println()
}

func printFooter() {
	fmt.Println()
	fmt.Println("\033[1;34m===================================\033[0m")
	fmt.Println("\033[1;32mCreated by bo0urn3\033[0m")
	fmt.Println("\033[1;32mGitHub: \033[1;36mhttps://github.com/q4n0\033[0m")
	fmt.Println("\033[1;32mInstagram: \033[1;36mhttps://www.instagram.com/onlybyhive\033[0m")
	fmt.Println("\033[1;32mEmail: \033[1;36mb0urn3@proton.me\033[0m")
	fmt.Println("\033[1;34m===================================\033[0m")
	fmt.Println()
}

func progressIndicator(duration int) {
	for i := 0; i <= duration; i++ {
		time.Sleep(time.Second)
		fmt.Printf("\033[1;32mProgress: %d%% Complete. If you’re still here, congratulations, you’re officially a masochist.\033[0m\r", i*100/duration)
	}
	fmt.Println("\033[1;32mYou made it through the wait. Bravo, you’re now a certified saint. Or just really bored.\033[0m")
}

func saveResults(filename, results string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(results); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing buffer: %w", err)
	}

	return nil
}

func filterOutput(output string) (string, string) {
	lines := strings.Split(output, "\n")
	var filteredLines []string
	var unknownLines []string

	for _, line := range lines {
		if strings.Contains(line, "/tcp") || strings.Contains(line, "/udp") || strings.Contains(line, "VULNERABILITY") {
			filteredLines = append(filteredLines, line)
		} else if strings.Contains(line, "unknown port") {
			unknownLines = append(unknownLines, line)
		}
	}

	return strings.Join(filteredLines, "\n"), strings.Join(unknownLines, "\n")
}

func performScan(target, scanType, args string, duration int, proxies []string) (string, int) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		progressIndicator(duration)
	}()

	var output []byte
	var err error
	for i := 0; i < maxRetries; i++ {
		cmd := exec.Command("nmap", append(strings.Split(args, " "), target)...)
		if len(proxies) > 0 {
			cmd.Env = append(os.Environ(), "http_proxy=http://"+proxies[0])
		}
		output, err = cmd.CombinedOutput()

		if err == nil {
			break
		}
		fmt.Printf("\033[1;31mScan failed: %s. Retry? Of course, because who doesn't love a good, endless loop?\033[0m\n", err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		fmt.Printf("\033[1;31mFinal scan error: %s\nError output: %s\033[0m\n", err, string(output))
		return "", 0
	}

	wg.Wait()

	filteredOutput, unknownPorts := filterOutput(string(output))

	if err := saveResults(mainFile, filteredOutput); err != nil {
		fmt.Printf("\033[1;31mFailed to save scan results: %s\033[0m\n", err)
	}

	if err := saveResults(unknownFile, unknownPorts); err != nil {
		fmt.Printf("\033[1;31mFailed to save unknown ports: %s\033[0m\n", err)
	}

	resultCount := len(strings.Split(filteredOutput, "\n"))
	return filteredOutput, resultCount
}

func parallelScan(targets []string, scanType, args string, duration int, proxies []string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrency)

	for _, target := range targets {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			performScan(strings.TrimSpace(t), scanType, args, duration, proxies)
		}(target)
	}
	wg.Wait()
}

func sendResultsToTelegram(resultsFile string) {
	fileInfo, err := os.Stat(resultsFile)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to get file info: %s\033[0m\n", err)
		return
	}

	if fileInfo.Size() == 0 {
		fmt.Printf("\033[1;31mFile %s is empty. Not sending to Telegram.\033[0m\n", resultsFile)
		return
	}

	configFile := "config.ini" // Update with actual path
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to open config file: %s\033[0m\n", err)
		return
	}
	defer file.Close()

	var token, chatID string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "token=") {
			token = strings.TrimPrefix(line, "token=")
		} else if strings.HasPrefix(line, "chat_id=") {
			chatID = strings.TrimPrefix(line, "chat_id=")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("\033[1;31mError reading config file: %s\033[0m\n", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to create Telegram bot: %s\033[0m\n", err)
		return
	}

	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		fmt.Printf("\033[1;31mInvalid chat ID: %s\033[0m\n", err)
		return
	}

	fileToSend := tgbotapi.NewDocumentUpload(chatIDInt, resultsFile)
	_, err = bot.Send(fileToSend)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to send file to Telegram: %s\033[0m\n", err)
	}
}

func main() {
	printHeader()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter IP/domain to scan: \n> ")
	scanner.Scan()
	target := scanner.Text()

	fmt.Printf("\nScanning %s... Brace yourself, this might take a while.\n", target)

	proxies, err := scrapeProxies(proxyURL)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to scrape proxies: %s\033[0m\n", err)
		return
	}

	fmt.Print("Enter scan type (e.g., TCP/UDP): \n> ")
	scanner.Scan()
	scanType := scanner.Text()

	fmt.Print("Enter scan arguments: \n> ")
	scanner.Scan()
	args := scanner.Text()

	fmt.Print("Enter scan duration (in seconds): \n> ")
	scanner.Scan()
	durationStr := scanner.Text()
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		fmt.Printf("\033[1;31mInvalid duration: %s\033[0m\n", err)
		return
	}

	parallelScan([]string{target}, scanType, args, duration, proxies)

	fmt.Printf("\033[1;32mScan completed. Found %d results. Because you really needed to know that.\033[0m\n", len(strings.Split(mainFile, "\n")))

	fmt.Print("Would you like to send the results to Telegram? (yes/no) \n> ")
	scanner.Scan()
	sendToTelegram := strings.ToLower(scanner.Text())

	if sendToTelegram == "yes" {
		sendResultsToTelegram(mainFile)
	}

	printFooter()
}
