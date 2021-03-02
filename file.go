package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// OperateFile ...
func OperateFile(fURL, name string) string {
	var logText string
	// download file.
	body := GetFile(fURL)

	// create dir.
	nowtime := time.Now()
	hash := CreateHash(nowtime.String())
	dir := filepath.Join(config.File.Path, hash)
	fmt.Println("dir path:" + dir)
	if !CreateDir(dir) {
		//
		logText = "dir create error."
		fmt.Println(logText)
		return logText
	}
	// join path.
	fullPath := filepath.Join(dir, name)
	fmt.Println("file path:" + fullPath)

	// create file.
	file := CreateFile(fullPath)
	if file == nil {
		logText = "file create error."
		fmt.Println(logText)
		return logText
	}

	// copy file
	_, err := file.Write(body)
	defer file.Close()
	if err != nil {
		logText = "file write error."
		fmt.Println(logText, err)
		return logText
	}
	//baseURL, _ := url.Parse(config.File.DLURL)
	//baseURL.Path = path.Join(baseURL.Path, hash, name)
	upURL := config.File.DLURL + "/" + hash + "/" + name
	return fmt.Sprintf("<%s|%s>", upURL, upURL)
}

// GetFile ...
func GetFile(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	// set params.
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Slack.BotToken))
	fmt.Println(req.URL.String())
	// do.
	client := &http.Client{}
	res, _ := client.Do(req)
	// respons body.
	resBody, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	//fmt.Println("GetFile response:" + string(resBody))
	return resBody
}

// CreateHash ...
func CreateHash(name string) (hash string) {
	sha1 := sha1.Sum([]byte(name))
	hash = fmt.Sprintf("%x", sha1)
	fmt.Println("hash:" + hash)
	return hash
}

// CreateDir ...
func CreateDir(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		fmt.Println("dir exists." + path)
		return true
	}
	if err := os.Mkdir(path, 0755); err != nil {
		fmt.Println("dir create error.", err)
		return false
	}
	return true
}

// CreateFile ...
func CreateFile(path string) *os.File {
	_, err := os.Stat(path)
	if err == nil {
		fmt.Println("file is exists.", err)
		return nil
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("file create error.", err)
		return nil
	}

	return file
}

// DirWalk ...
func DirWalk(dir string) (path string) {
	files, err := ioutil.ReadDir(filepath.Join(config.File.Path, dir))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			path = DirWalk(filepath.Join(dir, file.Name()))
			continue
		} else {
			path = filepath.Join((filepath.Base(dir)), file.Name())
		}
	}

	return path
}

// DelFile ...
func DelFile(target string) (result string) {
	// check.
	if target == "" {
		result = fmt.Sprintf("[dir name/file name] is empty.")
		fmt.Println(result)
		return result
	}
	// exists.
	fullPath := filepath.Join(config.File.Path, target)
	fInfo, err := os.Stat(fullPath)
	if err != nil {
		result = fmt.Sprintf("file not exists. %s", target)
		fmt.Println(result)
		return result
	}
	// is dir.
	var dirPath string
	if fInfo.IsDir() {
		dirPath = fullPath
	} else {
		dirPath = filepath.Dir(fullPath)
	}
	fmt.Println(fmt.Sprintf("dir path:%s", dirPath))
	// delete.
	if err := os.RemoveAll(dirPath); err != nil {
		result = fmt.Sprintf("failed to delete the file. %s", target)
		fmt.Println(result, err)
		return result
	}
	result = fmt.Sprintf("succesed to delete the file. %s", target)
	return result
}
