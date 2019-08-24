package gakki

import (
	"fmt"
	"strconv"
)

var freqMap = map[string]float64{
	"A0":  27.500,
	"A#0": 29.135,
	"B0":  30.868,
	"C1":  32.703,
	"C#1": 34.648,
	"D1":  36.708,
	"D#1": 38.891,
	"E1":  41.203,
	"F1":  43.654,
	"F#1": 46.249,
	"G1":  48.999,
	"G#1": 51.913,
	"A1":  55.000,
	"A#1": 58.270,
	"B1":  61.735,
	"C2":  65.406,
	"C#2": 69.296,
	"D2":  73.416,
	"D#2": 77.782,
	"E2":  82.407,
	"F2":  87.307,
	"F#2": 92.499,
	"G2":  97.999,
	"G#2": 103.826,
	"A2":  110.000,
	"A#2": 116.541,
	"B2":  123.471,
	"C3":  130.813,
	"C#3": 138.591,
	"D3":  146.832,
	"D#3": 155.563,
	"E3":  164.814,
	"F3":  174.614,
	"F#3": 184.997,
	"G3":  195.998,
	"G#3": 207.652,
	"A3":  220.000,
	"A#3": 233.082,
	"B3":  246.942,
	"C4":  261.626,
	"C#4": 277.183,
	"D4":  293.665,
	"D#4": 311.127,
	"E4":  329.628,
	"F4":  349.228,
	"F#4": 369.994,
	"G4":  391.995,
	"G#4": 415.305,
	"A4":  440.000,
	"A#4": 466.164,
	"B4":  493.883,
	"C5":  523.251,
	"C#5": 554.365,
	"D5":  587.330,
	"D#5": 622.254,
	"E5":  659.255,
	"F5":  698.456,
	"F#5": 739.989,
	"G5":  783.991,
	"G#5": 830.609,
	"A5":  880.000,
	"A#5": 932.328,
	"B5":  987.767,
	"C6":  1046.502,
	"C#6": 1108.731,
	"D6":  1174.659,
	"D#6": 1244.508,
	"E6":  1318.510,
	"F6":  1396.913,
	"F#6": 1479.978,
	"G6":  1567.982,
	"G#6": 1661.219,
	"A6":  1760.000,
	"A#6": 1864.655,
	"B6":  1975.533,
	"C7":  2093.005,
	"C#7": 2217.461,
	"D7":  2349.318,
	"D#7": 2489.016,
	"E7":  2637.020,
	"F7":  2793.826,
	"F#7": 2959.955,
	"G7":  3135.963,
	"G#7": 3322.438,
	"A7":  3520.000,
	"A#7": 3729.310,
	"B7":  3951.066,
	"C8":  4186.009,
}

type Chord struct {
	name  byte
	num   int8
	sharp bool
}

func NewChord(s string) Chord {
	var num int
	sharp := len(s) == 3 && s[1] == '#'
	if sharp {
		num, _ = strconv.Atoi(s[2:3])
	} else {
		num, _ = strconv.Atoi(s[1:2])
	}
	return Chord{
		name:  s[0],
		num:   int8(num),
		sharp: sharp,
	}
}

func (c Chord) String() string {
	if c.sharp {
		return fmt.Sprintf("%s#%d", string(c.name), c.num)
	}
	return fmt.Sprintf("%s%d", string(c.name), c.num)
}

func (c Chord) Freq() float64 {
	return freqMap[c.String()]
}

func (c Chord) Up() Chord {
	c.name++
	if c.name == 'H' {
		c.name = 'A'
	}
	if c.name == 'C' {
		c.num++
	}
	return c
}

func (c Chord) Down() Chord {
	c.name--
	if c.name == 'A'-1 {
		c.name = 'G'
	}
	if c.name == 'B' {
		c.num--
	}
	return c
}

func (c Chord) HalfUp() (Chord, bool) {
	if c.name == 'B' || c.name == 'E' {
		return c, false
	}
	c.sharp = !c.sharp
	return c, true
}

func (c Chord) HalfDown() (Chord, bool) {
	c = c.Down()
	return c.HalfUp()
}
