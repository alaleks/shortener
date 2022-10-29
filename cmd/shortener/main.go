package main

import (
	"log"
	"net/http"

	"github.com/alaleks/shortener/internal/app/serv"
)

func main() {
	server, err := serv.New(":8080")
	// здесь ошибка возвращается если недоступен порт
	// поэтому нужна проверка перед стартом сервера
	if err != nil {
		log.Fatal(err)
	}

	err = serv.Run(server)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
