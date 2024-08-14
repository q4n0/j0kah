package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	defaultScanDuration = 30                // Default duration for scan
	defaultScanType     = "SYN"             // Default scan type
	defaultArgs         = "-T4 -A"          // Default arguments for scan
	defaultTargets      = "target1,target2" // Default targets
	telegramTokenFile   = "config.ini"      // File containing Telegram bot token and chat ID
)

func printHeader() {
	fmt.Println("\033[1;31m=========================================================================\033[0m")
	fmt.Println("\033[1;31m       Welcome to the Insanity Show, Where Chaos Reigns Supreme,\033[0m")
	fmt.Println("\033[1;31m       And you're the Star of the Act\033[0m")
	fmt.Println("\033[1;31m           b0urn3 - The Mad Genius\033[0m")
	fmt.Println("\033[1;31m==========================================================================\033[0m")
	fmt.Println("\033[1;33mReady to dance with madness? Let the games begin.\033[0m")
	fmt.Println()
}

func printFooter() {
	fmt.Println()
	fmt.Println("\033[1;31m=================================================================\033[0m")
	fmt.Println("\033[1;32mCreated by b0urn3\033[0m")
	fmt.Println("\033[1;32mGitHub: \033[1;36mhttps://github.com/q4n0\033[0m")
	fmt.Println("\033[1;32mInstagram: \033[1;36mhttps://www.instagram.com/onlybyhive\033[0m")
	fmt.Println("\033[1;32mEmail: \033[1;36mb0urn3@proton.me\033[0m")
	fmt.Println("\033[1;31m==================================================================\033[0m")
	fmt.Println()
}

func progressIndicator(duration int) {
	for i := 0; i <= duration; i++ {
		time.Sleep(time.Second)
		fmt.Printf("\033[1;32mProgress: %d%% Complete. If you’re still here, you must enjoy torment. Either that or you’ve mastered the art of waiting. \033[0m\r", i*100/duration)
	}
	fmt.Println("\033[1;32mWow, you survived! Now go ahead and reward yourself for this grueling achievement. Or just admit you have nothing better to do. \033[0m")
}

func getUserInput(prompt, defaultValue string) string {
	fmt.Printf("\033[1;33m%s (default: %s)\033[0m: ", prompt, defaultValue)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return defaultValue
	}
	return input
}

func getTargets() string {
	return getUserInput("Enter comma-separated list of targets (default: target1,target2)", defaultTargets)
}

func getScanType() string {
	return getUserInput("Enter scan type (SYN/UDP/TCP/ACK/Xmas/Null/FIN/Window/Maimon)", defaultScanType)
}

func getConcurrency() int {
	for {
		input := getUserInput("Enter concurrency level (default: 10)", "10")
		concurrency, err := strconv.Atoi(input)
		if err == nil && concurrency > 0 {
			return concurrency
		}
		fmt.Println("\033[1;31mInvalid input. Enter a positive number.\033[0m")
	}
}

func getScanDuration() int {
	for {
		input := getUserInput("Enter scan duration in seconds (default: 30)", strconv.Itoa(defaultScanDuration))
		duration, err := strconv.Atoi(input)
		if err == nil && duration > 0 {
			return duration
		}
		fmt.Println("\033[1;31mInvalid duration. Must be a positive number.\033[0m")
	}
}

func getSaveOption() bool {
	input := getUserInput("Do you want to save the scan results locally? (y/n)", "n")
	return strings.ToLower(input) == "y"
}

func getTelegramOption() bool {
	input := getUserInput("Do you want to send the scan results to Telegram? (y/n)", "n")
	return strings.ToLower(input) == "y"
}

func getArgs() string {
	return getUserInput("Enter additional arguments for scan (default: -T4 -A)", defaultArgs)
}

func performScan(targets, scanType, args string, duration, concurrency int) string {
	fmt.Printf("\033[1;33mPreparing to perform a %s scan on %s with args '%s' and concurrency %d.\033[0m\n", scanType, targets, args, concurrency)
	progressIndicator(duration)

	cmd := exec.Command("nmap", args, "-p-", targets)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("\033[1;31mFailed to execute scan: %s\033[0m\n", err)
		return ""
	}

	return parseScanResults(out.String())
}

func parseScanResults(output string) string {
	var result strings.Builder
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "open") || strings.Contains(line, "closed") {
			parts := strings.Fields(line)
			port := parts[0]
			service := strings.Join(parts[1:], " ")
			status := "open"
			if strings.Contains(line, "closed") {
				status = "closed"
			}
			result.WriteString(fmt.Sprintf("  - Port %s (tcp): %s (Service: %s)\n", port, status, service))
		}
	}
	return result.String()
}

func sendResultsToTelegram(results string) {
	file, err := os.Open(telegramTokenFile)
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
			chatID = strings.TrimSpace(strings.TrimPrefix(line, "chat_id="))
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

	mailMessage := tgbotapi.NewMessage(chatIDInt, "Hellooo stranger! Your mail's here! You wanted chaos? You got it now have fun!")
	if _, err := bot.Send(mailMessage); err != nil {
		fmt.Printf("\033[1;31mFailed to send 'Mail's here!' message to Telegram: %s\033[0m\n", err)
		return
	}

	fileName := "scan_results.txt"
	if err := os.WriteFile(fileName, []byte(results), 0644); err != nil {
		fmt.Printf("\033[1;31mFailed to save scan results to file: %s\033[0m\n", err)
		return
	}

	fileToSend := tgbotapi.NewDocumentUpload(chatIDInt, fileName)
	if _, err := bot.Send(fileToSend); err != nil {
		fmt.Printf("\033[1;31mFailed to send scan results to Telegram: %s\033[0m\n", err)
		return
	}

	fmt.Println("\033[1;32mScan results have been sent to Telegram.\033[0m")
}

func main() {
	printHeader()

	targets := getTargets()
	scanType := getScanType()
	concurrency := getConcurrency()
	duration := getScanDuration()
	saveResults := getSaveOption()
	telegramOption := getTelegramOption()
	args := getArgs()

	results := performScan(targets, scanType, args, duration, concurrency)

	if results != "" {
		if saveResults {
			err := os.WriteFile("scan_results.txt", []byte(results), 0644)
			if err != nil {
				fmt.Printf("\033[1;31mFailed to save scan results: %s\033[0m\n", err)
			} else {
				fmt.Println("\033[1;32mScan results have been saved locally.\033[0m")
			}
		}

		if telegramOption {
			sendResultsToTelegram(results)
		}
	}

	printFooter()
}
