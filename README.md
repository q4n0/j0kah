j0kah Recon Tool

The j0kah Recon Tool is a network scanning utility designed to perform various types of scans with optional proxy support and Telegram integration for reporting results.
Table of Contents

    Prerequisites
    Setup
    Usage
    Proxy Management
    Telegram Integration
    Troubleshooting

Prerequisites

Before using the j0kah Recon Tool, ensure the following:

    Go Language: Make sure Go is installed. You can download it from golang.org.
    Nmap: Install Nmap for network scanning. Download it from nmap.org.
    Telegram Bot Token: Create a Telegram bot and obtain the token from BotFather.

Setup

    Clone the Repository:

    bash

git clone https://github.com/yourusername/j0kah
cd j0kah

Build the Project:

Build the Go application:

bash

go build -o j0kah

Create Configuration File:

Create a config.ini file with the following content:

ini

    token=put token here
    chat_id=put chat id here

    Replace put token here and put chat id here with your actual Telegram bot token and chat ID.

Usage

    Run j0kah Recon Tool:

    Execute the tool:

    bash

    ./j0kah

    Follow the prompts to:
        Enter the IP/domain to scan.
        Choose the scan type (e.g., SYN-ACK Scan, UDP Scan).
        Decide whether to fetch and use proxies.
        Enter the scan duration and other options.
        Choose how to handle results (save to file, send via Telegram, or both).

Proxy Management

To use proxies for scanning:

    Proxies are Integrated: The tool automatically fetches proxies from Proxy List API and saves them to proxy.list if you choose to use them.

    Proxy Configuration: If proxies are available, they will be used in the scan. The proxy list is fetched and saved automatically, and the proxy configuration is handled within the tool.

Telegram Integration

To receive scan results via Telegram:

    Create a Telegram Bot:
        Use BotFather to create a new bot and obtain the token.

    Configure the Bot:
        Save your bot token and chat ID in the config.ini file as described in the Setup section.

    Send Results:
        During tool usage, choose to send results to Telegram.

Troubleshooting

    Error: unexpected status code: 404: This indicates an issue with fetching proxies. Ensure the proxy URL is correct and reachable.
    Failed to create file: Check file permissions in the directory where you're trying to save files.
    Telegram errors: Verify the config.ini file for the correct token and chat ID. Ensure the bot has permission to send messages to the chat ID.
