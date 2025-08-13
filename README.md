
# macku-go-colorbot


This is a Go colorbot for a 2-PC Makcu setup. It captures a region of your screen, detects a specific color (HSV range), and sends mouse movement commands to a Makcu hardware device using the nullpkt/Makcu-Go library.

## What it does

- Captures a 200x200 region at the center of your screen and detects a target color (using a hardcoded HSV range).
- Moves the mouse towards the detected color (aimbot) by sending commands to a Makcu device (hardware mouse emulator) on another PC.
- All mouse actions are sent through Makcu using the nullpkt/Makcu-Go library.
- No settings.ini or CLI menu: all settings are in the code.

## Setup Instructions

1. **Install Go** (https://go.dev/dl/) and make sure it's in your PATH.
2. **Install dependency:**
	- Open a terminal in this folder and run:
	  - `go get github.com/nullpkt/Makcu-Go@latest`
	  - `go mod tidy`
3. **Connect your Makcu device** to your PC (the PC that will move the mouse).
4. **Build and run:**
	- `go build -o colorbot.exe main.go`
	- Double-click `colorbot.exe` or run `./colorbot.exe` in a terminal.

## Usage

1. Run `colorbot.exe`.
2. The bot will print `[INFO] Found Makcu device on ...` if the device is found and working.
3. Hold the right mouse button to activate aimbot movement.

## Features

- Color-based aimbot (no triggerbot)
- Uses Makcu hardware for all mouse actions
- HSV color range is hardcoded in main.go:
	- H: 28–53
	- S: 38–255
	- V: 200–255
- FOV is 200x200 pixels at the center of the screen (change ROIHalfSize in main.go to adjust)

## Credits
- Go port: yrludev
	(nullpkt/Makcu-Go integration)

## DISCLAIMER

- Compatibility: This software is designed for Makcu boards only. Only to be used with 2 PC setup.
- Responsibility: This software is intended for educational purposes only. I am not responsible for any account bans, penalties, or any other consequences that may result from using this tool. Use it at your own risk and be aware of the potential implications.
