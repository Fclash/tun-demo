package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

var (
	tun *water.Interface
	closeHandler chan os.Signal
	msgChan chan []byte
)

func startRead() {
	buffer := [4096]byte{}
	for {
		fmt.Println("Start reading from", tun.Name())
		read, err := tun.Read(buffer[:])
		if err != nil {
			fmt.Println("warn: ", err)
		} else {
			fmt.Println("write ", read, "b")
		}
		var frame ethernet.Frame = buffer[:read]
		fmt.Println(frame.Source().String(), "->",frame.Destination().String())
		msgChan <- frame
	}
}

func startWrite() {
	for msg := range msgChan {
		write, err := tun.Write(msg)
		if err != nil {
			fmt.Println("warn in write: ", err)
		} else {
			fmt.Println("write to chan: ", write, "b")
		}
	}
}

func main()  {
	dev, err := water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: "tun-demo",
		},
	})
	if err != nil {
		fmt.Println("create tun device tun-demo failed", err)
		return
	}
	tun = dev
	defer tun.Close()
	go startRead()
	go startWrite()
	// ctrl+c handler
	closeHandler = make(chan os.Signal, 1)
	msgChan = make(chan []byte)
	signal.Notify(closeHandler, os.Interrupt)
	for range closeHandler {
		fmt.Println("recv close signal, close tun now.")
		break
	}
}