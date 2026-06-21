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

### 🤖 Automation System

* Background scheduler
* Queue-based notification delivery
* Configurable check intervals

### 🔒 Private Bot

* Access restricted to allowed Telegram user IDs
* Environment-based configuration

## Planned Features

* 🌦 Weather reports
* 📰 RSS feeds
* 📰 Hacker News monitoring
* Additional automation modules

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
