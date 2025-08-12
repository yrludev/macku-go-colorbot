package mouse

import (
	"log"

	makcu "github.com/nullpkt/Makcu-Go"
)

type Mouse struct {
	conn *makcu.MakcuHandle
}

func NewMouse(conn *makcu.MakcuHandle) *Mouse {
	return &Mouse{conn: conn}
}

func (m *Mouse) Move(dx, dy int) {
	if m.conn != nil {
		m.conn.MoveMouse(dx, dy)
	}
}

func (m *Mouse) Click() {
	if m.conn != nil {
		err := m.conn.LeftClick()
		if err != nil {
			log.Println("Mouse click error:", err)
		}
	}
}
