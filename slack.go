package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// URLReq is request.
type URLReq struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Event     Event  `json:"event"`
	Text      string `json:"text"`
}

// Event ...
type Event struct {
	Type    string `json:"type"`
	FileID  string `json:"file_id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// URLRes is response.
type URLRes struct {
	Challenge string `json:"challenge"`
}

// FileReq is request.
type FileReq struct {
	Token string `json:"token"`
	File  string `json:"file"`
}

// FileRes is request.
type FileRes struct {
	OK   bool `json:"ok"`
	File File `json:"file"`
}

// File is request.
type File struct {
	Name string `json:"name"`
	URL  string `json:"url_private_download"`
	Type string `json:"mimetype"`
}

// ReqCmd is request.
type ReqCmd struct {
	Token     string `json:"token"`
	ChannelID string `json:"channel_id"`
	Command   string `json:"command"`
	Text      string `json:"text"`
	ResURL    string `json:"response_url"`
}

// ResCmd is request.
type ResCmd struct {
	Chnl   string `json:"channel"`
	Text   string `json:"text"`
	Parse  string `json:"parse"`
	AsUsr  bool   `json:"as_user"`
	UnfURL bool   `json:"unfurl_links"`
}

const (
	// reqTypeURL ...
	reqTypeURL string = "url_verification"
	// reqTypeEvt ...
	reqTypeEvt string = "event_callback"
	// evtTypeShr
	evtTypeShr string = "file_shared"
)

// EventCallback ...
func EventCallback(req URLReq) {
	logText := fmt.Sprintf("Event info, event_type:%s file_id:%s", req.Event.Type, req.Event.FileID)
	fmt.Println(logText)
	// check.
	if (req.Event.FileID == "") || (req.Event.Type != evtTypeShr) {
		fmt.Println(fmt.Sprintf("type:%s text:%s file_iD:%s ", req.Event.Type, req.Event.Text, req.Event.FileID))
		return
	}
	// req params.
	var fileReq FileReq
	fileReq.Token = config.Slack.BotToken
	fileReq.File = req.Event.FileID
	result := FileUpload(fileReq)
	// send chat.
	SendChat(result)
}

// FileUpload ...
func FileUpload(fileReq FileReq) (result string) {
	logText := fmt.Sprintf("File download info, api_url:%s token:%s", config.Slack.FileURL, fileReq.Token)
	fmt.Println(logText)
	// get url of file.
	res := GetFileURL(config.Slack.FileURL, fileReq.Token, fileReq.File)
	if !res.OK {
		result = "fail download."
		fmt.Println(result)
		return result
	}
	logText = fmt.Sprintf("file_url:%s", res.File.URL)
	fmt.Println(logText)
	// file operation.
	result = OperateFile(res.File.URL, res.File.Name)
	fmt.Println(result)
	return result
}

// GetFileURL ...
func GetFileURL(url, token, file string) (result *FileRes) {
	req, _ := http.NewRequest("GET", url, nil)
	// set params.
	params := req.URL.Query()
	params.Add("token", token)
	params.Add("file", file)
	req.URL.RawQuery = params.Encode()
	fmt.Println(req.URL.String())
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
	fmt.Println("GetFileURL response:" + string(resBody))
	return resJSON
}

// DispFileList ...
func DispFileList() string {
	files, err := ioutil.ReadDir(config.File.Path)
	if err != nil {
		panic(err)
	}
	var paths string = ""
	for _, file := range files {
		if file.IsDir() {
			paths += DirWalk(file.Name()) + "\n"
			continue
		}
	}
	fmt.Println(paths)
	return paths
}

// SendChat ...
func SendChat(msg string) {
	// set body.
	var res ResCmd
	res.Chnl = config.Slack.Channel
	res.Text = msg
	res.Parse = "none"
	res.AsUsr = true
	res.UnfURL = true
	body, _ := json.Marshal(res)
	fmt.Println(fmt.Sprintf("request body:%s", string(body)))
	// send.
	SendPost(body)
}
