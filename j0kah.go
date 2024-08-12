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

func getUserInput(prompt string, defaultValue string) string {
	fmt.Printf("\033[1;33m%s (default: %s)\033[0m: ", prompt, defaultValue)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return defaultValue
	}
	return input
}

func getTargets() string {
	return getUserInput("Enter comma-separated list of targets (default: target1,target2) - Or just put a single target", defaultTargets)
}

func getScanType() string {
	return getUserInput("Enter scan type (SYN/UDP/TCP/ACK/Xmas/Null/FIN/Window/Maimon) - Or stick with the default SYN", defaultScanType)
}

func getConcurrency() int {
	for {
		input := getUserInput("Enter concurrency level (default: 10) - Or how many threads you can handle before you lose your sanity", "10")
		concurrency, err := strconv.Atoi(input)
		if err == nil && concurrency > 0 {
			return concurrency
		}
		fmt.Println("\033[1;31mInvalid input. Enter a positive number. Do you even know how to count?\033[0m")
	}
}

func getScanDuration() int {
	for {
		input := getUserInput("Enter scan duration in seconds (default: 30) - Or just sit back and relax while it runs", strconv.Itoa(defaultScanDuration))
		duration, err := strconv.Atoi(input)
		if err == nil && duration > 0 {
			return duration
		}
		fmt.Println("\033[1;31mInvalid duration. Must be a positive number. Did you forget how to use a timer?\033[0m")
	}
}

func getSaveOption() bool {
	input := getUserInput("Do you want to save the scan results locally? (y/n) - Or are you too lazy to bother?", "n")
	return strings.ToLower(input) == "y"
}

func getTelegramOption() bool {
	input := getUserInput("Do you want to send the scan results to Telegram? (y/n) - Or would you rather keep it a secret?", "n")
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
	err := cmd.Run()
	if err != nil {
		fmt.Printf("\033[1;31mFailed to execute scan: %s. Check your command and targets.\033[0m\n", err)
		return ""
	}

	scanResults := parseScanResults(out.String())

	fmt.Printf("\033[1;32mScan complete! Results will be sent to Telegram or saved locally based on your choice.\033[0m\n")

	return scanResults
}

func parseScanResults(output string) string {
	var result string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "open") {
			parts := strings.Fields(line)
			port := parts[0]
			service := strings.Join(parts[1:], " ")
			result += fmt.Sprintf("  - Port %s (tcp): open (Service: %s)\n", port, service)
		} else if strings.Contains(line, "closed") {
			parts := strings.Fields(line)
			port := parts[0]
			service := strings.Join(parts[1:], " ")
			result += fmt.Sprintf("  - Port %s (tcp): closed (Service: %s)\n", port, service)
		}
	}
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
			chatID = strings.TrimSpace(strings.TrimPrefix(line, "chat_id="))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("\033[1;31mError reading config file: %s. Was it in the shredder?\033[0m\n", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to create Telegram bot: %s. Did you enter the token right, or are you fudging with me?\033[0m\n", err)
		return
	}
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		fmt.Printf("\033[1;31mInvalid chat ID: %s. Seriously? Did you just pull that out of a rat's arse?\033[0m\n", err)
		return
	}

	// Send the "Mail's here!" message to Telegram
	mailMessage := tgbotapi.NewMessage(chatIDInt, "Hellooo stranger! Your mail's here! You wanted chaos? You got it now have fun!")
	_, err = bot.Send(mailMessage)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to send 'Mail's here!' message to Telegram: %s. Did the pigeons get lost?\033[0m\n", err)
		return
	}

	// Now save the results to a file and send the file to Telegram
	fileName := "scan_results.txt"
	err = os.WriteFile(fileName, []byte(results), 0644)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to save scan results to file: %s. Did you forget to write it?\033[0m\n", err)
		return
	}

	// Send the results file to Telegram
	fileToSend := tgbotapi.NewDocumentUpload(chatIDInt, fileName)
	_, err = bot.Send(fileToSend)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to send results file to Telegram: %s. Did you lose the file in the void?\033[0m\n", err)
		return
	}

	fmt.Println("\033[1;32mResults have been delivered to Telegram. Brace yourself—because you just invited chaos into your chat.\033[0m")
}

func saveResultsLocally(results string) {
	fileName := getUserInput("Enter file name to save results (default: scan_results.txt)", "scan_results.txt")
	err := os.WriteFile(fileName, []byte(results), 0644)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to save results to file: %s. You probably messed something up. \033[0m\n", err)
		return
	}
	fmt.Printf("\033[1;32mResults saved to file: %s\033[0m\n", fileName)
}

func main() {
	printHeader()

	// Collect scan parameters from the user
	targets := getTargets()
	scanType := getScanType()
	args := getArgs()
	duration := getScanDuration()
	concurrency := getConcurrency()

	// Perform the scan
	results := performScan(targets, scanType, args, duration, concurrency)

	// Check if the user wants to send results to Telegram
	if getTelegramOption() {
		sendResultsToTelegram(results)
	} else if getSaveOption() {
		saveResultsLocally(results)
	}

	printFooter()
}
