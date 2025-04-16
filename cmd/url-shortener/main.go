package main

import (
	"fmt"
	"restApiService/internal/config"
)

func main() {
	// TODO : init config cleanenv,
	// TODO : init logger slog
	// TODO: init db POSTGRES
	// TODO : init router HZ KAKOY
	// TODO : run server

	cfg := config.MustLoad()
	fmt.Println(cfg)

}
