package main

import "github.com/pgorczyca/todo-list/internal/app"

func main() {
	app := app.New()
	app.ServeHTTP()
}
