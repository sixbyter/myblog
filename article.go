package main

import (
	"bufio"
	"encoding/json"
	_ "github.com/realint/dbgutil"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"
)

type Article struct {
	Filename  string `json:"filename"`
	Date      string `json:"date"`
	Title     string `json:"title"`
	Timestamp int64  `json:"timestamp"`
}

type ArticleSlice []Article

func (a ArticleSlice) Len() int {
	return len(a)
}

func (a ArticleSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ArticleSlice) Less(i, j int) bool {
	return a[j].Timestamp < a[i].Timestamp
}

func ArticleIndex(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		articles := getArticles()
		b, err := json.Marshal(articles)
		checkError(err)
		resp.Write(b)
	} else {
		http.NotFound(resp, req)
	}
}

func getArticles() ArticleSlice {
	files, err := ioutil.ReadDir(*articlesPath)
	checkError(err)
	articles := make(ArticleSlice, len(files))
	for i, file := range files {
		f, err := os.Open(*articlesPath + "/" + file.Name())
		checkError(err)
		defer f.Close()
		b := bufio.NewReader(f)
		line, err := b.ReadBytes('\n')
		articles[i].Title = string(line[13 : len(line)-1])
		articles[i].Filename = file.Name()
		articles[i].Date = string(line[3:12])
		the_time, err := time.Parse("02 Jan 06", articles[i].Date)
		if err == nil {
			articles[i].Timestamp = the_time.Unix()
		}
	}
	// 排序
	sort.Sort(articles)
	return articles
}

func ArticleShow(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {

		name := req.FormValue("name")
		filename := *articlesPath + "/" + name
		_, err := os.Stat(filename)
		if err != nil {
			resp.Write([]byte(`{"status":404,"message":"article not found."}`))
		}
		article := make(map[string]string)
		content, err := ioutil.ReadFile(filename)
		// dbgutil.FormatDisplay("content", content)
		article["content"] = string(content)
		b, err := json.Marshal(article)
		checkError(err)
		resp.Write(b)
	} else {
		http.NotFound(resp, req)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
