package controllers

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
	router.GET("/api/finance/record", Finance.Records)
	router.POST("/api/finance/record", Finance.Record)
	router.GET("/api/finance/account", Finance.Accounts)
	router.POST("/api/finance/account", Finance.Account)
}

func get(uri string, params string, router *gin.Engine) (int, []byte) {
	req := httptest.NewRequest("GET", uri+"?"+params, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	body, _ := ioutil.ReadAll(result.Body)
	return result.StatusCode, body
}

func post(uri string, jsonSS string, router *gin.Engine) (int, []byte) {
	req := httptest.NewRequest("POST", uri, bytes.NewBuffer([]byte(jsonSS)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	body, _ := ioutil.ReadAll(result.Body)
	return result.StatusCode, body
}
