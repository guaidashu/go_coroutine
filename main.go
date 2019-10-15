package main

import (
	"go_coroutine/data_struct"
	"time"
)

func main() {
	r := &data_struct.ReadFromFile{
		Path: "./files/access.log",
	}

	w := &data_struct.WriteToFluxDB{
		InfluxDBDsn: "username&password..",
	}

	lp := &data_struct.LogProcess{
		Rc:    make(chan []byte),
		Wc:    make(chan *data_struct.Message),
		Read:  r,
		Write: w,
	}

	go lp.Read.Read(lp.Rc)
	go lp.Process()
	go lp.Write.Write(lp.Wc)

	for {
		time.Sleep(500 * time.Second)
	}
}
