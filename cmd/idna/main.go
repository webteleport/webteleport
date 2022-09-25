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
	showIdna("sudo")
}

func showIdna(s string) {
	ascii, err := idna.ToASCII(s)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(ascii, s)
}
