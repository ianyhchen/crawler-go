package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
	"strconv"
)

type Topic struct {
	Title    string
	Author   string
	URL      string
	Date     string
	Comments int
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

	c := colly.NewCollector()

	topics := make([]Topic, 0)
	//// Set the "over18" cookie in the request headers
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "over18=1")
	})
	// Configure collector options
	c.OnHTML(".r-ent", func(e *colly.HTMLElement) {
		if len(topics) >= 50 {
			return
		}

		title := e.ChildText(".title a")
		author := e.ChildText(".meta .author")
		url := e.ChildAttr(".title a", "href")
		date := e.ChildText(".meta .date")
		comments := e.ChildText(".nrec span")
		topic := Topic{
			Title:    title,
			Author:   author,
			URL:      url,
			Date:     date,
			Comments: parseComments(comments),
		}
		topics = append(topics, topic)
		//c.Visit(url)
	})
	// Extract additional information from the topic page
	//c.OnHTML("#main-content", func(e *colly.HTMLElement) {
	//	// Extract and process additional information from the topic page
	//	// Modify this part to extract the desired information from the topic page
	//})

	// Start crawling
	err := c.Visit(mainUrl)

	if err != nil {
		fmt.Println(fmt.Errorf("visit err: %v", err))
	}
	fmt.Printf("Total topic: %d\n", len(topics))
	//for _, topic := range topics {
	//	fmt.Println("Title:", topic.Title)
	//	fmt.Println("Author:", topic.Author)
	//	fmt.Println("URL:", topic.URL)
	//	fmt.Println("Date:", topic.Date)
	//	fmt.Println("Comments:", topic.Comments)
	//	fmt.Println("--------------")
	//}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(topics)
}
