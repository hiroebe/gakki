package main

import (
	"math"

	"github.com/hiroebe/gakki"
	"github.com/hiroebe/gakki/gui"
)

func main() {
	g, err := gakki.NewGakki()
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.UI = gui.GetDefaultGUI()
	g.Keyboard = gakki.GetDefaultUSKeyboard()

	g.WaveFunc = func(x, wavelength float64) float64 {
		return math.Sin(2 * math.Pi / wavelength * x)
	}

	g.ChordMap = gakki.GetDefaultChordMap(gakki.NewChord("C4"), g.Keyboard[1], g.Keyboard[2])

	if err = g.Run(); err != nil {
		panic(err)
	}
}
