package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)


func Fetch(url string, ch chan string)  {
	start := time.Now()
	if strings.HasPrefix(url, "http://") == false {
		url = "http://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("fetch: %v\n", err)
		return
	}
	defer resp.Body.Close()
	byteCnt, err := io.Copy(ioutil.Discard, resp.Body)
	ch <- fmt.Sprintf("%.2f %s %d %d", time.Since(start).Seconds(), url, byteCnt, resp.StatusCode)
}


func Gets(urls []string) {
	start := time.Now()
	ch := make(chan string)
	for _, url := range urls {
		go Fetch(url, ch)
	}
	for range urls {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2f elapsed", time.Since(start).Seconds())
}