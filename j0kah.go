package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	maxConcurrency = 80
	maxRetries     = 3
	retryDelay     = 2 * time.Second
	outputFile     = "proxy.list"
	testURL        = "http://www.google.com"
)

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
	fmt.Println("\033[1;34m       You merely adopted the scans.\033[0m")
	fmt.Println("\033[1;34m       I was born in it, molded by it.\033[0m")
	fmt.Println("\033[1;34m===================================\033[0m")
	fmt.Println("\033[1;33mWelcome to the j0kah Recon Tool, the only tool worthy of your time.\033[0m")
	fmt.Println()
}

func printFooter() {
	fmt.Println()
	fmt.Println("\033[1;34m===================================\033[0m")
	fmt.Println("\033[1;32mWhen the scan is complete, you have my permission to gloat.\033[0m")
	fmt.Println("\033[1;32mCreated by bo0urn3, a master of the craft.\033[0m")
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

func performScan(target, scanType, args string, duration int) string {
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

	mainFile := "scan_results.txt"

	err = saveResults(mainFile, filteredOutput)
	if err != nil {
		fmt.Printf("\033[1;31mFailed to save scan results: %s\033[0m\n", err)
	}

	return mainFile
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
	if fileInfo, err := os.Stat(resultsFile); err != nil || fileInfo.Size() == 0 {
		fmt.Printf("\033[1;31mFile %s is empty or does not exist. Not sending to Telegram.\033[0m\n", resultsFile)
		return
	}

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
		fmt.Printf("\033[1;31mFailed to create Telegram bot: %s. Did you enter the token correctly, genius?\033[0m\n", err)
		return
	}

	chatIDInt64, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		fmt.Printf("\033[1;31mInvalid chat ID: %s. Can’t send messages to a black hole, you know.\033[0m\n", chatID)
		return
	}

	message := tgbotapi.NewMessage(chatIDInt64, fmt.Sprintf("Scan results saved to %s", resultsFile))
	if _, err := bot.Send(message); err != nil {
		fmt.Printf("\033[1;31mFailed to send Telegram message: %s. Guess the bots are revolting.\033[0m\n", err)
	}
}

func getTarget() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the target URL or IP: ")
	target, _ := reader.ReadString('\n')
	return strings.TrimSpace(target)
}

func main() {
	printHeader()
	fmt.Println("i. SYN-ACK Scan - Because poking the bear is fun")
	fmt.Println("ii. UDP Scan - Unfiltered and full of chaos")
	fmt.Println("iii. AnonScan - Sneaky like a thief in the night")
	fmt.Println("iv. Regular Scan - The vanilla flavor for the boring folks")
	fmt.Println("v. OS Detection - Guessing what OS they're running, like a pro")
	fmt.Println("vi. Multiple IP inputs - Because one target is never enough")
	fmt.Println("vii. Ping Scan - Hello? Is anybody home?")
	fmt.Println("viii. Comprehensive Scan - The whole shebang, go big or go home")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your choice: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var scanType, args string
	var duration int
	switch choice {
	case "i":
		scanType = "SYN-ACK"
		args = "-sS"
		duration = 60
	case "ii":
		scanType = "UDP"
		args = "-sU"
		duration = 120
	case "iii":
		scanType = "AnonScan"
		args = "-sS -sU -O -A"
		duration = 180
	case "iv":
		scanType = "Regular"
		args = "-A"
		duration = 60
	case "v":
		scanType = "OS Detection"
		args = "-O"
		duration = 90
	case "vi":
		scanType = "Multiple IP inputs"
		args = "-sS -sU"
		duration = 120
	case "vii":
		scanType = "Ping Scan"
		args = "-sn"
		duration = 10
	case "viii":
		scanType = "Comprehensive"
		args = "-sS -sU -O -A -T4"
		duration = 240
	default:
		fmt.Println("\033[1;31mInvalid choice. Please follow the script, it's not that hard.\033[0m")
		return
	}

	target := getTarget()

	mainFile := performScan(target, scanType, args, duration)
	if mainFile == "" {
		return
	}

	fmt.Println("\033[1;32mScan results saved to:", mainFile, "\033[0m")

	fmt.Print("Would you like to send the results to Telegram? (yes/no): ")
	sendToTelegram, _ := reader.ReadString('\n')
	sendToTelegram = strings.TrimSpace(sendToTelegram)

	if strings.ToLower(sendToTelegram) == "yes" {
		sendResultsToTelegram(mainFile)
	} else {
		fmt.Println("\033[1;33mAlright, results not sent to Telegram. Your secrets are safe. For now.\033[0m")
	}

	printFooter()
}
