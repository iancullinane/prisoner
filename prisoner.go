package main

import (
	"flag"
	"log"
)

var value string

func init() {
	flag.StringVar(&value, "value", "", "")
}

func main() {

	log.Println("Prisoners dilemna")
	// fmt.Printf("%s\n", dst[:n])

}
