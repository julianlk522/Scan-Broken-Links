package main

import (
	"fmt"
	"strings"
	"time"

	screen "github.com/aditya43/clear-shell-screen-golang"
	"github.com/fatih/color"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

// text styles
var title_text_styles *color.Color = color.New(color.FgHiWhite, color.Underline)
var highlight_text_styles *color.Color = color.New(color.FgRed, color.Bold)

func main() {
	// collect Root URL and desired number of consumer
	// threads from user
	root_url, ct := collect_user_input()

	// start timing
	start := time.Now()

	// scan site
	bad_links := scan_site(root_url, ct)

	// report results
	report_results(bad_links)

	// end timing
	fmt.Printf("\n(Total execution time: %v)\n", time.Since(start))
	return
}

func collect_user_input() (string, int) {
	// clear terminal
	clear_chars()
	
	// ask for URL
	var root_url string = "https://www.google.com"
	fmt.Printf("Enter site root URL\n\n(or just hit 'Enter' for default: %s)\n\n: ", root_url)

	// overwrite root_url if user inputs one
	fmt.Scanf("%s\n", &root_url)

	// clear
	clear_chars()

	// ask for number of consumer threads
	var ct int = 2
	fmt.Printf("Enter number of threads\n\n(or just hit 'Enter' for default: %d)\n\n: ", ct)

	// overwrite ct if user inputs one
	fmt.Scanf("%d\n", &ct)

	// clear
	clear_chars()

	return root_url, ct
}

func clear_chars() {
	screen.Clear()
	screen.MoveTopLeft()
}

func scan_site(root_url string, ct int) []string {
	title_text_styles.Printf("Scanning: %s...\n", root_url)

	c := colly.NewCollector(
		// max depth 4 is arbitrary... sets reasonable limit
		colly.MaxDepth(4),
	)

	// new queue with selected number of consumer threads
	// and default queue memory
	q, err := queue.New(ct, &queue.InMemoryQueueStorage{MaxSize: 10000})
	q.AddURL(root_url)
	if err != nil {
		panic(err)
	}

	// prepare slices for storing URLs
	visited_urls := []string{root_url + "/"}
	bad_links := []string{}

	c.OnResponse(func(r *colly.Response) {
		// if bad link
		if r.StatusCode >= 400 && r.StatusCode <= 499 {

			// highlight URL
			highlight_text_styles.Printf("Bad link: %s (response code: %d)\n", r.Request.URL, r.StatusCode)

			// save it to bad_links
			bad_links = append(bad_links, fmt.Sprint(r.Request.URL))
		} else {
			// else print out URL
			color.HiBlue(fmt.Sprint(r.Request.URL))
		}

		// save visited URLs
		visited_urls = append(visited_urls, fmt.Sprint(r.Request.URL))
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// skip if current URL is outside root URL domain
		if !strings.Contains(e.Request.URL.String(), root_url) {
			return
		}

		// skip if already visited
		for _, v := range visited_urls {
			if v == e.Attr("href") {
				return
			}
		}

		// else add new linked page to queue
		q.AddURL(e.Attr("href"))
	})

	q.Run(c)
	return bad_links
}

func report_results(bad_links []string) {
	if len(bad_links) == 0 {
		title_text_styles.Println("\nAll links up to date.")
	} else {
		title_text_styles.Println("\nNeeds updating:")
	
		// print bad links
		for _, v := range bad_links {
			highlight_text_styles.Println(v)
		}
	}
}