package main

import (
	"fmt"
	log "github.com/llimllib/loglevel"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"os"
	"strings"
	"time"
)

type Link struct {
	url string
	text string
	depth int
}

type HttpError struct {
	original string
}

func createNewLink(tag html.Token, text string, depth int) Link {
	link := Link{text: strings.TrimSpace(text), depth: depth}

	for i := range tag.Attr {
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)
		}
	}
	return link
}

func linkReader(resp *http.Response, depth int) []Link {
	page := html.NewTokenizer(resp.Body)
	var links []Link

	var start *html.Token
	var text string

	for {
		_ = page.Next()
		token := page.Token()
		if token.Type == html.ErrorToken {
			break
		}
		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", token, token.Data)
		}
		if token.DataAtom == atom.A {
			switch token.Type {
			case html.StartTagToken:
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken:
				if start == nil {
					log.Warn("We found the end of the link but there was no start")
					continue
				}
				link := createNewLink(*start, text, depth)
				if link.Valid() {
					links = append(links, link)
					log.Debug("Link found")
				}

				start = nil
				text = ""

			}
		}
	}

	log.Debug(links)
	return links
}

func (link Link) String() string {
	spacer := strings.Repeat("\t", link.depth)
	return fmt.Sprintf("%s%s (%d) - %s", spacer, link.text, link.depth, link.url)
}

var MaxDepth = 2

func (link Link) Valid() bool {
	if link.depth >= MaxDepth {return false}
	if len(link.text) == 0 {return false}
	if len(link.url) == 0 || strings.Contains(strings.ToLower(link.url), "javascript") {
		return false
	}
	return true
}

func (httpErr HttpError) Error() string {return httpErr.original}

func recurDownloader(url string, depth int, channel chan bool) {
	page, err := get(url)
	if err != nil {
		log.Error(err)
		channel <- false
		return
	}
	links := linkReader(page, depth)
	for _, link := range links {
		fmt.Println(link)
		if depth + 1 < MaxDepth {
			recurDownloader(link.url, depth + 1, channel)
		}
	}
	channel <- true
}

func get(url string) (resp *http.Response, err error) {
	log.Debug("Downloading %s", url)
	resp, err = http.Get(url)
	if err != nil {
		log.Debug("Error $s", err)
		return
	}
	if resp.StatusCode > 299 {
		err = HttpError{fmt.Sprintf("Error (%d) %s", resp.StatusCode, url)}
		log.Debug(err)
		return
	}
	return
}

func main() {
	println("Starting...")
    e := log.SetPriorityString("info")
    if e != nil {println("unable to find priority")} // by default, makes no sense
    log.SetPrefix("crawler")

    log.Debug(os.Args)
    if len(os.Args) < 2 {
    	log.Fatalln("Missing url arg")
	}

	doneStatus := make(chan bool)

	go recurDownloader(os.Args[1], 0, doneStatus)
	select {
		case res := <-doneStatus:
			switch res {
			case true:
				fmt.Println("No failures")
			case false:
				fmt.Println("Error")
			}
		case <-time.After(60 * time.Second):
			fmt.Println("Timeout one minute reached")
	}
	fmt.Println("Done")

}
