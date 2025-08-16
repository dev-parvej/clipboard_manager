# Clipboard Manager

A simple and efficient clipboard manager built with Go and Fyne that helps you keep track of your clipboard history.

## Features

- 📋 Automatically captures text copied to clipboard
- 🖼️ Supports image clipboard content (Not yet working)
- 🕒 Maintains history with timestamps
- 🗑️ Auto-cleanup of entries older than 30 days
- 💻 System tray integration
- 🔍 View and manage clipboard history
- 🚀 Lightweight and efficient

## Screenshot

![Clipboard Manager Screenshot](Screenshot%202025-08-16%20135230.png)

## Installation

1. Make sure you have Go 1.24.5 or later installed
2. Clone the repository:
```sh
git clone https://github.com/dev-parvej/clipboard_manager.git
```
3. Install dependencies:
```sh
go mod download
```
4. Build and run:
```sh
go run .
```

## Usage

1. The app runs in the system tray
2. Click the tray icon to:
   - View clipboard history
   - Clear all entries
   - Exit the application
3. Copy text or images as usual - they'll be automatically captured
4. Click "View History" to see and manage your clipboard entries

## Development

Built with:
- [Go](https://golang.org/)
- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- SQLite for storage

## Storage

The application stores its data in:
- `~/.clipboard_manager/clipboard.db` - SQLite database
- `~/.clipboard_manager/images` - Captured images

## License

MIT License
