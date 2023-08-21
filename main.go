package main

import (
	"crawler_go/client"
	"crawler_go/pkg/ptt"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

var service *client.MongoDBService // Declare the MongoDBService instance as a package-level variable
func main() {
	var err error
	service, err = client.NewMongoDBService()
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	router := gin.Default()
	//// Use your custom logger middleware
	//router.Use(client.CustomLoggerMiddleware())
	logger := client.SetupCustomLogger()
	// Define your routes here
	router.GET("/board/:name", func(c *gin.Context) {
		boardName := c.Param("name")
		count := c.DefaultQuery("count", strconv.Itoa(ptt.DefaultGetDataLimit))
		getBoardData(c, logger, boardName, count) // Pass the logger to the handler function
	})
	router.POST("/board/:name/update", func(c *gin.Context) {
		boardName := c.Param("name")
		updateLatestBoardData(c, logger, boardName)
	})

	err = router.Run(":8080")
	if err != nil {
		fmt.Println("Fail to start api server.")
	}
}

func getBoardData(c *gin.Context, logger *log.Logger, boardName string, count string) {
	// Implement the logic to retrieve board data from MongoDB
	// and respond with the data in the HTTP response.
	// You may need to marshal the data into JSON before sending it.
	if !ptt.ValidateBoardName(boardName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Board not found"})
		return
	}
	counter, err := strconv.Atoi(count)
	if err != nil {
		counter = ptt.DefaultGetDataLimit
	}
	data, err := service.GetBoardDataFromDB(boardName, int64(counter))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//// Retrieve the custom logger from the Gin context
	//logger, _ := c.Get("logger")
	//customLogger, ok := logger.(*log.Logger)
	//if !ok {
	//	// Fallback to the default logger if something goes wrong
	//	customLogger = log.Default()
	//}
	logger.Printf("Retrieved %d topics for board: %s\n", len(data), boardName)

	c.JSON(http.StatusOK, data)
}

func updateLatestBoardData(c *gin.Context, logger *log.Logger, boardName string) {
	// Implement the logic to update board data from PTT and
	// store it in MongoDB.
	if !ptt.ValidateBoardName(boardName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Board not found"})
		return
	}
	URL := ptt.GetBoardURL(boardName)
	if URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Board not found"})
		return
	}
	topicData, err := client.GetLatestBoardData(URL)
	if err != nil {
		logger.Printf("Error retrieving board data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve board data"})
		return
	}
	if len(topicData) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No new data found"})
		return
	}
	insertCount, updateCount := service.UpdateTopic(boardName, topicData)
	logger.Printf("Retrieved %d topic, insert %d new topics, update %d topics", len(topicData), insertCount, updateCount)
	c.JSON(http.StatusOK, gin.H{"message": "Board data updated successfully"})
}
