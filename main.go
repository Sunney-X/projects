package main

import app "github.com/sunney-x/projects/cmd"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
