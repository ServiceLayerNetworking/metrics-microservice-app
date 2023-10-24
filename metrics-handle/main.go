package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type processorResp struct {
	Anomalies string `json:"anomaliesDetected"`
}

type processorReq struct {
	Metrics []byte `json:"metrics"`
}

func main() {
	r := gin.Default()
	r.GET("/start", DetectAnomalies)
	r.Run(":8080")
}

func DetectAnomalies(c *gin.Context) {
	// call metrics processing service
	// return whether anomaly was detected
	data := make([]byte, 5*1024*1024)
	rand.Read(data)
	//r := &processorReq{
	//	Metrics: data,
	//}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://metrics-processing:8000/detectAnomalies", bytes.NewBuffer(data))
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
