package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type processorResp struct {
	Anomalies string `json:"anomaliesDetected"`
}

type processorReq struct {
	Metrics []byte `json:"metrics"`
}

var propogate = []string{"X-Request-Id", "X-B3-Traceid", "X-B3-Spanid", "X-B3-ParentSpanid", "X-B3-Sampled"}

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
	// propogate tracing headers

	for _, v := range propogate {
		if val, ok := c.Request.Header[v]; ok {
			if v == "X-B3-Traceid" {
				fmt.Printf("traceid: %v\n", val)
			}
			req.Header[v] = val
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// save body to string
	var r *processorResp
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r = &processorResp{
			Anomalies: "false",
		}
	} else {
		r = &processorResp{}
		if err = json.Unmarshal(body, r); err != nil {
			panic(err)
		}
	}
	fmt.Printf("detected: %v\n", r.Anomalies)
	c.JSON(200, gin.H{
		"anomaliesDetected": r.Anomalies,
	})
}
