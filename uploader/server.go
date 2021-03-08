package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func service() {
	router := gin.Default()
	router.Static("/assets", "./assets")

	// GET
	router.GET("/file-uploader", RecvGet)
	// POST
	router.POST("/file-uploader", RecvPostUp)
	router.POST("/file-cmd", RecvPostCmd)

	router.Run(":" + config.Service.Port)
}

// RecvGet ...
func RecvGet(ctx *gin.Context) {
	logText := "GET."
	fmt.Println(logText)
}

// RecvPostUp ...
func RecvPostUp(ctx *gin.Context) {
	// recv data to json.
	var req URLReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logText := fmt.Sprintf("Recv POST, type:%s token:%s", req.Type, req.Token)
	fmt.Println(logText)
	// type.
	if req.Type == reqTypeURL {
		ctx.JSON(http.StatusOK, URLRes{req.Challenge})
	} else if req.Type == reqTypeEvt {
		EventCallback(req)
	} else {
		//
	}
	fmt.Println()
}

// RecvPostCmd ...
func RecvPostCmd(ctx *gin.Context) {
	fmt.Println("Recv POST file-cmd.")
	// recv data to json.
	var req ReqCmd
	req.ChannelID = ctx.PostForm("channel_id")
	req.Command = ctx.PostForm("command")
	req.Text = ctx.PostForm("text")
	logText := fmt.Sprintf("Recv POST, ch_id:%s cmd:%s text:%s", req.ChannelID, req.Command, req.Text)
	fmt.Println(logText)
	// exclude all but certain channels.
	if req.ChannelID != config.Slack.Channel {
		return
	}
	// type.
	if req.Command == "/list" {
		// list of files.
		result := DispFileList()
		// send chat.
		msg := "List of Uploaded Files.\n" + result
		SendChat(msg)
	} else if req.Command == "/del" {
		msg := DelFile(req.Text)
		SendChat(msg)
	} else {
		//
	}
	fmt.Println()
}

// SendPost ...
func SendPost(body []byte) {
	fmt.Println(fmt.Sprintf("request body:%s", string(body)))
	// ChatURL
	req, _ := http.NewRequest("POST", config.Slack.ChatURL, bytes.NewBuffer(body))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Slack.BotToken))
	req.Header.Add("Content-Type", "application/json")
	fmt.Println(fmt.Sprintf("request:%s", req.URL.String()))
	// do.
	client := &http.Client{}
	res, _ := client.Do(req)
	// respons body.
	resBody, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	resJSON := new(FileRes)
	if err := json.Unmarshal(resBody, resJSON); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return
	}
	fmt.Println("SendChat response:" + string(resBody))
	return
}
