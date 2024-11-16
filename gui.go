package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"
)

type webApp struct {
	ToClient      chan string
	blockOn       bool
	ActiveClients int
}

func (s *webApp) SendToClient(text string) {
	if s.ActiveClients > 0 {
		go func() {
			s.ToClient <- text
		}()
	}
}

func (s *webApp) HandlePOST(w http.ResponseWriter, r *http.Request) { //client sent requests
	fmt.Println(r.Method, r.RequestURI, r.Body)
	if r.Method == "POST" {
		switch r.URL.Query().Get("q") {
		case "stop":
			os.Exit(0)
		}
	}
}

func (s *webApp) ServeHTTP(w http.ResponseWriter, r *http.Request) { //sys updates
	flusher, ok := w.(http.Flusher)
	ctx := r.Context()

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	for {
		select {
		case msg := <-s.ToClient:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-ctx.Done(): //terminates goroutine on client disconnect
			return
		}
	}
}

func (s *webApp) GorotuineReporting(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	ctx := r.Context()
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.ActiveClients++
	for {
		select {
		case <-ctx.Done(): //terminates goroutine on client disconnect
			s.ActiveClients--
			fmt.Println("Client disconnected.")
			return
		default:
			time.Sleep(1 * time.Second)
			fmt.Fprintf(w, "data: %d\n\n", runtime.NumGoroutine())
			flusher.Flush()
		}

	}
}

func FileHandling(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.RequestURI)
	if r.Method == "GET" {
		_, fileOk := os.Open(r.RequestURI)
		fmt.Println(r.RequestURI, fileOk)

		http.ServeFile(w, r, r.RequestURI)
	}
}

func HandleFile(filename string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.RequestURI)
		if r.Method == "GET" {
			_, fileOk := os.Open(filename)
			fmt.Println(r.RequestURI, fileOk, filename)

			http.ServeFile(w, r, filename)
		}
	}
}

func StartWAServer(App *webApp) {
	//GUI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.RequestURI)
		if r.Method == "GET" {
			App.SendToClient(fmt.Sprintf("b%t", App.blockOn))
			_, fileOk := os.Open("web-app/index.html")
			fmt.Println(r.RequestURI, fileOk, "web-app/index.html")

			http.ServeFile(w, r, "web-app/index.html")
		}
	})
	http.HandleFunc("/eventmanager.js", HandleFile("web-app/eventmanager.js"))
	http.HandleFunc("/chart.umd.js", HandleFile("web-app/chart.js/dist/chart.umd.js"))
	http.HandleFunc("/chart.umd.js.map", HandleFile("web-app/chart.js/dist/chart.umd.js.map"))
	//Communications
	http.Handle("/updates", App)
	http.Handle("/reqs", http.HandlerFunc(App.HandlePOST))
	http.Handle("/goroutines", http.HandlerFunc(App.GorotuineReporting))
	fmt.Println("Server starting on:", pwettyPwint("http://127.0.0.1:8080", textProperties{Bold: true, Color: "#346eeb"}))
	http.ListenAndServe("127.0.0.1:8080", nil)
}
