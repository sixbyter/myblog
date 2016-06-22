package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	rootPath     = flag.String("r", "./public", "where is the website root path?")
	articlesPath = flag.String("a", "./articles", "where is the articles path?")
)

func main() {
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir(*rootPath)))

	http.HandleFunc("/api/articles", ArticleIndex)
	http.HandleFunc("/api/article", ArticleShow)
	http.HandleFunc("/test", test)

	port := "8888"
	fmt.Println("开启http服务, 端口是: " + port + "网站根目录为: " + *rootPath)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func test(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Del("Content-Length")
	resp.Header().Del("Content-Type")
	resp.Header().Del("Date")
}
