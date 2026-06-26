# 🚀 FeedBot

FeedBot is a personal Telegram automation bot written in Go.

The goal of FeedBot is to collect updates from different sources and deliver them directly to Telegram in a clean and simple format.

## Features

### 🎥 YouTube

* Add YouTube channels by link or handle
* Automatically track new uploads
* Receive notifications when new videos are published
* List saved channels
* Remove channels


### 🌦 Weather

* Search and select a weather location
* Save the selected location permanently
* Fallback location when no location is configured
* Get the weather report for today
* Get the weather report for tomorrow
* Receive an automatic weather report every morning for the current day
* Receive an automatic weather report every evening for the next day
* Shows current temperature, feels-like temperature, min/max temperature, sunrise, sunset and hourly forecast

### 🤖 Automation System

* Background scheduler
* Queue-based notification delivery
* Configurable check intervals
* Scheduled daily automation jobs

### 🔒 Private Bot

* Access restricted to allowed Telegram user IDs
* Environment-based configuration


## Requirements

* Go 1.26+
* Docker (optional)
* Telegram Bot Token

## Configuration

Create a `.env` file in the project root:

```env
TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN
ALLOWED_USER_IDS=123456789
```

## Run Locally

Install dependencies:

```bash
go mod download
```

Start FeedBot:

```bash
go run .
```

## Docker

Build the image:

```bash
docker build -t feedbot .
```

Start with Docker Compose:

```bash
docker compose up -d
```

Stop the container:

```bash
docker compose down
```

View logs:

```bash
docker logs -f feedbot
```

## License

MIT License

---

Built with Go and Telegram Bot API.
