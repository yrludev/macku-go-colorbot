package colorbot

import (
	"image"
	"macku-go-colorbot/capture"
	"macku-go-colorbot/config"
	"macku-go-colorbot/mouse"
	"math"
	"time"
)

type HSV struct {
	H, S, V uint8
}

type Colorbot struct {
	capturer           *capture.Capture
	mouse              *mouse.Mouse
	settings           *config.Settings
	lowerColor         HSV
	upperColor         HSV
	xSpeed, ySpeed     float64
	xFov, yFov         int
	targetOffset       float64
	xRange, yRange     int
	minDelay, maxDelay int
	toggleButton       int
	screenCenter       image.Point
}

func NewColorbot(settings *config.Settings, mouse *mouse.Mouse, capturer *capture.Capture) *Colorbot {
	lower := settings.GetIntList("ColorRange", "lower_color")
	upper := settings.GetIntList("ColorRange", "upper_color")
	return &Colorbot{
		capturer:     capturer,
		mouse:        mouse,
		settings:     settings,
		lowerColor:   HSV{uint8(lower[0]), uint8(lower[1]), uint8(lower[2])},
		upperColor:   HSV{uint8(upper[0]), uint8(upper[1]), uint8(upper[2])},
		xSpeed:       settings.GetFloat("Aimbot", "xSpeed"),
		ySpeed:       settings.GetFloat("Aimbot", "ySpeed"),
		xFov:         settings.GetInt("Aimbot", "xFov"),
		yFov:         settings.GetInt("Aimbot", "yFov"),
		targetOffset: settings.GetFloat("Aimbot", "targetOffset"),
		xRange:       settings.GetInt("Triggerbot", "xRange"),
		yRange:       settings.GetInt("Triggerbot", "yRange"),
		minDelay:     settings.GetInt("Triggerbot", "minDelay"),
		maxDelay:     settings.GetInt("Triggerbot", "maxDelay"),
		toggleButton: settings.GetInt("ToggleButton", "button"),
		screenCenter: image.Point{settings.GetInt("Aimbot", "xFov") / 2, settings.GetInt("Aimbot", "yFov") / 2},
	}
}

func (c *Colorbot) Process(action string) {
	img, err := c.capturer.GetScreen()
	if err != nil {
		return
	}
	target := findColorTarget(img, c.lowerColor, c.upperColor, c.screenCenter)
	if target != nil {
		cX, cY := target.X, target.Y-int(c.targetOffset)
		xDiff := cX - c.screenCenter.X
		yDiff := cY - c.screenCenter.Y
		if action == "move" {
			c.mouse.Move(int(float64(xDiff)*c.xSpeed), int(float64(yDiff)*c.ySpeed))
		} else if action == "click" {
			if abs(xDiff) <= c.xRange && abs(yDiff) <= c.yRange {
				time.Sleep(time.Duration(randInt(c.minDelay, c.maxDelay)) * time.Millisecond)
				c.mouse.Click()
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func randInt(min, max int) int {
	return min + int(time.Now().UnixNano())%(max-min+1)
}

func findColorTarget(img *image.RGBA, lower, upper HSV, center image.Point) *image.Point {
	var minDist float64 = math.MaxFloat64
	var closest *image.Point
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			h, s, v := rgbToHSV(r8, g8, b8)
			if inHSVRange(h, s, v, lower, upper) {
				dist := math.Hypot(float64(x-center.X), float64(y-center.Y))
				if dist < minDist {
					minDist = dist
					closest = &image.Point{X: x, Y: y}
				}
			}
		}
	}
	return closest
}

func inHSVRange(h, s, v uint8, lower, upper HSV) bool {
	return h >= lower.H && h <= upper.H &&
		s >= lower.S && s <= upper.S &&
		v >= lower.V && v <= upper.V
}

func rgbToHSV(r, g, b uint8) (h, s, v uint8) {
	fr := float64(r) / 255.0
	fg := float64(g) / 255.0
	fb := float64(b) / 255.0
	max := math.Max(fr, math.Max(fg, fb))
	min := math.Min(fr, math.Min(fg, fb))
	delta := max - min

	var hh float64
	if delta == 0 {
		hh = 0
	} else if max == fr {
		hh = 60 * math.Mod(((fg-fb)/delta), 6)
	} else if max == fg {
		hh = 60 * (((fb - fr) / delta) + 2)
	} else {
		hh = 60 * (((fr - fg) / delta) + 4)
	}
	if hh < 0 {
		hh += 360
	}

	var ss float64
	if max == 0 {
		ss = 0
	} else {
		ss = delta / max
	}

	return uint8(hh / 2), uint8(ss * 255), uint8(max * 255)
}
