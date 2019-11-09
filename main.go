package main

import (
	"flag"
	"go_coroutine/data_struct"
	"time"
)

func main() {

	var path, influxdbDsn string
	flag.StringVar(&path, "path", "./files/access.log", "read file path")
	flag.StringVar(&influxdbDsn, "influxdbDsn", "http://127.0.0.1:8999@yy@wyysdsa!@test@s", "influxdb data source")
	flag.Parse()

	r := &data_struct.ReadFromFile{
		Path: path,
	}

	w := &data_struct.WriteToFluxDB{
		InfluxDBDsn: influxdbDsn,
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

	m := data_struct.Monitor{
		StartTime: time.Now(),
		Data:      data_struct.SystemInfo{},
	}

	m.Start(lp)

	for {
		time.Sleep(500 * time.Second)
	}
}
