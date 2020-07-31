package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	links "github.com/oucema001/link"
)

func main() {
	f := flag.String("url", "https://www.calhoun.io", "url to use for the sitemap")
	flag.Parse()
	_ = f
	//getLinks("https://www.calhoun.io")
	siteTraverse(3, "https://www.calhoun.io")
}

func getLinks(urlStr string) []string {
	content, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("the error is ", urlStr)
		log.Println(err)
	}
	defer content.Body.Close()

	l, err := links.Parse(content.Body)
	if err != nil {
		log.Println(err)
	}
	var list []string
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("the error is ", urlStr)
		log.Fatal(err)
	}
	domain := u.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	list = make([]string, 1)
	for _, v := range l {
		link := v.Href
		if strings.HasPrefix(link, "/") {
			list = append(list, "http://"+domain+link)
			continue
		}
		ur, err := url.Parse(link)
		if err != nil {
			continue
		}
		domain2 := ur.Hostname()
		if domain != domain2 {
			continue
		}

		list = append(list, link)
	}
	return list
}

func siteTraverse(depth int, urlStr string) []string {
	res := make(map[string]struct{})

	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}

	for i := 0; i < depth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for k := range q {
			if k == "" {
				continue
			}
			list := getLinks(k)
			for _, v := range list {
				if _, ok := q[v]; !ok {
					nq[v] = struct{}{}
					res[v] = struct{}{}
				}
			}
		}
		i++

	}
	ret := make([]string, 1)

	for k := range res {
		ret = append(ret, k)
	}
	return ret
}
