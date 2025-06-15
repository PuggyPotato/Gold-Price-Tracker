# 🪙 Gold & Silver Price Tracker Telegram Bot (Golang)

A Telegram bot written in Go that fetches **real-time gold and silver prices**, allows users to set **automatic update intervals**, and receive **price alerts** when a target is met. Lightweight, fast, and perfect for anyone interested in precious metal market tracking.

## 📌 Features

- `/start` – Introduces the bot and explains available features.
- `gold`, `silver`, or `price` – Instantly returns current prices.
- `/interval` – Sets an update interval (e.g. `10min`, `2hour`, `1day`) to get price updates periodically.
- `/settarget` – Lets users set a price threshold (e.g. `gold exceed 2500`, `silver below 22`) to trigger notifications.
- `/stop` – Stops all interval updates and price tracking for the user.

## 🚀 Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/PuggyPotato/Gold-Price-Tracker.git
cd Gold-Price-Tracker
```

### 2. Environment Setup
Create a `.env` file in the root directory:
```
API=your_telegram_bot_token
```
Get your bot token from [@BotFather](https://t.me/botfather) on Telegram.

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Run the Bot
```bash
go run main.go
```

## 🛠 Technologies Used

- **Go (Golang)** – Main programming language
- **Telegram Bot API** – For sending/receiving messages
- **Real-time data** – Fetched from `https://data-asg.goldprice.org/dbXRates/USD`
- **dotenv** – For secure environment variable management

## 📂 Commands Overview

| Command       | Description                                             |
|---------------|---------------------------------------------------------|
| `/start`      | Intro message and help                                  |
| `gold`        | Shows current gold price in USD                         |
| `silver`      | Shows current silver price in USD                       |
| `price`       | Shows both gold and silver prices                       |
| `/interval`   | Sets periodic update interval (e.g. `10min`, `2hour`)   |
| `/settarget`  | Notifies when price exceeds or drops below a target     |
| `/stop`       | Clears intervals and alerts for the user                |

