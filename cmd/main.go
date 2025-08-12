package main

func main() {
	app := InitializeApp()
	err := app.Run(":9090")
	if err != nil {
		println("运行异常!", err.Error())
		return
	}
}
