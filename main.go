package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// 都市のデータ構造
type City struct {
	Name        string
	Temperature float64
	Weather     string
}

// 共有する気象データ（並行処理安全のためMutexを使用）
var (
	weatherData []City
	dataMutex   sync.Mutex
)

// HTMLテンプレート（画面の見た目）
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>気象情報シミュレーター</title>
    <meta http-equiv="refresh" content="3">
    <style>
        body { font-family: Arial, sans-serif; background: #f0f2f5; text-align: center; padding: 20px; }
        .container { display: flex; justify-content: center; gap: 20px; margin-top: 20px; }
        .card { background: white; padding: 20px; border-radius: 10px; box-shadow: 0 4px 8px rgba(0,0,0,0.1); width: 150px; }
        .temp { font-size: 24px; font-weight: bold; color: #ff5722; }
    </style>
</head>
<body>
    <h1>リアルタイム気象情報シミュレーター (Cgo-Free)</h1>
    <p>※3秒ごとに自動でデータが更新されます</p>
    <div class="container">
        {{range .}}
        <div class="card">
            <h2>{{.Name}}</h2>
            <p>天気: {{.Weather}}</p>
            <p class="temp">{{printf "%.1f" .Temperature}} °C</p>
        </div>
        {{end}}
    </div>
</body>
</html>
`

func main() {
	// 初期データ
	weatherData = []City{
		{Name: "札幌", Temperature: 4.0, Weather: "雪 ❄️"},
		{Name: "東京", Temperature: 16.0, Weather: "晴れ ☀️"},
		{Name: "沖縄", Temperature: 23.0, Weather: "曇り ☁️"},
	}

	// バックグラウンドでデータを常に更新し続ける (Goroutine)
	go func() {
		weathers := []string{"晴れ ☀️", "曇り ☁️", "雨 ☔", "雪 ❄️"}
		for {
			time.Sleep(3 * time.Second)
			dataMutex.Lock()
			for i := range weatherData {
				// 気温をランダムに増減
				weatherData[i].Temperature += (rand.Float64() * 2) - 1
				// 天気をランダムに変更
				weatherData[i].Weather = weathers[rand.Intn(len(weathers))]
			}
			dataMutex.Unlock()
		}
	}()

	// Webサーバーの設定
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.New("webpage").Parse(htmlTemplate)
		dataMutex.Lock()
		t.Execute(w, weatherData) // HTMLにデータを流し込む
		dataMutex.Unlock()
	})

	fmt.Println("サーバーを起動しました: http://localhost:8080")
	fmt.Println("ブラウザで上記のURLを開いてください。")
	http.ListenAndServe(":8080", nil)
}