package main

import (
    "bytes"
    "fmt"
    "log"
    "os"
    "os/exec"
    "sync"
)

// Enhanced error handling for performing the scan
func performScan(targets, scanType, args string, duration, concurrency int) string {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var results string

    fmt.Printf("\033[1;33mPreparing to perform a %s scan on %s with args '%s' and concurrency %d.\033[0m\n", scanType, targets, args, concurrency)
    progressIndicator(duration)

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            cmdArgs := []string{args, "-p-", targets}
            
            // Support additional scan types
            switch scanType {
            case "SYN":
                cmdArgs = append(cmdArgs, "-sS")
            case "UDP":
                cmdArgs = append(cmdArgs, "-sU")
            case "TCP":
                cmdArgs = append(cmdArgs, "-sT")
            case "ACK":
                cmdArgs = append(cmdArgs, "-sA")
            case "Xmas":
                cmdArgs = append(cmdArgs, "-sX")
            case "Null":
                cmdArgs = append(cmdArgs, "-sN")
            case "FIN":
                cmdArgs = append(cmdArgs, "-sF")
            case "Window":
                cmdArgs = append(cmdArgs, "-sW")
            case "Maimon":
                cmdArgs = append(cmdArgs, "-sM")
            default:
                cmdArgs = append(cmdArgs, "-sS") // Default to SYN scan
            }

            cmd := exec.Command("nmap", cmdArgs...)
            var out bytes.Buffer
            cmd.Stdout = &out
            err := cmd.Run()
            if err != nil {
                fmt.Printf("\033[1;31mFailed to execute scan: %s. Check your command and targets.\033[0m\n", err)
                return
            }

            mu.Lock()
            results += parseScanResults(out.String())
            mu.Unlock()
        }(i)
    }

    wg.Wait()

    fmt.Printf("\033[1;32mScan complete! Results will be sent to Telegram or saved locally based on your choice.\033[0m\n")

    return results
}

// Improved function to save results locally with error handling
func saveResultsLocally(results string) {
    fileName := getUserInput("Enter file name to save results (default: scan_results.txt)", "scan_results.txt")
    err := os.WriteFile(fileName, []byte(results), 0644)
    if err != nil {
        fmt.Printf("\033[1;31mFailed to save results to file: %s. You probably messed something up.\033[0m\n", err)
        return
    }
    fmt.Printf("\033[1;32mResults saved to file: %s\033[0m\n", fileName)
}

// Enhanced user input function with better feedback
func getUserInput(prompt string, defaultValue string) string {
    fmt.Printf("\033[1;33m%s (default: %s)\033[0m: ", prompt, defaultValue)
    var input string
    _, err := fmt.Scanln(&input)
    if err != nil {
        fmt.Printf("\033[1;31mError reading input: %s. Using default value.\033[0m\n", err)
        return defaultValue
    }
    if input == "" {
        return defaultValue
    }
    return input
}

// Simple progress indicator for simulating scan progress
func progressIndicator(duration int) {
    fmt.Printf("\033[1;36m[+] Scanning...\033[0m\n")
    for i := 0; i < duration; i++ {
        fmt.Printf(".")
        time.Sleep(time.Second)
    }
    fmt.Printf("\n")
}

// Function to parse scan results (dummy implementation)
func parseScanResults(output string) string {
    // Implement the logic to parse the Nmap output
    return output // For now, just return the raw output
}

// Improved logging function
func logEvent(message string) {
    logFile, err := os.OpenFile("j0kah.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("\033[1;31mFailed to open log file: %s\033[0m\n", err)
        return
    }
    defer logFile.Close()
    logger := log.New(logFile, "", log.LstdFlags)
    logger.Println(message)
}

// Function to integrate with other tools (placeholder)
func integrateWithOtherTools(results string) {
    // Implement integration logic with other tools or platforms
    // For example, send results to an external API or tool
}

func main() {
    // Example usage of the tool
    targets := getUserInput("Enter the target(s) for the scan", "192.168.1.0/24")
    scanType := getScanType()
    args := "-Pn"
    duration := 10
    concurrency := 5

    results := performScan(targets, scanType, args, duration, concurrency)

    // Optional: Save results locally or integrate with other tools
    saveResultsLocally(results)
    integrateWithOtherTools(results)
}
