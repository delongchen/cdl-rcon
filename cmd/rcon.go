package main

import (
	"fmt"
	"github.com/delongchen/cdl-rcon/pkg/mc"
	"github.com/delongchen/cdl-rcon/pkg/rcon"
	"time"
)

const (
	ADDR string = "localhost:25575"
)

func main() {
	s := rcon.NewSession(ADDR)
	h := mc.CoverSession(s)
	quite := make(chan struct{})
	go h.Start(quite)

	for i := 0; i < 20; i++ {
		time.Sleep(time.Second)
		r := h.ExecCMD("time set 0")
		fmt.Println(r)
	}

	<-quite
}
