package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	defaultScanDuration = 30 // Default duration for scan
	defaultScanType     = "SYN" // Default scan type
	defaultArgs         = "-T4 -A"  // Default arguments for scan
)

// printHeader prints the header of the tool.
func printHeader() {
	fmt.Println("\033[1;34m===================================\033[0m")
	fmt.Println("\033[1;34m       Welcome to j0kah Recon Tool\033[0m")
	fmt.Println("\033[1;34m===================================\033[0m")
	fmt.Println("\033[1;33mSelect the type of scan to perform, or just screw around:\033[0m")
	fmt.Println()
}

// printFooter prints the footer with creator info.
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

// progressIndicator shows the progress of the scan.
func progressIndicator(duration int) {
	for i := 0; i <= duration; i++ {
		time.Sleep(time.Second)
		fmt.Printf("\r\033[1;32mProgress: %d%% Complete. If you’re still here, congratulations, you’re officially a masochist.\033[0m", i*100/duration)
	}
	fmt.Println("\n\033[1;32mYou made it through the wait. Bravo, you’re now a certified saint. Or just really bored.\033[0m")
}

// performScan performs the scan with the given parameters.
func performScan(target, scanType, args string, duration int) (string, error) {
	fmt.Printf("\033[1;33mStarting %s scan on %s with arguments: %s\033[0m\n", scanType, target, args)

	var output []byte
	var err error

	cmd := exec.Command("nmap", append(strings.Split(args, " "), target)...)
	output, err = cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("\033[1;31mScan failed: %s. Retry? Of course, because who doesn't love a good, endless loop?\033[0m\n", err)
		return "", err
	}

	progressIndicator(duration)

	mainFile := "scan_results.txt"
	err = saveResults(mainFile, string(output))
	if err != nil {
		fmt.Printf("\033[1;31mFailed to save scan results: %s. Well, that’s just perfect, isn’t it?\033[0m\n", err)
		return "", err
	}

	fmt.Println("\n\033[1;33mScan Results:\033[0m")
	fmt.Printf("\033[1;33mTarget:\033[0m %s\n", target)
	fmt.Printf("\033[1;33mFiltered Output:\033[0m\n%s\n", string(output))

	return mainFile, nil
}

// saveResults saves the scan results to a file.
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

// getUserInput prompts the user for various inputs and returns them.
func getUserInput() (string, string, string, int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter target (e.g., example.com): ")
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)

	fmt.Print("Enter scan type (e.g., SYN, UDP): ")
	scanType, _ := reader.ReadString('\n')
	scanType = strings.TrimSpace(scanType)
	if scanType == "" {
		scanType = defaultScanType
	}

	fmt.Print("Enter scan arguments (e.g., -T4 -A): ")
	args, _ := reader.ReadString('\n')
	args = strings.TrimSpace(args)
	if args == "" {
		args = defaultArgs
	}

	fmt.Print("Enter scan duration in seconds (e.g., 30): ")
	durationStr, _ := reader.ReadString('\n')
	durationStr = strings.TrimSpace(durationStr)
	duration := defaultScanDuration
	if durationStr != "" {
		duration, _ = strconv.Atoi(durationStr)
	}

	return target, scanType, args, duration
}

// main is the entry point of the application.
func main() {
	printHeader()

	target, scanType, args, duration := getUserInput()

	_, err := performScan(target, scanType, args, duration)
	if err != nil {
		fmt.Printf("\033[1;31mScan failed: %s\033[0m\n", err)
	}

	printFooter()
}
