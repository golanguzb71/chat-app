package main

import (
	"chat-app/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.Modules).Run()
}
