package main

import (
	"macku-go-colorbot/capture"
	"macku-go-colorbot/colorbot"
	"macku-go-colorbot/config"
	"macku-go-colorbot/mouse"
	"sync/atomic"
	"time"

	makcu "github.com/nullpkt/Makcu-Go"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var running int32

func main() {
	myApp := app.New()
	win := myApp.NewWindow("Macku Go Colorbot")

	status := widget.NewLabel("Status: Idle")
	startBtn := widget.NewButton("Start", nil)
	stopBtn := widget.NewButton("Stop", nil)
	stopBtn.Disable()

	quit := make(chan struct{})

	startBtn.OnTapped = func() {
		if atomic.CompareAndSwapInt32(&running, 0, 1) {
			status.SetText("Status: Running")
			startBtn.Disable()
			stopBtn.Enable()
			go runBot(status, quit)
		}
	}

	stopBtn.OnTapped = func() {
		if atomic.CompareAndSwapInt32(&running, 1, 0) {
			status.SetText("Status: Stopped")
			stopBtn.Disable()
			startBtn.Enable()
			close(quit)
		}
	}

	win.SetContent(container.NewVBox(
		widget.NewLabel("Macku Go Colorbot"),
		status,
		container.NewHBox(startBtn, stopBtn),
	))

	win.Resize(fyne.NewSize(320, 120))
	win.ShowAndRun()
}

func runBot(status *widget.Label, quit chan struct{}) {
	settings, err := config.NewSettings("settings.ini")
	if err != nil {
		status.SetText("Failed to load settings.ini")
		return
	}

	port, err := makcu.Find()
	if err != nil {
		status.SetText("Makcu not found!")
		return
	}
	conn, err := makcu.Connect(port, 115200)
	if err != nil {
		status.SetText("Makcu connect error!")
		return
	}
	defer conn.Close()

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

	status.SetText("Status: Running")
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
	status.SetText("Status: Idle")
}
