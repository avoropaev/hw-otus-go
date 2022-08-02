package main

import (
	"context"

	_ "github.com/joho/godotenv/autoload"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/cmd"
)

func main() {
	cmd.Execute(context.Background())
}
