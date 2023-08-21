package client

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

type Topic struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string
	Author     string
	URL        string
	Date       string
	Comments   int
	UpdateTime time.Time
}

// TopicContent is for new feature in the future, store the related topic content
type TopicContent struct {
	TopicID primitive.ObjectID `bson:"_id,omitempty"`
	Content string
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

func parseOnePage(c *colly.Collector, url string, topicsList []Topic) ([]Topic, string, error) {
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
			Title:      title,
			Author:     author,
			URL:        e.Request.AbsoluteURL(url),
			Date:       date,
			Comments:   parseComments(comments),
			UpdateTime: time.Now(),
		}
		topicsList = addTopic(topicsList, topic)

	})
	c.OnHTML("#action-bar-container > div > div.btn-group.btn-group-paging > a:nth-child(2)", func(e *colly.HTMLElement) {
		nextPage = e.Request.AbsoluteURL(e.Attr("href"))
	})

	// Start crawling
	err := c.Visit(url)
	if err != nil {
		return nil, nextPage, fmt.Errorf("visit err: %v", err)
	}
	//fmt.Printf("Total topic: %d\n", len(topicsList))
	return topicsList, nextPage, nil
}

func GetLatestBoardData(url string) ([]Topic, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +https://www.google.com/bot.html)"),
	)
	err := c.Limit(&colly.LimitRule{DomainGlob: "*", RandomDelay: time.Second * 1})
	if err != nil {
		return nil, fmt.Errorf("set limit err: %v", err)
	}

	topics := make([]Topic, 0)

	for len(topics) <= 50 && url != "" {
		topics, url, err = parseOnePage(c, url, topics)
		if err != nil {
			return nil, fmt.Errorf("parse page err: %v", err)
		}
	}
	return topics, nil

	//Write json result to standard output
	//enc := json.NewEncoder(os.Stdout)
	//enc.SetIndent("", "  ")
	//
	//// Dump json to the standard output
	//enc.Encode(topics)

	////----Write result to file----
	//file, err := os.OpenFile("ptt_title.json", os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	fmt.Println(fmt.Errorf("open file err: %v", err))
	//}
	//defer file.Close()
	//enc := json.NewEncoder(file)
	//enc.SetIndent("", "  ")
	//
	//// Dump json to the output file
	//_ = enc.Encode(topics)
	//----------
}
