package gakki

import (
	"errors"
	"io"

	"github.com/hajimehoshi/oto"
)

type UI interface {
	SetKeyboard([][]rune)
	SetKeyDisplayFunc(func(rune) string)

	Run(keydownCh, keyupCh chan<- rune) error
	Close()
}

type Gakki struct {
	UI       UI
	Keyboard [][]rune
	WaveFunc func(x, wavelength float64) float64

	FreqMap  map[rune]float64
	ChordMap map[rune]Chord

	OnKeyPress     func(key rune)
	KeyDisplayFunc func(key rune) string

	otoCtx *oto.Context
}

func NewGakki() (*Gakki, error) {
	ctx, err := oto.NewContext(SampleRate, 2, 2, 4096)
	if err != nil {
		return nil, err
	}

	gakki := Gakki{
		FreqMap:  map[rune]float64{},
		ChordMap: map[rune]Chord{},
		otoCtx:   ctx,
	}
	return &gakki, nil
}

func (g *Gakki) Close() {
	g.UI.Close()
	g.otoCtx.Close()
}

func (g *Gakki) Run() error {
	if g.UI == nil {
		return errors.New("Gakki: UI is not set")
	}
	if g.Keyboard == nil {
		return errors.New("Gakki: Keyboard is not set")
	}
	if g.WaveFunc == nil {
		return errors.New("Gakki: WaveFunc is not set")
	}

	g.UI.SetKeyboard(g.Keyboard)
	g.UI.SetKeyDisplayFunc(g.KeyDisplayFunc)

	keydownCh := make(chan rune, 10)
	keyupCh := make(chan rune, 10)

	go g.RunOto(keydownCh, keyupCh)

	return g.UI.Run(keydownCh, keyupCh)
}

func (g *Gakki) RunOto(keydownCh, keyupCh <-chan rune) {
	stopChannels := map[rune]chan struct{}{}
	for {
		select {
		case key := <-keydownCh:
			if stopChannels[key] != nil {
				break
			}
			if g.OnKeyPress != nil {
				g.OnKeyPress(key)
			}
			stopCh := make(chan struct{}, 1)
			stopChannels[key] = stopCh
			go g.play(key, stopCh)

		case key := <-keyupCh:
			stopCh, ok := stopChannels[key]
			if !ok {
				break
			}
			close(stopCh)
			delete(stopChannels, key)
		}
	}
}

func (g *Gakki) play(key rune, stopCh <-chan struct{}) {
	freq, ok := g.FreqMap[key]
	if !ok {
		chord, ok := g.ChordMap[key]
		if !ok {
			return
		}
		freq = chord.Freq()
	}

	p := g.otoCtx.NewPlayer()
	defer p.Close()

	w := NewWave(freq, g.WaveFunc)
LOOP:
	for {
		select {
		case _, ok := <-stopCh:
			if !ok {
				break LOOP
			}
		default:
			if _, err := io.CopyN(p, w, 32*1024); err != nil {
				panic(err)
			}
		}
	}
}
