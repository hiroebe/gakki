package gakki

func GetDefaultUSKeyboard() [][]rune {
	return [][]rune{
		[]rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '='},
		[]rune{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', '[', ']'},
		[]rune{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', ';', '\'', '\\'},
		[]rune{'z', 'x', 'c', 'v', 'b', 'n', 'm', ',', '.', '/'},
	}
}

func GetDefaultJISKeyboard() [][]rune {
	return [][]rune{
		[]rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '^'},
		[]rune{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', '@', '['},
		[]rune{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', ';', ':', ']'},
		[]rune{'z', 'x', 'c', 'v', 'b', 'n', 'm', ',', '.', '/'},
	}
}

func GetDefaultChordMap(start Chord, topKeys, midKeys []rune) map[rune]Chord {
	m := make(map[rune]Chord, len(topKeys)+len(midKeys))
	chord := start
	for i := 0; i < len(midKeys); i++ {
		mid := midKeys[i]
		m[mid] = chord
		if i == 0 {
			if hd, ok := chord.HalfDown(); ok {
				m[topKeys[i]] = hd
			}
		}
		if i < len(midKeys)-1 {
			if hu, ok := chord.HalfUp(); ok {
				m[topKeys[i+1]] = hu
			}
		}
		chord = chord.Up()
	}
	return m
}
