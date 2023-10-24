package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type processorResp struct {
	Anomalies string `json:"anomaliesDetected"`
}

func main() {
	r := gin.Default()
	r.GET("/start", DetectAnomalies)
	r.Run(":8080")
}

func DetectAnomalies(c *gin.Context) {
	// call metrics processing service
	// return whether anomaly was detected
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://metrics-processing:8000/detectAnomalies", nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// save body to string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	r := &processorResp{}
	if err = json.Unmarshal(body, r); err != nil {
		panic(err)
	}
	fmt.Printf("detected: %v\n", r.Anomalies)
	c.JSON(200, gin.H{
		"anomaliesDetected": r.Anomalies,
	})
}
