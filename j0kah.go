package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	defaultScanDuration = 30  // Default duration for scan
	defaultScanType     = "SYN" // Default scan type
	defaultArgs         = "-T4 -A" // Default arguments for scan
	telegramTokenFile   = "config.ini" // File containing Telegram bot token and chat ID
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

func getUserInput(prompt string) string {
	fmt.Printf("\033[1;33m%s\033[0m: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func getTarget() string {
	return getUserInput("Enter the target to scan (e.g., example.com) - Or maybe you just want to watch the world burn")
}

func getScanType() string {
	return getUserInput("Enter scan type (e.g., SYN, ACK) - Because a regular scan just wouldn’t be enough")
}

func getConcurrency() int {
	for {
		input := getUserInput("Enter concurrency level (e.g., 80) - Or how many threads you can handle before you lose your sanity")
		concurrency, err := strconv.Atoi(input)
		if err == nil && concurrency > 0 {
			return concurrency
		}
		fmt.Println("\033[1;31mInvalid input. Enter a positive number. Do you even know how to count?\033[0m")
	}
}

func getScanDuration() int {
	for {
		input := getUserInput("Enter scan duration in seconds (default is 30) - Or just sit back and relax while it runs")
		duration, err := strconv.Atoi(input)
		if err == nil && duration > 0 {
			return duration
		}
		fmt.Println("\033[1;31mInvalid duration. Must be a positive number. Did you forget how to use a timer?\033[0m")
	}
}

func getSaveOption() bool {
	input := getUserInput("Do you want to save the scan results locally? (yes/no) - Or are you too lazy to bother?")
	return strings.ToLower(input) == "yes"
}

func getTelegramOption() bool {
	input := getUserInput("Do you want to send the scan results to Telegram? (yes/no) - Or would you rather keep it a secret?")
	return strings.ToLower(input) == "yes"
}

func performScan(target, scanType, args string, duration, concurrency int) string {
	fmt.Printf("\033[1;33mPreparing to perform a %s scan on %s with args '%s' and concurrency %d.\033[0m\n", scanType, target, args, concurrency)
	progressIndicator(duration)

	// Simulate scan result
	result := fmt.Sprintf("Simulated scan result for target: %s\nScan Type: %s\nDuration: %d seconds\nArgs: %s\nConcurrency: %d\n", target, scanType, duration, args, concurrency)
	fmt.Printf("\033[1;32mScan complete! Here are the results:\033[0m\n%s\n", result)

	return result
}

func sendResultsToTelegram(results string) {
	file, err := os.Open(telegramTokenFile)
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

	message := tgbotapi.NewMessage(chatIDInt, "Scan Results:\n"+results)
	_, err = bot.Send(message)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to send results to Telegram: %s. Did the bot get lost?\033[0m\n", err)
		return
	}

	fmt.Println("\033[1;32mResults successfully sent to Telegram.\033[0m")
}

func main() {
	printHeader()

	target := getTarget()
	scanType := getScanType()
	args := defaultArgs
	duration := getScanDuration()
	concurrency := getConcurrency()

	results := performScan(target, scanType, args, duration, concurrency)

	saveLocally := getSaveOption()
	if saveLocally {
		resultsFile := "scan_results.txt"
		err := os.WriteFile(resultsFile, []byte(results), 0644)
		if err != nil {
			fmt.Printf("\033[1;31mFailed to save scan results: %s. Maybe try not to mess things up next time?\033[0m\n", err)
			return
		}
		fmt.Println("\033[1;32mResults successfully saved locally.\033[0m")
	}

	sendToTelegram := getTelegramOption()
	if sendToTelegram {
		sendResultsToTelegram(results)
	}

	printFooter()
}
