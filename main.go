package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	link "github.com/oucema001/link"
)

func main() {
	f := flag.String("url", "caluhon.io", "url to use for the sitemap")
	flag.Parse()
}

func getLinks(url string) {
	content, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer content.Body.Close()
	l, err := link.Parse(content.Body)
	if err != nil {
		log.Println(err)
	}
	for k, v := range l {

	}
	fmt.Sprintf("%s : %s \n")
}
