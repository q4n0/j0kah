j0kah

j0kah is a powerful network scanning tool designed for both performance and style. With features like customizable scan types, adjustable concurrency, and integration with Telegram for notifications, j0kah brings a touch of chaos to network discovery.
Features

    Customizable Scans: Choose from various scan types including SYN, UDP, TCP, ACK, and more.
    Adjustable Concurrency: Set the number of concurrent threads for your scan.
    Configurable Scan Duration: Define how long the scan should run.
    Local File Saving: Optionally save scan results to a local file.
    Telegram Integration: Send scan results to a Telegram chat for remote notification.
    Styled Output: Enjoy a touch of humor and style with colorful and provocative output.

Installation

To use j0kah, follow these steps:

    Clone the Repository:

    bash

git clone https://github.com/q4n0/j0kah.git
cd j0kah

Install Dependencies:

Ensure you have Go installed on your system. Install the required Go packages:

bash

go mod tidy

You also need nmap installed on your system, as j0kah relies on it for performing network scans.

Create and Configure the Telegram Bot:

    Create a Telegram bot using BotFather and get the bot token.

    Get the chat ID for the chat where you want to receive scan results.

    Create a config.ini file in the project directory with the following content:

    makefile

        token=YOUR_TELEGRAM_BOT_TOKEN
        chat_id=YOUR_TELEGRAM_CHAT_ID

Usage

    Run the Tool:

    bash

    go run j0kah.go

    Follow the Prompts:
        Enter the comma-separated list of targets.
        Choose the scan type.
        Set the concurrency level.
        Define the scan duration.
        Decide if you want to save results locally.
        Choose if you want to send results to Telegram.
        Provide any additional arguments for the scan.

    View Results:
        If you chose to save results locally, they will be saved in scan_results.txt.
        If you opted for Telegram notifications, you’ll receive the scan results in your specified chat.

Example Output

vbnet

=========================================================================
       Welcome to the Insanity Show, Where Chaos Reigns Supreme,
       And you're the Star of the Act
           b0urn3 - The Mad Genius
=========================================================================
Ready to dance with madness? Let the games begin.

Preparing to perform a SYN scan on target1,target2 with args '-T4 -A' and concurrency 10.
Progress: 10% Complete. If you’re still here, you must enjoy torment. Either that or you’ve mastered the art of waiting. 
...
Wow, you survived! Now go ahead and reward yourself for this grueling achievement. Or just admit you have nothing better to do.

Scan complete! Results will be sent to Telegram or saved locally based on your choice.

Results have been delivered to Telegram. Brace yourself—because you just invited chaos into your chat.

=================================================================
Created by b0urn3
GitHub: https://github.com/q4n0
Instagram: https://www.instagram.com/onlybyhive
Email: b0urn3@proton.me
=================================================================
