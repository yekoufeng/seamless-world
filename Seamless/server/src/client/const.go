package main

import (
	"fmt"
	"math/rand"
)

var GmTips = []string{
	fmt.Sprintf("%-10v:%-10v", "GM", "example"),
	fmt.Sprintf("%-10v:%-10v", "moveto", "moveto 500 500"),
	" ",
	fmt.Sprintf("%-10v:%-10v", "skill", "desc"),
	fmt.Sprintf("%-10v:%-10v", "1", "normal attack"),
	fmt.Sprintf("%-10v:%-10v", "2", "reduce defence"),
}

var rgbs [][3]uint8

func init() {
	for len(rgbs) < 255 {
		r := uint8(rand.Int31() % 255)
		g := uint8(rand.Int31() % 255)
		b := uint8(rand.Int31() % 255)
		rgbs = append(rgbs, [3]uint8{r, g, b})
	}
}

func GetRGB(i int) (uint8, uint8, uint8) {
	if i < 0 || i >= len(rgbs) {
		return 0, 0, 0
	}
	return rgbs[i][0], rgbs[i][1], rgbs[i][2]
}
