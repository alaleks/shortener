package service_test

import (
	"fmt"

	"github.com/alaleks/shortener/internal/app/service"
)

func ExampleIsURL() {
	fmt.Println(service.IsURL("https://ya.ru") == nil)
	fmt.Println(service.IsURL("htts://ya.ru") == nil)
	// Output:
	// true
	// false
}

func ExampleGenUID() {
	uid := service.GenUID(5)
	fmt.Println(len(uid))
	// Output:
	// 5
}
