package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	links "github.com/oucema001/link"
)

const (
	// A generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

func main() {
	urlFlag := flag.String("url", "https://www.calhoun.io", "url to use for the sitemap")
	depth := flag.Int("depth", 2, "depth to traverse for the sitemap")
	filename := flag.String("filename", "sitemap.xml", "file name to output sitemap to")
	flag.Parse()

	l := siteTraverse(*depth, *urlFlag)

	BuildSiteMapXML(*filename, l)
}

type URLSiteMap struct {
	//URL string `xml:"url"`
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns:prefix,attr"`
	URL     []u      `xml:"url"`
}

type u struct {
	Loc string `xml:"loc"`
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
				if _, ok := res[k]; ok {
					continue
				}
				res[v] = struct{}{}
				nq[v] = struct{}{}
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

func BuildSiteMapXML(filename string, l []string) {
	m := make([]u, 0)
	for _, v := range l {
		if v == "" {
			continue
		}
		k := u{
			v,
		}
		m = append(m, k)
	}
	u := URLSiteMap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL:   m,
	}
	//b, err := xml.Marshal(u)
	var bz []byte
	bz = []byte(xml.Header)
	b, err := xml.MarshalIndent(u, "", " ")
	if err != nil {
		log.Println(err)
	}
	bz = append(bz, b...)
	err = ioutil.WriteFile(filename, bz, 0644)
	if err != nil {
		log.Println(err)
	}
}
