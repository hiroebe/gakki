package main

import (
	"flag"
	"math"

	"github.com/hiroebe/gakki"
	"github.com/hiroebe/gakki/cui"
	"github.com/hiroebe/gakki/gui"
)

var (
	flagGUI = flag.Bool("gui", false, "GUI mode")
)

func sin(x, wavelength float64) float64 {
	r := 1 - x/gakki.SampleRate
	if r < 0 {
		r = 0
	}
	return math.Sin(2 * math.Pi / wavelength * x)
}

func sin2(x, wavelength float64) float64 {
	v1 := math.Sin(2 * math.Pi / wavelength * x)
	v2 := math.Sin(2 * math.Pi / wavelength * 2 * x)
	v3 := math.Sin(2 * math.Pi / wavelength * 1.5 * x)
	v4 := math.Sin(2 * math.Pi / wavelength * 2.5 * x)
	v := (v1 + v2 + v3 + v4) / 4
	return v
}

func main() {
	flag.Parse()

	g, err := gakki.NewGakki()
	if err != nil {
		panic(err)
	}
	defer g.Close()

	if *flagGUI {
		g.UI = gui.GetDefaultGUI()
		g.Keyboard = gakki.GetDefaultUSKeyboard()
	} else {
		g.UI = cui.GetDefaultCUI()
		g.Keyboard = gakki.GetDefaultJISKeyboard()
	}

	g.WaveFunc = sin

	startChord := gakki.NewChord("C4")
	g.ChordMap = gakki.GetDefaultChordMap(startChord, g.Keyboard[1], g.Keyboard[2])

	g.OnKeyPress = func(key rune) {
		switch key {
		case '=', '^':
			startChord = startChord.Up()
			g.ChordMap = gakki.GetDefaultChordMap(startChord, g.Keyboard[1], g.Keyboard[2])
		case '-':
			startChord = startChord.Down()
			g.ChordMap = gakki.GetDefaultChordMap(startChord, g.Keyboard[1], g.Keyboard[2])
		case 'z':
			g.WaveFunc = sin
		case 'x':
			g.WaveFunc = sin2
		}
	}
	g.KeyDisplayFunc = func(key rune) string {
		switch key {
		case '=', '^':
			return ">"
		case '-':
			return "<"
		case 'z':
			return "1"
		case 'x':
			return "2"
		}

		chord, ok := g.ChordMap[key]
		if !ok {
			return ""
		}
		return chord.String()
	}

	if err = g.Run(); err != nil {
		panic(err)
	}
}
