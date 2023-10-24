package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/detectAnomalies", DetectAnomalies)
	r.Run(":8080")

}

func DetectAnomalies(c *gin.Context) {
	promUrl := "http://10.244.0.31:9090/api/v1/query?query=istio_requests_total"
	resp, err := http.Get(promUrl)
	if err != nil {
		c.JSON(500, gin.H{
			"anomaliesDetected": "error",
		})
		return
	}
	// save body to string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"anomaliesDetected": "error",
		})
		return
	}
	fmt.Printf("Prometheus response code: %v, body len: %v\n", resp.StatusCode, len(body))

	// we'll assume there are always no anomalies
	c.JSON(500, gin.H{
		"anomaliesDetected": "false",
	})
}
