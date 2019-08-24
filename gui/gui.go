package gui

import (
	"github.com/andlabs/ui"
	"github.com/hiroebe/gakki"
)

const (
	keyWidth   = 30.0
	keyHeight  = 30.0
	keyPadding = 5.0
)

func GetDefaultGUI() gakki.UI {
	return &defaultGUI{}
}

type defaultGUI struct {
	keyboard       [][]rune
	keyDisplayFunc func(rune) string
}

func (g *defaultGUI) SetKeyboard(keyboard [][]rune) {
	g.keyboard = keyboard
}

func (g *defaultGUI) SetKeyDisplayFunc(f func(rune) string) {
	g.keyDisplayFunc = f
}

func (g *defaultGUI) Run(keydownCh, keyupCh chan<- rune) error {
	return ui.Main(func() {
		area := ui.NewArea(areaHandler{
			keydownCh:      keydownCh,
			keyupCh:        keyupCh,
			keyboard:       g.keyboard,
			keyDisplayFunc: g.keyDisplayFunc,
			typedKeys:      map[rune]bool{},
		})

		winWidth := (keyWidth + keyPadding*2) * (len(g.keyboard[0]) + 1)
		winHeight := (keyHeight + keyPadding*2) * 4
		window := ui.NewWindow("Gakki", int(winWidth), int(winHeight), false)
		window.SetChild(area)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
}

func (g *defaultGUI) Close() {}

type areaHandler struct {
	keydownCh      chan<- rune
	keyupCh        chan<- rune
	keyboard       [][]rune
	keyDisplayFunc func(rune) string

	typedKeys map[rune]bool
}

func (h areaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	brush := &ui.DrawBrush{
		Type: ui.DrawBrushTypeSolid,
		R:    0.9,
		G:    0.9,
		B:    0.9,
		A:    1,
	}
	brushTyped := &ui.DrawBrush{
		Type: ui.DrawBrushTypeSolid,
		R:    0.9,
		G:    0.4,
		B:    0.0,
		A:    1,
	}

	for i, keys := range h.keyboard {
		y := float64(i)*(keyHeight+keyPadding*2) + keyPadding

		for j, r := range keys {
			x := float64(j)*(keyWidth+keyPadding*2) + keyPadding + float64(i)*(keyWidth/2)

			path := ui.DrawNewPath(ui.DrawFillModeWinding)
			path.AddRectangle(x, y, keyWidth, keyHeight)
			path.End()
			if h.typedKeys[r] {
				p.Context.Fill(path, brushTyped)
			} else {
				p.Context.Fill(path, brush)
			}
			path.Free()

			var s string
			if h.keyDisplayFunc != nil {
				s = h.keyDisplayFunc(r)
			} else {
				s = string(r)
			}
			if s == "" {
				continue
			}
			attrStr := ui.NewAttributedString(s)
			tlParams := &ui.DrawTextLayoutParams{
				String:      attrStr,
				DefaultFont: &ui.FontDescriptor{Size: ui.TextSize(keyWidth / len(s))},
				Width:       keyWidth,
				Align:       ui.DrawTextAlignCenter,
			}
			tl := ui.DrawNewTextLayout(tlParams)
			p.Context.Text(tl, x+keyPadding/2, y)

			attrStr.Free()
			tl.Free()
		}
	}
}

func (areaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {}

func (areaHandler) MouseCrossed(a *ui.Area, left bool) {}

func (areaHandler) DragBroken(a *ui.Area) {}

func (h areaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	handled = false
	if ke.Modifier > 0 {
		return
	}

	handled = true
	if ke.Up {
		h.typedKeys[ke.Key] = false
		h.keyupCh <- ke.Key
	} else {
		if h.typedKeys[ke.Key] {
			return
		}
		h.typedKeys[ke.Key] = true
		h.keydownCh <- ke.Key
	}

	a.QueueRedrawAll()
	return
}
