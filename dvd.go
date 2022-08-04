package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	dvdChar      = "\xF0\x9F\x93\x80"
	confettiChar = "*"
	sleep        = 200
	duration     = 10
)

type confettiFlake struct {
	startingXCord int
	style         tcell.Style
	x             int
	y             int
}

func (c *confettiFlake) drawConfetti(s tcell.Screen, r *rand.Rand) {
	// Only 3 states the confetti could fall: left, right, neutral.
	f := r.Intn(3)
	switch {
	case f == 0:
		c.x = c.startingXCord - 1
	case f == 1:
		c.x = c.startingXCord + 1
	case f == 2:
		c.x = c.startingXCord
	}
	c.y = c.y + 1
	emitString(s, confettiChar, c.x, c.y, c.style)
}

func emitString(s tcell.Screen, text string, x int, y int, style tcell.Style) {
	for _, t := range text {
		s.SetContent(x, y, t, nil, style)
		x++
	}
}

func edgeDetected(x int, y int, maxWidth int, maxHeight int) bool {
	switch {
	case x == 0 && y == 0, x == 0 && y == maxHeight, x == maxWidth && y == 0, x == maxWidth && y == maxHeight:
		return true
		// Testing the end
		//case x == 4 && y == 5:
		//	return true
	}
	return false
}

func bounce(s tcell.Screen, x int, y int, style tcell.Style) {
	emitString(s, fmt.Sprintf("%v", dvdChar), x, y, style)
}

func dvd(s tcell.Screen, x, y, w, h int, style tcell.Style, quit func()) {
	bounce(s, x, y, style)

	xForward := true
	yForward := true

	// A literal goto, because why not. I'm not going to adhere to reasonable architecture at this point.
	// Technically, I'm labeling this loop, but still..
out:
	for {
		s.Clear()
		if edgeDetected(x, y, w, h) {
			break
		}
		// Cheap FPS
		time.Sleep(300 * time.Millisecond)

		// Debugging..
		//emitString(s, fmt.Sprintf("%v, %v", w, h), 10, 0, style)
		//emitString(s, fmt.Sprintf("%v, %v", x, y), 0, 0, style)

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
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				break out
			}
		}

		bounce(s, x, y, style)

		// Insert shitty math and logic bc I'm hungry
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
		err := s.PostEvent(tcell.NewEventInterrupt(nil))
		if err != nil {
			log.Fatalf("%+v", err)
		}
		s.Show()
	}

}

func confetti(s tcell.Screen, w, h, duration int, quit func()) {
	frames := duration * 10

	colorStlyes := []tcell.Color{
		tcell.ColorBlue,
		tcell.ColorPink,
		tcell.ColorHotPink,
		tcell.ColorYellow,
		tcell.ColorLimeGreen,
		tcell.ColorCornflowerBlue,
		tcell.ColorRed,
		tcell.ColorOrange,
		tcell.ColorOrchid,
		tcell.ColorAqua,
		tcell.ColorBlack,
		tcell.ColorFuchsia,
		tcell.ColorNavy,
		tcell.ColorWhite,
		tcell.ColorPlum,
		tcell.ColorViolet,
	}
	var allConfetti = make([]*confettiFlake, 0, w*h/3)
	for i := 0; i < frames; i++ {
		s.Clear()
		r := rand.New(rand.NewSource(int64(i)))

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

		// Create a random number of Confetti
		newConfetti := r.Intn(w / 10)

		for i2 := 0; i2 < newConfetti; i2++ {
			color := colorStlyes[r.Intn(len(colorStlyes))]
			style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(color)
			allConfetti = append(allConfetti, &confettiFlake{
				// Avoid having flakes fall out of bounds
				r.Intn(w-2) + 1,
				style,
				2,
				0,
			})
		}

		copyAllConfetti := allConfetti[:0]
		for _, c := range allConfetti {
			if !edgeDetected(0, c.y+1, w, h) {
				copyAllConfetti = append(copyAllConfetti, c)
			}
			c.drawConfetti(s, r)
		}
		allConfetti = copyAllConfetti

		err := s.PostEvent(tcell.NewEventInterrupt(nil))
		if err != nil {
			log.Fatalf("%+v", err)
		}
		time.Sleep(sleep * time.Millisecond)
		s.Show()
	}
}

func main() {
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	s.Clear()
	w, h := s.Size()
	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	dvd(s, 0, 1, w, h, defStyle, quit)
	s.Clear()
	confetti(s, w, h, duration, quit)
	s.Clear()
	quit()

}
