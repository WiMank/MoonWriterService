package main

import "./app"

func main() {
	r := app.InitApp("Alarm Service")
	r.Start()
	//r.Stop()
}
