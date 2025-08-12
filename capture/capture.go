package capture

import (
	"image"

	"github.com/kbinani/screenshot"
)

type Capture struct {
	x, y, width, height int
}

func NewCapture(x, y, width, height int) *Capture {
	return &Capture{x: x, y: y, width: width, height: height}
}

func (c *Capture) GetScreen() (*image.RGBA, error) {
	return screenshot.Capture(c.x, c.y, c.width, c.height)
}
