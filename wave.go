package gakki

import (
	"math"
)

const (
	SampleRate = 44100
)

type Wave struct {
	freq      float64
	f         func(x, lambda float64) float64
	pos       int64
	remaining []byte
}

func NewWave(freq float64, f func(x, lambda float64) float64) *Wave {
	return &Wave{
		freq: freq,
		f:    f,
	}
}

func (w *Wave) Read(buf []byte) (int, error) {
	if len(w.remaining) > 0 {
		n := copy(buf, w.remaining)
		w.remaining = w.remaining[n:]
		return n, nil
	}

	fullBufLen := len(buf)
	if mod := fullBufLen % 4; mod > 0 {
		fullBufLen += 4 - mod
	}
	fullBuf := make([]byte, fullBufLen)

	p := w.pos / 4
	for i := 0; i < fullBufLen/4; i++ {
		val := w.f(float64(p), float64(SampleRate)/w.freq)
		b := int16(val * 0.3 * math.MaxInt16)
		idx := i * 4
		fullBuf[idx] = byte(b)
		fullBuf[idx+1] = byte(b >> 8)
		fullBuf[idx+2] = byte(b)
		fullBuf[idx+3] = byte(b >> 8)
		p++
	}

	w.pos += int64(fullBufLen)

	n := copy(buf, fullBuf)
	w.remaining = fullBuf[n:]

	return n, nil
}
