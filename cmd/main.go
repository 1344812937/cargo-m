package main

import "cargo-m/internal/until"

func main() {
	stop := make(chan struct{})
	app := InitializeApp()
	app.Start()
	defer func() {
		app.Close()
		until.Log.Infof("Application stopped")
	}()
	<-stop
}
