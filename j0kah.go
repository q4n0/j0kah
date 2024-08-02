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
	maxConcurrency        = 10
	maxRetries            = 3
	retryDelay            = 2 * time.Second
	proxyURL              = "https://www.proxy-list.download/api/v1/get?type=https"
	outputFile            = "proxy.list"
	defaultScanDuration   = 30 // Default duration for scan
	defaultScanType       = "SYN" // Default scan type
	defaultArgs           = "-T4 -A" // Default arguments for scan
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

func performScan(target, scanType, args string, duration int, proxies []string, concurrency int) (string, string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		progressIndicator(duration)
	}()

	sem := make(chan struct{}, concurrency)
	var output []byte
	var err error
	for i := 0; i < maxRetries; i++ {
		sem <- struct{}{}
		cmd := exec.Command("nmap", append(strings.Split(args, " "), target)...)
		if len(proxies) > 0 {
			cmd.Env = append(os.Environ(), "http_proxy=http://"+proxies[0])
		}
		output, err = cmd.CombinedOutput()
		<-sem

		if err == nil {
			break
		}
		fmt.Printf("\033[1;31mScan failed: %s. Retry? Of course, because who doesn't love a good, endless loop?\033[0m\n", err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		fmt.Printf("\033[1;31mFinal scan error: %s\nError output: %s\033[0m\n", err, string(output))
		return "", ""
	}

	wg.Wait()

	filteredOutput := filterOutput(string(output))
	unknownPorts := filterUnknownPorts(string(output))

	mainFile := "scan_results.txt"
	unknownFile := "unknown_ports.txt"

	err = saveResults(mainFile, filteredOutput)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to save scan results: %s. Well, that’s just perfect, isn’t it?\033[0m\n", err)
	}

	err = saveResults(unknownFile, unknownPorts)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to save unknown ports: %s. Maybe try not to mess things up next time?\033[0m\n", err)
	}

	fmt.Println("\n\033[1;33mScan Results:\033[0m")
	fmt.Printf("\033[1;33mTarget:\033[0m %s\n", target)
	fmt.Printf("\033[1;33mFiltered Output:\033[0m\n%s\n", filteredOutput)

	return mainFile, unknownFile
}

func filterOutput(output string) string {
	lines := strings.Split(output, "\n")
	var filteredLines []string
	for _, line := lines {
		if strings.Contains(line, "/tcp") || strings.Contains(line, "/udp") || strings.Contains(line, "VULNERABILITY") {
			filteredLines = append(filteredLines, line)
		}
	}
	return strings.Join(filteredLines, "\n")
}

func filterUnknownPorts(output string) string {
	lines := strings.Split(output, "\n")
	var unknownPorts []string
	for _, line := range lines {
		if !strings.Contains(line, "/tcp") && !strings.Contains(line, "/udp") {
			unknownPorts = append(unknownPorts, line)
		}
	}
	return strings.Join(unknownPorts, "\n")
}

func saveResults(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(content); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing buffer: %w", err)
	}

	return nil
}

func atoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("\033[1;31mError converting string to int: %s. Did you forget how to count?\033[0m\n", err)
		return 0
	}
	return val
}

func sendResultsToTelegram(resultsFile string) {
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
		fmt.Printf("\033[1;31mFailed to create Telegram bot: %s. Did you really think this would work?\033[0m\n", err)
		return
	}

	msg := tgbotapi.NewMessage(atoi(chatID), "Scan results:")
	_, err = bot.Send(msg)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to send Telegram message: %s. Maybe the bot is just shy?\033[0m\n", err)
		return
	}

	fileContent, err := os.ReadFile(resultsFile)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to read results file: %s. Did the file just vanish into thin air?\033[0m\n", err)
		return
	}

	msg = tgbotapi.NewMessage(atoi(chatID), string(fileContent))
	_, err = bot.Send(msg)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to send Telegram message: %s. Did you think it would be that easy?\033[0m\n", err)
	}
}

func main() {
	printHeader()

	proxies, err := scrapeProxies(proxyURL)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to scrape proxies: %s. Better luck next time?\033[0m\n", err)
		return
	}

	err = saveProxies(outputFile, proxies)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to save proxies: %s. Proxies? What proxies?\033[0m\n", err)
		return
	}

	var target, scanType, args string
	var duration, concurrency int

	fmt.Print("\033[1;33mEnter the target (IP or URL): \033[0m")
	fmt.Scanln(&target)
	if target == "" {
		fmt.Printf("\033[1;31mInvalid target. You had one job.\033[0m\n")
		return
	}

	fmt.Print("\033[1;33mEnter the scan duration (seconds): \033[0m")
	fmt.Scanln(&duration)
	if duration == 0 {
		duration = defaultScanDuration
	}

	fmt.Print("\033[1;33mEnter the scan type (SYN, UDP, etc.): \033[0m")
	fmt.Scanln(&scanType)
	if scanType == "" {
		scanType = defaultScanType
	}

	fmt.Print("\033[1;33mEnter any additional arguments (e.g., -T4 -A -v): \033[0m")
	fmt.Scanln(&args)
	if args == "" {
		args = defaultArgs
	}

	fmt.Print("\033[1;33mEnter the concurrency level: \033[0m")
	fmt.Scanln(&concurrency)
	if concurrency == 0 {
		concurrency = maxConcurrency
	}

	resultsFile, unknownFile := performScan(target, scanType, args, duration, proxies, concurrency)

	fmt.Print("\033[1;33mDo you want to send results to Telegram? (yes/no): \033[0m")
	var sendToTelegram string
	fmt.Scanln(&sendToTelegram)

	if strings.ToLower(sendToTelegram) == "yes" {
		sendResultsToTelegram(resultsFile)
		sendResultsToTelegram(unknownFile)
	}

	printFooter()
}
