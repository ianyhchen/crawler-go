package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
	"strconv"
	"time"
)

type Topic struct {
	Title    string
	Author   string
	URL      string
	Date     string
	Comments int
}

func addTopic(topics []Topic, topic Topic) []Topic {
	return append(topics, topic)
}

func parseComments(comments string) int {
	if comments == "" {
		return 0
	}
	count, err := strconv.Atoi(comments)
	if err != nil {
		return 0
	}
	return count
}

func parseOnePage(c *colly.Collector, url string, topicsList []Topic) ([]Topic, string) {
	var nextPage string
	// Set the "over18" cookie in the request headers
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "over18=1")
	})
	// Configure collector options
	c.OnHTML(".r-ent", func(e *colly.HTMLElement) {
		title := e.ChildText(".title a")
		author := e.ChildText(".meta .author")
		url := e.ChildAttr(".title a", "href")
		date := e.ChildText(".meta .date")
		comments := e.ChildText(".nrec span")
		topic := Topic{
			Title:    title,
			Author:   author,
			URL:      e.Request.AbsoluteURL(url),
			Date:     date,
			Comments: parseComments(comments),
		}
		topicsList = addTopic(topicsList, topic)

	})
	c.OnHTML("#action-bar-container > div > div.btn-group.btn-group-paging > a:nth-child(2)", func(e *colly.HTMLElement) {
		nextPage = e.Request.AbsoluteURL(e.Attr("href"))
	})

	// Start crawling
	err := c.Visit(url)

	if err != nil {
		fmt.Println(fmt.Errorf("visit err: %v", err))
	}
	fmt.Printf("Total topic: %d\n", len(topicsList))
	return topicsList, nextPage
}

func main() {
	var mainUrl = "https://www.ptt.cc/bbs/Gossiping/index.html"
	//var preCheckUrl = "https://www.ptt.cc/ask/over18"
	//// Create an HTTP client with a session
	//client := &http.Client{}
	//
	//// Create a POST request to set the "over18" cookie
	//data := url.Values{}
	//data.Set("from", "/bbs/Gossiping/index.html")
	//data.Set("yes", "yes")
	//req, err := http.NewRequest("POST", preCheckUrl, nil)
	//if err != nil {
	//	fmt.Println("Failed to create request:", err)
	//	return
	//}
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.PostForm = data
	//
	//// Send the POST request to set the cookie
	//res, err := client.Do(req)
	//fmt.Println("Ask 18 response:", res)
	//if err != nil {
	//	fmt.Println("Failed to set cookie:", err)
	//	return
	//}
	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//		fmt.Println("Failed to close response body:", err)
	//	}
	//}(res.Body)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +https://www.google.com/bot.html)"),
	)
	err := c.Limit(&colly.LimitRule{Delay: time.Second})
	if err != nil {
		fmt.Println(fmt.Errorf("set limit err: %v", err))
	}

	topics := make([]Topic, 0)
	var url = mainUrl
	for len(topics) <= 50 && url != "" {
		topics, url = parseOnePage(c, url, topics)
	}
	//Write json result to standard output
	//enc := json.NewEncoder(os.Stdout)
	//enc.SetIndent("", "  ")
	//
	//// Dump json to the standard output
	//enc.Encode(topics)

	//Write result to file
	file, err := os.OpenFile("ptt_title.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(fmt.Errorf("open file err: %v", err))
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the output file
	_ = enc.Encode(topics)
}
