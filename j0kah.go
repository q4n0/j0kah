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
		fmt.Printf("\033[1;32mProgress: %d%% Complete. If you’re still here, congratulations, you’re officially a masochist.\033[0m\r", i*100/duration)
	}
	fmt.Println("\033[1;32mYou made it through the wait. Bravo, you’re now a certified saint. Or just really bored.\033[0m")
}

func performScan(target, scanType, args string, duration int, proxies []string) (string, string) {
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
		return "", ""
	}

	wg.Wait()

	filteredOutput := filterOutput(string(output))
	unknownPorts := filterUnknownPorts(string(output))

	mainFile := "scan_results.txt"
	unknownFile := "unknown_ports.txt"

	err = saveResults(mainFile, filteredOutput)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to save scan results: %s\033[0m\n", err)
	}

	err = saveResults(unknownFile, unknownPorts)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to save unknown ports: %s\033[0m\n", err)
	}

	fmt.Println("\n\033[1;33mScan Results:\033[0m")
	fmt.Printf("\033[1;33mTarget:\033[0m %s\n", target)
	fmt.Printf("\033[1;33mFiltered Output:\033[0m\n%s\n", filteredOutput)

	return mainFile, unknownFile
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
		fmt.Printf("\033[1;31mError converting string to int: %s\033[0m\n", err)
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
		fmt.Printf("\033[1;31mFailed to create Telegram bot: %s. Did you enter the token right, or are you messing with me?\033[0m\n", err)
		return
	}

	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		fmt.Printf("\033[1;31mInvalid chat ID: %s. Did you forget how to count?\033[0m\n", err)
		return
	}

	fileToSend := tgbotapi.NewDocumentUpload(chatIDInt, resultsFile)
	_, err = bot.Send(fileToSend)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to send file to Telegram: %s. You sure the chat ID isn’t a black hole?\033[0m\n", err)
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
	fmt.Print("  iii. AnonScan - Sneaky like a thief in the night\n")
	fmt.Print("  iv. Regular Scan - The vanilla flavor for the boring folks\n")
	fmt.Print("  v. OS Detection - Guessing what OS they're running, like a pro\n")
	fmt.Print("  vi. Multiple IP inputs - Because one target is never enough\n")
	fmt.Print("  vii. Ping Scan - Hello? Is anybody home?\n")
	fmt.Print("  viii. Comprehensive Scan - The whole shebang, go big or go home\n> ")
	scanner.Scan()
	scanType := scanner.Text()

	scanOptions := map[string]string{
		"i":    "-sS",
		"ii":   "-sU",
		"iii":  "-sS -Pn",
		"iv":   "-sT",
		"v":    "-O",
		"vi":   "-iL",
		"vii":  "-sn",
		"viii": "-A",
	}

	args, ok := scanOptions[scanType]
	if !ok {
		fmt.Println("\033[1;31mInvalid scan type selected. Try again if you’re feeling lucky!\033[0m")
		return
	}

	fmt.Print("Enter duration for progress indicator (seconds): \n> ")
	scanner.Scan()
	duration := atoi(scanner.Text())

	fmt.Print("Do you want to use proxies? (y/n): \n> ")
	scanner.Scan()
	useProxies := strings.ToLower(scanner.Text()) == "y"

	var proxies []string
	if useProxies {
		fmt.Println("\033[1;33mFetching proxies... or not, depending on how the internet feels today.\033[0m")
		proxies, _ = scrapeProxies(proxyURL)
		if err := saveProxies(outputFile, proxies); err != nil {
			fmt.Printf("\033[1;31mFailed to save proxies: %s\033[0m\n", err)
		} else {
			fmt.Printf("\033[1;32mProxies saved to %s. Because anonymity is a thing.\033[0m\n", outputFile)
		}
	}

	fmt.Print("How would you like to handle results? (1. Save to file, 2. Send to Telegram, 3. Both): \n> ")
	scanner.Scan()
	resultHandling := scanner.Text()

	fmt.Println("\033[1;33mStarting scan. Hold tight!\033[0m")
	if strings.Contains(resultHandling, "1") {
		mainFile, unknownFile := performScan(target, scanType, args, duration, proxies)
		if mainFile != "" {
			fmt.Printf("\033[1;32mScan results saved to %s and %s.\033[0m\n", mainFile, unknownFile)
		}
	}
	if strings.Contains(resultHandling, "2") {
		mainFile, _ := performScan(target, scanType, args, duration, proxies)
		if mainFile != "" {
			sendResultsToTelegram(mainFile)
		}
	}
	if strings.Contains(resultHandling, "3") {
		mainFile, unknownFile := performScan(target, scanType, args, duration, proxies)
		if mainFile != "" {
			sendResultsToTelegram(mainFile)
		}
		if unknownFile != "" {
			sendResultsToTelegram(unknownFile)
		}
	}

	printFooter()
}
