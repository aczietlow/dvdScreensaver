package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
	"time"
)

const (
	dvd = "\xF0\x9F\x93\x80"
)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func emitString(s tcell.Screen, text string, x int, y int, style tcell.Style) {
	for _, t := range text {
		s.SetContent(x, y, t, nil, style)
		x++
	}
}

func edgeDetected(x int, y int, maxWidth int, maxHeight int) bool {
	switch {
	case x == 0 && y == 0:
	case x == 0 && y == maxHeight:
	case x == maxWidth && y == 0:
	case x == maxWidth && y == maxHeight:
		return true
	}
	return false
}

func bounce(s tcell.Screen, x int, y int) {
	// @TODO play with this.
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	emitString(s, fmt.Sprintf("%v", dvd), x, y, defStyle)
}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

	w, h := s.Size()
	//emitString(s, fmt.Sprintf("%v, %v", w, h), 0, 0, defStyle)

	x, y := 0, 1
	bounce(s, x, y)
	xForward := true
	yForward := true

	// Event loop
	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	for {
		s.Show()

		time.Sleep(500 * time.Millisecond)
		s.Clear()
		emitString(s, fmt.Sprintf("%v, %v", w, h), 10, 0, defStyle)
		emitString(s, fmt.Sprintf("%v, %v", x, y), 0, 0, defStyle)

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			}
		}

		bounce(s, x, y)
		if edgeDetected(x, y, w, h) {
			break
		}
		// insert shitty math and logic bc I'm hungry
		if x == w-1 {
			xForward = false
		}
		if x == 0 {
			xForward = true
		}
		if y == h-1 {
			yForward = false
		}
		if y == 0 {
			yForward = true
		}

		if xForward {
			x++
		} else {
			x--
		}
		if yForward {
			y++
		} else {
			y--
		}

		// Post an event to the queue to continue the flow.
		s.PostEvent(tcell.NewEventInterrupt(nil))
	}

}
