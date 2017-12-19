package fastgo

import (
	"strconv"
)

var server *HttpServer;

func init() {
	Conf.Prepare("./config/config.ini")
	port, err := strconv.Atoi(Conf.Get("APPPORT"))
	if err != nil {
		panic(err)
	}
	server = InitServer("", port, 5)
}

func Run() {
	server.Run()
}

func Router(c interface{}) {
	server.AddController(c)
}