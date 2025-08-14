package main

func main() {
	app := InitializeApp()
	app.Start()
	defer app.Close()
}
