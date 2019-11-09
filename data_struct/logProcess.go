/**
  create by yy on 2019-08-13
*/

package data_struct

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogProcess struct {
	Rc    chan []byte
	Wc    chan *Message
	Read  Reader
	Write Writer
}

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

type SystemInfo struct {
	HandleLine   int     `json:"handleLine"`   // 总处理行数
	Tps          float64 `json:"tps"`          // 系统吞吐量
	ReadChanLen  int     `json:"readChanLen"`  // read channel 长度
	WriteChanLen int     `json:"writeChanLen"` // write channel 长度
	RunTime      string  `json:"runTime"`      // 运行总时间
	ErrNum       int     `json:"errNum"`       // 错误数
}

type Monitor struct {
	StartTime time.Time
	Data      SystemInfo
}

func (m *Monitor) Start(lp *LogProcess) {
	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		m.Data.RunTime = time.Now().Sub(m.StartTime).String()
		m.Data.ReadChanLen = len(lp.Rc)
		m.Data.WriteChanLen = len(lp.Wc)
		data, _ := json.MarshalIndent(m.Data, "", "\t")

		_, _ = io.WriteString(writer, string(data))
	})

	_ = http.ListenAndServe("127.0.0.1:8998", nil)
}

func (l *LogProcess) Process() {
	// 解析模块

	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	for v := range l.Rc {
		ret := r.FindStringSubmatch(string(v))

		if len(ret) != 14 {
			log.Println("FindStringSubmatch failed: ", string(v))
			continue
		}

		location, _ := time.LoadLocation("Asia/Shanghai")
		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], location)
		if err != nil {
			log.Println("ParseInLocation failed, error: ", err.Error(), ret[4])
		}
		message.TimeLocal = t

		message.BytesSent, _ = strconv.Atoi(ret[8])

		reqSli := strings.Split(ret[6], " ")

		if len(reqSli) != 3 {
			log.Println("string split failed: ", ret[6])
			continue
		}

		message.Method = reqSli[0]
		u, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("url Parse failed: ", err)
			continue
		}

		message.Path = u.Path

		message.Scheme = ret[5]
		message.Status = ret[7]

		message.UpstreamTime, _ = strconv.ParseFloat(ret[12], 64)
		message.RequestTime, _ = strconv.ParseFloat(ret[13], 64)

		l.Wc <- message
	}
	//for {
	//	data := <-l.Rc
	//	l.Wc <- strings.ToUpper(string(data))
	//}
}
