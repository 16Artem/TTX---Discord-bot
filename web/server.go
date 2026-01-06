// web/server.go
package web

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"TTX/bot"
)

type Server struct {
	Bot *bot.Bot
}

type PageData struct {
	Title       string
	Servers     string
	Commands    string
	Users       string
	CurrentYear int
}

func NewServer(botInstance *bot.Bot) *Server {
	return &Server{
		Bot: botInstance,
	}
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/status" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "ok")
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	servers, commands, users := s.Bot.GetStats()

	data := PageData{
		Title:       "Лучший Дискорд бот",
		Servers:     fmt.Sprintf("%d", servers),
		Commands:    fmt.Sprintf("%d", commands),
		Users:       fmt.Sprintf("%d", users),
		CurrentYear: time.Now().Year(),
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
	}
}

func (s *Server) docsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/docs" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/docs.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона документации", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "TTX-Bot Документация",
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
	}
}

func (s *Server) Start() error {
	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets")))
	http.Handle("/assets/", assetsHandler)

	http.HandleFunc("/", s.homeHandler)
	http.HandleFunc("/status", s.statusHandler)
	http.HandleFunc("/docs", s.docsHandler)

	certFile := "./certs/server.crt"
	keyFile := "./certs/server.key"

	fmt.Println("HTTPS сервер запущен на https://localhost:443")
	return http.ListenAndServeTLS(":443", certFile, keyFile, nil)
}
