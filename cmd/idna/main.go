package main

import (
	"fmt"
	"log"

	"golang.org/x/net/idna"
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
	ascii, err := idna.ToASCII(s)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(ascii, s)
}
