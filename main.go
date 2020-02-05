package main

import (
	"github.com/WiMank/AlarmService/app"
)

func main() {
	r := app.InitApp("Alarm Service")
	r.Start()
	//r.Stop()
}
