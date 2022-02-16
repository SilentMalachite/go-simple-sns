package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"time"
)

// データはJSONとして保存

const logFile = "logs.json"

//  保存するデータ定義

type Log struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Body  string `json:"body"`
	CTime int64  `json:"ctime"`
}

// サーバーを起動

func main() {
	println("server start - http://localhost:8888")  // 　環境によって書き換えてください
	println("if server stop ctl+c or ctl+d")

	http.HandleFunc("/", showHandler)
	http.HandleFunc("/write", writeHandler)

	http.ListenAndServe(":8888", nil)
}

// ログ読み出しとHTML生成

func showHandler(w http.ResponseWriter, r *http.Request) {
	htmlLog := ""
	logs := loadLogs()
	for _, i := range logs {
		htmlLog += fmt.Sprintf(
			"<p>(%d) <span>%s</span>: %s --- %s</p>",
			i.ID,
			html.EscapeString(i.Name),
			html.EscapeString(i.Body),
			time.Unix(i.CTime, 0).Format("2006/1/2 15:04"))
	}
	// HTMLテンプレ定義

	htmlBody := "<html><head><style>" +
		"p { border: 1px solid silver; padding: 1em;} " +
		"span { background-color: #eef; } " +
		"</style></head><body><h1>ローカルで動く簡易SNS</h1>" +
		getForm() + htmlLog + "</body></html>"
	w.Write([]byte(htmlBody))
}

// 送信した文言を書き込み

func writeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var log Log
	log.Name = r.Form["name"][0]
	log.Body = r.Form["body"][0]
	if log.Name == "" {
		log.Name = "名無しさん"
	}
	logs := loadLogs()
	log.ID = len(logs) + 1
	log.CTime = time.Now().Unix()
	logs = append(logs, log)
	saveLogs(logs)
	http.Redirect(w, r, "/", 302)
}

// 書き込んだフォームを示す

func getForm() string {
	return "<div><form action='/write' method='POST'>" +
		"名前: <input type='text' name='name'><br>" +
		"本文: <input type='text' name='body' style='width:60em;'><br>" +
		"<input type='submit' value='書込'>" +
		"</form></div><hr>"
}

// ファイルからログを読み込む

func loadLogs() []Log {
	text, err := ioutil.ReadFile(logFile)
	if err != nil {
		return make([]Log, 0)
	}

	var logs []Log
	json.Unmarshal([]byte(text), &logs)
	return logs
}

// ログファイルの書き込み

func saveLogs(logs []Log) {
	bytes, _ := json.Marshal(logs)
	ioutil.WriteFile(logFile, bytes, 0644)
}
