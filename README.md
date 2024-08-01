j0kah Recon Tool

Welcome to the j0kah Recon Tool, a versatile network scanning utility with optional proxy support and Telegram integration. Follow the steps below to set up and use the tool effectively.
Table of Contents

    Prerequisites
    Setup
    Usage
    Proxy Management
    Telegram Integration
    Troubleshooting

Prerequisites

Before you begin, ensure you have the following:

    Go Language: Ensure Go is installed on your system. You can download it from golang.org.
    Nmap: Install Nmap for network scanning. You can get it from nmap.org.
    Telegram Bot Token: Create a Telegram bot and get the token from BotFather.
    Config File: Create a config.ini file with the following structure:

    makefile

    token=YOUR_TELEGRAM_BOT_TOKEN
    chat_id=YOUR_CHAT_ID

Setup

    Clone the Repository:

    bash

git clone https://github.com/yourusername/j0kah
cd j0kah

Build the Project:

Navigate to the project directory and build the Go application:

bash

go build -o j0kah

Create Proxy Scraper:

Save the proxy scraper script as proxy_scraper.go:

go

// Paste the proxy scraper code here

Build the proxy scraper:

bash

    go build -o proxy_scraper proxy_scraper.go

Usage

    Run the Proxy Scraper:

    To fetch and save proxies:

    bash

./proxy_scraper

This will save the fetched proxies to proxy.list.

Run j0kah Recon Tool:

Start the tool with:

bash

    ./j0kah

    Follow the on-screen prompts to:
        Enter IP/domain to scan.
        Choose scan type (e.g., SYN-ACK Scan, UDP Scan).
        Decide whether to use proxies.
        Enter proxy settings if applicable.
        Enter the scan duration and other options.
        Choose how to handle results (save to file, send via Telegram, or both).

Proxy Management

If you wish to use proxies for scanning, ensure:

    Run Proxy Scraper:

    Fetch and save the list of proxies:

    bash

    ./proxy_scraper

    Configure Proxies:

    When prompted by j0kah, provide the path to the proxy.list file.

Telegram Integration

To receive scan results on Telegram:

    Create a Telegram Bot:
        Use BotFather on Telegram to create a bot and get the token.

    Configure Bot:
        Save the bot token and your chat ID in config.ini.

    Send Results:
        During tool usage, choose the option to send results to Telegram.

Troubleshooting

    Error: unexpected status code: 404: This may indicate an issue with the proxy list URL. Verify the URL and ensure it's correct.
    Failed to create file: Ensure you have the necessary permissions to create files in the directory.
    Telegram errors: Check your config.ini file for correct token and chat ID.
