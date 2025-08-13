package main

import (
	"fmt"
	"image"
	"syscall"
	"time"
	"unsafe"
)

// --- Config ---
const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
	CenterX      = ScreenWidth / 2
	CenterY      = ScreenHeight / 2
	ROIHalfSize  = 100 // Increased for easier detection
)

var (
	// HSV for new color range
	HSVLower = [3]uint8{28, 38, 200}
	HSVUpper = [3]uint8{53, 255, 255}
)

func main() {
	fmt.Println("[INFO] Starting CLI aimbot using Windows syscalls. Hold right mouse button to drag to target. Press Ctrl+C to exit.")

	user32 := syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState := user32.NewProc("GetAsyncKeyState")
	const VK_RBUTTON = 0x02

	for {
		// Check if right mouse button is held
		state, _, _ := getAsyncKeyState.Call(VK_RBUTTON)
		if state&0x8000 != 0 {
			img, err := captureScreen()
			if err != nil {
				fmt.Println("[ERROR] Screen capture failed:", err)
				continue
			}

			// Scan center region for color
			centerRect := image.Rect(CenterX-ROIHalfSize, CenterY-ROIHalfSize, CenterX+ROIHalfSize, CenterY+ROIHalfSize)
			found := false
			var foundX, foundY int
			for y := centerRect.Min.Y; y < centerRect.Max.Y; y++ {
				for x := centerRect.Min.X; x < centerRect.Max.X; x++ {
					r, g, b, _ := img.At(x, y).RGBA()
					h, s, v := rgbToHSV(uint8(r>>8), uint8(g>>8), uint8(b>>8))
					if inHSVRange(h, s, v, HSVLower, HSVUpper) {
						found = true
						foundX, foundY = x, y
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				dx := foundX - CenterX
				dy := foundY - CenterY
				// Clamp dx, dy to -100..100
				if dx > 100 {
					dx = 100
				}
				if dx < -100 {
					dx = -100
				}
				if dy > 100 {
					dy = 100
				}
				if dy < -100 {
					dy = -100
				}
				fmt.Printf("[AIM] Found at (%d,%d), center (%d,%d), dx=%d dy=%d\n", foundX, foundY, CenterX, CenterY, dx, dy)
				moveMouseRelative(dx, dy)
			} else {
				fmt.Println("[AIM] No target color found in center region.")
			}
		}
		// Small sleep to avoid 100% CPU
		time.Sleep(10 * time.Millisecond)
	}
}

// --- Windows screen capture using syscalls ---
func captureScreen() (image.Image, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	gdi32 := syscall.NewLazyDLL("gdi32.dll")

	getDC := user32.NewProc("GetDC")
	releaseDC := user32.NewProc("ReleaseDC")
	createCompatibleDC := gdi32.NewProc("CreateCompatibleDC")
	createCompatibleBitmap := gdi32.NewProc("CreateCompatibleBitmap")
	selectObject := gdi32.NewProc("SelectObject")
	bitBlt := gdi32.NewProc("BitBlt")
	getDIBits := gdi32.NewProc("GetDIBits")
	deleteObject := gdi32.NewProc("DeleteObject")
	deleteDC := gdi32.NewProc("DeleteDC")

	hdc, _, _ := getDC.Call(0)
	hdcMem, _, _ := createCompatibleDC.Call(hdc)
	hBitmap, _, _ := createCompatibleBitmap.Call(hdc, uintptr(ScreenWidth), uintptr(ScreenHeight))
	selectObject.Call(hdcMem, hBitmap)
	bitBlt.Call(hdcMem, 0, 0, uintptr(ScreenWidth), uintptr(ScreenHeight), hdc, 0, 0, 0x00CC0020)

	type BITMAPINFOHEADER struct {
		Size          uint32
		Width         int32
		Height        int32
		Planes        uint16
		BitCount      uint16
		Compression   uint32
		SizeImage     uint32
		XPelsPerMeter int32
		YPelsPerMeter int32
		ClrUsed       uint32
		ClrImportant  uint32
	}
	type BITMAPINFO struct {
		Header BITMAPINFOHEADER
		Colors [1]uint32
	}

	var bmi BITMAPINFO
	bmi.Header.Size = 40
	bmi.Header.Width = int32(ScreenWidth)
	bmi.Header.Height = -int32(ScreenHeight) // top-down
	bmi.Header.Planes = 1
	bmi.Header.BitCount = 32
	bmi.Header.Compression = 0

	buf := make([]byte, ScreenWidth*ScreenHeight*4)
	getDIBits.Call(hdcMem, hBitmap, 0, uintptr(ScreenHeight), uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&bmi)), 0)

	img := image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))
	copy(img.Pix, buf)

	deleteObject.Call(hBitmap)
	deleteDC.Call(hdcMem)
	releaseDC.Call(0, hdc)

	return img, nil
}

// --- HSV helpers ---
func rgbToHSV(r, g, b uint8) (uint8, uint8, uint8) {
	max := float64(max3(r, g, b))
	min := float64(min3(r, g, b))
	v := max
	s := 0.0
	h := 0.0
	if max != 0 {
		s = (max - min) / max
	}
	if max == min {
		h = 0
	} else if max == float64(r) {
		h = (float64(g) - float64(b)) / (max - min)
	} else if max == float64(g) {
		h = 2 + (float64(b)-float64(r))/(max-min)
	} else {
		h = 4 + (float64(r)-float64(g))/(max-min)
	}
	h = h * 60
	if h < 0 {
		h += 360
	}
	return uint8(h / 2), uint8(s * 255), uint8(v)
}

func inHSVRange(h, s, v uint8, lower, upper [3]uint8) bool {
	return h >= lower[0] && h <= upper[0] &&
		s >= lower[1] && s <= upper[1] &&
		v >= lower[2] && v <= upper[2]
}

func max3(a, b, c uint8) uint8 {
	if a > b && a > c {
		return a
	}
	if b > c {
		return b
	}
	return c
}

func min3(a, b, c uint8) uint8 {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

// --- Windows mouse movement using syscalls ---
func moveMouseRelative(dx, dy int) {
	user32 := syscall.NewLazyDLL("user32.dll")
	getCursorPos := user32.NewProc("GetCursorPos")
	setCursorPos := user32.NewProc("SetCursorPos")

	type POINT struct {
		X int32
		Y int32
	}

	var pt POINT
	getCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	newX := int(pt.X) + dx
	newY := int(pt.Y) + dy
	setCursorPos.Call(uintptr(newX), uintptr(newY))
}
