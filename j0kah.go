package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	maxConcurrency      = 10
	maxRetries          = 3
	retryDelay          = 2 * time.Second
	defaultScanDuration = 30 // Default duration for scan
	defaultScanType     = "SYN" // Default scan type
	defaultArgs         = "-T4 -A"  // Default arguments for scan
)

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

func performScan(target, scanType, args string, duration int) (string, string) {
	var output []byte
	var err error
	for i := 0; i < maxRetries; i++ {
		cmd := exec.Command("nmap", append(strings.Split(args, " "), target)...)
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

	progressIndicator(duration)

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
		fmt.Printf("\033[1;31mFailed to send results to Telegram: %s. Did the bot get lost?\033[0m\n", err)
		return
	}

	fmt.Println("\033[1;32mResults successfully sent to Telegram.\033[0m")
}

func main() {
	printHeader()

	// Prompt for scan parameters
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\033[1;33mAlright, you brave soul. Enter the target (e.g., example.com) or face my wrath: \033[0m")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)

	fmt.Print("\033[1;33mSelect your scan type (e.g., SYN). Choose wisely, or prepare for chaos: \033[0m")
	scanType, _ := reader.ReadString('\n')
	scanType = strings.TrimSpace(scanType)

	fmt.Print("\033[1;33mHow long should we endure this pain? Enter scan duration in seconds (default 30): \033[0m")
	durationStr, _ := reader.ReadString('\n')
	durationStr = strings.TrimSpace(durationStr)
	duration := defaultScanDuration
	if durationStr != "" {
		duration = atoi(durationStr)
	}

	fmt.Print("\033[1;33mAdditional scan arguments (default '-T4 -A'). Make it spicy: \033[0m")
	args, _ := reader.ReadString('\n')
	args = strings.TrimSpace(args)
	if args == "" {
		args = defaultArgs
	}

	fmt.Print("\033[1;33mHow many warriors should we send into the fray? Enter concurrency level (default 10): \033[0m")
	concurrencyStr, _ := reader.ReadString('\n')
	concurrencyStr = strings.TrimSpace(concurrencyStr)
	if concurrencyStr == "" {
		concurrencyStr = strconv.Itoa(maxConcurrency)
	}
	maxConcurrency = atoi(concurrencyStr)

	fmt.Print("\033[1;33mRetries, retries, and more retries! Enter number of retries (default 3): \033[0m")
	retriesStr, _ := reader.ReadString('\n')
	retriesStr = strings.TrimSpace(retriesStr)
	if retriesStr == "" {
		retriesStr = strconv.Itoa(maxRetries)
	}
	maxRetries = atoi(retriesStr)

	fmt.Print("\033[1;33mHow long should we wait before another try? Enter retry delay in seconds (default 2): \033[0m")
	retryDelayStr, _ := reader.ReadString('\n')
	retryDelayStr = strings.TrimSpace(retryDelayStr)
	if retryDelayStr != "" {
		retryDelay = time.Duration(atoi(retryDelayStr)) * time.Second
	}

	// Perform the scan
	mainFile, unknownFile := performScan(target, scanType, args, duration)
	if mainFile != "" {
		fmt.Printf("\033[1;33mScan results saved to: %s\033[0m\n", mainFile)
		fmt.Printf("\033[1;33mUnknown ports saved to: %s\033[0m\n", unknownFile)
	}

	// Send results to Telegram
	sendResultsToTelegram(mainFile)

	printFooter()
}
