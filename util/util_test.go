package util

import (
	"net/http"
	"testing"
	"time"
)

var (
	validHtml = []byte("<html><body><main></main></body></html>")
	host      = []byte("http://0.0.0.0")
	port      = []byte("5900")
	validUrl  = []byte("/valid")
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc(string(validUrl), func(w http.ResponseWriter, r *http.Request) {
		w.Write(validHtml)
	})
	url := append([]byte(":"), port...)
	server := &http.Server{
		Addr:           string(url),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go server.ListenAndServe()

	m.Run()
	server.Close()
}

func TestGetPage(t *testing.T) {
	client := &http.Client{}

	t.Run("Valid html", func(t *testing.T) {
		url := append(host, []byte(":")...)
		url = append(url, port...)
		url = append(url, validUrl...)

		headers := []Header{
			{
				Key:   "User-Agent",
				Value: UserAgents["win_ff"],
			},
		}
		reqOpt := &RequestOptions{
			Client:  client,
			Url:     string(url),
			Headers: headers,
		}

		html, err := GetPage(reqOpt)
		t.Log(html)
		if err != nil && html != string(validHtml) {
			t.Errorf("This should return valid html (%s)", err)
		}
	})
}
