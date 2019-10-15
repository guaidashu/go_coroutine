/**
  create by yy on 2019-08-13
*/

package data_struct

import (
	"bufio"
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"io"
	"os"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}

type ReadFromFile struct {
	Path string //文件路径
}

func (r *ReadFromFile) Read(rc chan []byte) {
	// 读取模块
	// 打开文件
	//fmt.Println("执行到了这里")
	file, err := os.Open(r.Path)
	if err != nil {
		panic(fmt.Sprintf("open file failed, errors: %s", err.Error()))
	}
	_, _ = file.Seek(0, 2)
	// 从文件 末尾开始逐行读取文件内容
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("read file failed, errors: %s", err.Error()))
		}
		rc <- line[:len(line)-1]
	}
}

type WriteToFluxDB struct {
	InfluxDBDsn string // influx data source
}

func (w *WriteToFluxDB) Write(wc chan *Message) {
	// 写入模块
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	q := client.NewQuery("SELECT count(value) FROM cpu_load", "mydb", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}

	for {
		fmt.Println(<-wc)
	}
}
