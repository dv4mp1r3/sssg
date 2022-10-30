package main

func ReplaceAtIndex(in string, r []rune, i int, l int) string {
	return in[:i] + string(r) + in[i+l:]
}
