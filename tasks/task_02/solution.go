package main

func rotateRunes(s string, shift int) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	res := make([]rune, len(runes))

	shifted := shift % len(runes)
	if shifted < 0 {
		shifted += len(runes)
	}

	for i := range runes {
		destination := (i - shifted)
		if destination < 0 {
			destination += len(runes)
		}
		res[destination] = runes[i]
	}

	return string(res)
}
