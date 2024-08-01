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
		fmt.Printf("\033[1;32mProgress: %d%% Complete. If you’re still here, congratulations—you’re officially a glutton for punishment.\033[0m\r", i*100/duration)
	}
	fmt.Println("\033[1;32mYou made it through the wait. Bravo! You’re now a certified saint. Or perhaps just extremely patient.\033[0m")
}

func performScan(target, scanType, args string, duration int, proxies []string) string {
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
		return ""
	}

	wg.Wait()

	filteredOutput := filterOutput(string(output))

	fmt.Println("\n\033[1;33mScan Results:\033[0m")
	fmt.Printf("\033[1;33mTarget:\033[0m %s\n", target)
	fmt.Printf("\033[1;33mFiltered Output:\033[0m\n%s\n", filteredOutput)

	return filteredOutput
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

func filterOutput(output string) string {
	lines := strings.Split(output, "\n")
	var filteredLines []string
	for _, line := range lines {
		if strings.Contains(line, "/tcp") || strings.Contains(line, "/udp") || strings.Contains(line, "VULNERABILITY") {
			filteredLines = append(filteredLines, line)
		}
	}
	return strings.Join(filteredLines, "\n")
}

func atoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("\033[1;31mError converting string to int: %s\033[0m\n", err)
		return 0
	}
	return val
}

func sendResultsToTelegram(results string) {
	configFile := "config.ini" // Update with actual path
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to open config file: %s. Maybe try not screwing it up next time?\033[0m\n", err)
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
		fmt.Printf("\033[1;31mError reading config file: %s. Was it in the shredder?\033[0m\n", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to create Telegram bot: %s. Did you enter the token right, or are you messing with me?\033[0m\n", err)
		return
	}

	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		fmt.Printf("\033[1;31mInvalid chat ID: %s. Did you forget how to count?\033[0m\n", err)
		return
	}

	if results == "" {
		results = "No scan results to show. Maybe try a different hobby?"
	}

	msg := tgbotapi.NewMessage(chatIDInt, results)
	_, err = bot.Send(msg)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to send message to Telegram: %s. You sure the chat ID isn’t a black hole?\033[0m\n", err)
	}
}

func main() {
	printHeader()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter IP/domain to scan: \n> ")
	scanner.Scan()
	target := scanner.Text()

	fmt.Printf("\nScanning %s... Brace yourself, this is gonna be a wild ride!\n", target)
	fmt.Print("\nSelect scan type:\n")
	fmt.Print("  i. SYN-ACK Scan - Because poking the bear is fun\n")
	fmt.Print("  ii. UDP Scan - Unfiltered and full of chaos\n")
	fmt.Print("  iii. Version Scan - Because who doesn’t love details?\n")
	fmt.Print("  iv. OS Detection - Or should we say, ‘Guessing game’?\n")
	fmt.Print("> ")

	scanner.Scan()
	scanType := strings.TrimSpace(scanner.Text())

	fmt.Print("\nEnter additional arguments (or just hit Enter to skip):\n> ")
	scanner.Scan()
	args := scanner.Text()

	fmt.Print("\nEnter duration for progress indicator (seconds): \n> ")
	scanner.Scan()
	duration := atoi(scanner.Text())

	fmt.Println("\nFetching proxies...")
	proxies, err := scrapeProxies(proxyURL)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to fetch proxies: %s\033[0m\n", err)
		return
	}

	if err := saveProxies(outputFile, proxies); err != nil {
		fmt.Printf("\033[1;31mFailed to save proxies: %s\033[0m\n", err)
	} else {
		fmt.Printf("\033[1;32mProxies saved to %s\033[0m\n", outputFile)
	}

	fmt.Print("How would you like to handle results? (1. Save to file, 2. Send to Telegram, 3. Both): \n> ")
	scanner.Scan()
	resultHandling := scanner.Text()

	fmt.Println("\033[1;33mStarting scan. Hold tight!\033[0m")
	if strings.Contains(resultHandling, "1") {
		performScan(target, scanType, args, duration, proxies)
	}
	if strings.Contains(resultHandling, "2") {
		results := performScan(target, scanType, args, duration, proxies)
		if results != "" {
			sendResultsToTelegram(results)
		}
	}
	if strings.Contains(resultHandling, "3") {
		results := performScan(target, scanType, args, duration, proxies)
		if results != "" {
			sendResultsToTelegram(results)
		}
	}

	printFooter()
}
