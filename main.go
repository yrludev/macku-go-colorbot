package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"macku-go-colorbot/capture"
	"macku-go-colorbot/colorbot"
	"macku-go-colorbot/config"
	"macku-go-colorbot/mouse"

	makcu "github.com/nullpkt/Makcu-Go"
)

var running int32

func main() {
	reader := bufio.NewReader(os.Stdin)
	quit := make(chan struct{})

	for {
		fmt.Println("\nMacku Go Colorbot")
		fmt.Println("1. Start Bot")
		fmt.Println("2. Stop Bot")
		fmt.Println("3. Exit")
		fmt.Print("Select an option: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			if atomic.CompareAndSwapInt32(&running, 0, 1) {
				go runBot(quit)
			} else {
				fmt.Println("[WARN] Bot is already running.")
			}
		case "2":
			if atomic.CompareAndSwapInt32(&running, 1, 0) {
				fmt.Println("[INFO] Bot stopped.")
				close(quit)
				quit = make(chan struct{}) // reset for next run
			} else {
				fmt.Println("[WARN] Bot is not running.")
			}
		case "3":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func runBot(quit chan struct{}) {
	settings, err := config.NewSettings("settings.ini")
	if err != nil {
		fmt.Println("[ERROR] Failed to load settings.ini")
		atomic.StoreInt32(&running, 0)
		return
	}

	makcu.Debug = true 
	port, err := makcu.Find()
	if err != nil {
		fmt.Println("[ERROR] Makcu not found! Please connect your device.")
		atomic.StoreInt32(&running, 0)
		return
	}
	conn, err := makcu.Connect(port, 115200)
	if err != nil {
		fmt.Println("[ERROR] Makcu connect error!")
		atomic.StoreInt32(&running, 0)
		return
	}
	defer conn.Close()

	fmt.Printf("[INFO] Connected to Makcu on %s\n", port)

	m := mouse.NewMouse(conn)
	centerX := settings.GetInt("Aimbot", "xFov") / 2
	centerY := settings.GetInt("Aimbot", "yFov") / 2
	capt := capture.NewCapture(
		centerX-settings.GetInt("Aimbot", "xFov")/2,
		centerY-settings.GetInt("Aimbot", "yFov")/2,
		settings.GetInt("Aimbot", "xFov"),
		settings.GetInt("Aimbot", "yFov"),
	)
	cbot := colorbot.NewColorbot(settings, m, capt)

	fmt.Println("[INFO] Bot is running. Press '2' in the menu to stop.")
	for atomic.LoadInt32(&running) == 1 {
		select {
		case <-quit:
			return
		default:
			cbot.Process("move")
			cbot.Process("click")
			time.Sleep(10 * time.Millisecond)
		}
	}
	fmt.Println("[INFO] Bot stopped.")
}
