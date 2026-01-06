// main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"TTX/bot"
	"TTX/web"
)

func loadTokenFromFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "DISCORD_BOT_TOKEN=") {
			token := strings.TrimPrefix(line, "DISCORD_BOT_TOKEN=")
			token = strings.Trim(token, `"`)
			return token
		}
	}
	return ""
}

func main() {
	token := loadTokenFromFile(".env")
	
	if token == "" {
		token = os.Getenv("DISCORD_BOT_TOKEN")
	}
	
	if token == "" {
		fmt.Println("Ошибка: DISCORD_BOT_TOKEN не установлен")
		fmt.Println("Создайте файл .env с содержанием:")
		fmt.Println("DISCORD_BOT_TOKEN=\"ваш_токен\"")
		return
	}

	fmt.Println("Токен успешно загружен!")

	discordBot, err := bot.NewBot(token)
	if err != nil {
		log.Fatal("Ошибка создания бота: ", err)
	}

	webServer := web.NewServer(discordBot)

	errChan := make(chan error, 2)

	go func() {
		fmt.Println("Запуск Discord бота...")
		if err := discordBot.Start(); err != nil {
			errChan <- fmt.Errorf("ошибка бота: %v", err)
		}
	}()

	go func() {
		fmt.Println("Запуск веб-сервера...")
		if err := webServer.Start(); err != nil {
			errChan <- fmt.Errorf("ошибка сервера: %v", err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case sig := <-sc:
		fmt.Printf("Получен сигнал %v, завершаем...\n", sig)
		time.Sleep(2 * time.Second)

	case err := <-errChan:
		fmt.Printf("Ошибка: %v\n", err)
	}

	fmt.Println("Программа завершена.")
}