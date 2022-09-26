package main

import (
	"fmt"

	"github.com/btwiuse/ufo"
)

func main() {
	showIdna("👽")
	showIdna("I😏")
	showIdna("i😏.ws")
	showIdna("i❤️")
	showIdna("i❤️")
	showIdna("❤️")
	showIdna("i❤.ws")
	showIdna("i❤")
	showIdna("i❤:80")
	showIdna("sudo")
	showIdna("https://😂.ufo.k0s.io")
}

func showIdna(s string) {
	t := ufo.ToIdna(s)
	fmt.Println(s, "~>", t)
}
