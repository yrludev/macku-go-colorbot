
# macku-go-colorbot

This is a Go port of the ValorantMakcuColorbot, using Makcu-Go for mouse control. It detects a specific color on your screen and moves/clicks the mouse using a Makcu hardware device.

## What it does

- Captures a region of your screen and detects a target color (using HSV range from `settings.ini`).
- Moves the mouse towards the detected color (aimbot).
- Clicks the mouse when the color is within a trigger region (triggerbot).
- All mouse actions are sent through a Makcu device (hardware mouse emulator).
- All settings (color, FOV, speed, etc.) are configurable in `settings.ini`.
- Simple command-line menu: start/stop the bot, see connection status, and exit.

## Setup Instructions

1. **Install Go** (https://go.dev/dl/) and make sure it's in your PATH.
2. **Install dependencies:**
	- Open a terminal in this folder and run:
	  - `go get github.com/kbinani/screenshot`
	  - `go get github.com/nullpkt/Makcu-Go`
	  - `go get gopkg.in/ini.v1`
3. **Connect your Makcu device** to your PC.
4. **Edit `settings.ini`** to match your color and aimbot preferences.
5. **Build and run:**
	- `go build -o colorbot.exe`
	- Double-click `colorbot.exe` or run `./colorbot.exe` in a terminal.

## Usage

1. Run `colorbot.exe`.
2. Use the menu:
	- `1` to start the bot (will print connection status to Makcu)
	- `2` to stop the bot
	- `3` to exit
3. The bot will print `[INFO] Connected to Makcu on ...` if the device is found and working.

## Features

- Color-based aimbot and triggerbot
- Uses Makcu hardware for all mouse actions
- All settings in `settings.ini`
- Simple CLI menu for control

## Credits
- Go port: yrludev

## DISCLAIMER

- Compatibility: This software is designed for Makcu boards only. Only to be used with 2 PC setup.
- Responsibility: This software is intended for educational purposes only. I am not responsible for any account bans, penalties, or any other consequences that may result from using this tool. Use it at your own risk and be aware of the potential implications.
