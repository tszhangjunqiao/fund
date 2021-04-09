package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Fund struct {
	Code   string `json:"fundcode"`
	Name   string `json:"name"`
	Change string `json:"gszzl"`
}

type Config struct {
	Code []string `json:"code"`
	Pushplus string `json:"pushplus"`
}

func main() {
	config := getConfig()
	fund := config.Code
	str := ""

	for _, code := range fund {
		url := fmt.Sprintf("https://fundgz.1234567.com.cn/js/%s.js?rt=%d", code, time.Now().Unix()*1000)
		client := http.Client{}

		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36 Edg/88.0.705.63")
		resp, _ := client.Do(request)

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		result := string(body)
		result = strings.Replace(result, "jsonpgz(", "", -1)
		result = strings.Replace(result, ");", "", -1)

		r := &Fund{}
		_ = json.Unmarshal([]byte(result), r)

		change, _ := strconv.ParseFloat(r.Change, 32)
		color := "red"
		if change < 0 {
			color = "green"
		}

		chaColor := fmt.Sprintf("<font color=\"%s\">%s</font>", color, r.Change)
		rs := fmt.Sprintf("%s(%s)，涨跌幅：%s；\n", r.Name, r.Code, chaColor)
		str += rs

		randSleep()
	}

	sendPushPlus(config.Pushplus , "今日行情", str)
}

func sendPushPlus(token, title, content string) {
	url := "https://www.pushplus.plus/send"
	ma := make(map[string]interface{})
	ma["token"] = token
	ma["title"] = title
	ma["content"] = content
	js, _ := json.Marshal(ma)
	param := bytes.NewReader(js)

	req, _ := http.NewRequest("POST", url, param)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}

func randSleep() {
	s := RandInt64(3000, 5000)
	time.Sleep(time.Millisecond * time.Duration(s))
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func getConfig() *Config {
	dir, _ := os.Getwd()
	content, _ := ioutil.ReadFile(dir + "/main/config.json")
	r := &Config{}
	json.Unmarshal(content, r)
	return r
}
