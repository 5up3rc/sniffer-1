package http
import (
	"sync"
	"bufio"
	"net/http"
	"io"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/binlaniua/sniffer/logger"
)

var logHRS2 = logger.Logger{"http.response.stream"}

type HttpResponseStream struct {
	wg          *sync.WaitGroup
	httpRequest *HttpRequestStream
	reader      tcpreader.ReaderStream
	response    *http.Response
	complete    bool
}

func NewHttpResponseStream(httpRequest *HttpRequestStream, wg *sync.WaitGroup) *HttpResponseStream {
	//
	stream := HttpResponseStream{}
	stream.reader = tcpreader.NewReaderStream()
	stream.wg = wg;
	stream.httpRequest = httpRequest

	//
	go stream.start()

	//
	return &stream
}

func (stream HttpResponseStream) start() {
	//
	logHRS2.Debug("等待请求获取完毕")
	stream.wg.Wait()
	logHRS2.Debug("请求获取完毕, 开始获取影响数据")

	//
	stream.wg.Add(1)

	buf := bufio.NewReader(&stream.reader)
	for {
		if request := stream.httpRequest.request; request != nil {
			res, err := http.ReadResponse(buf, request)
			if err == io.EOF {
				logHRS2.Debug("end response eof")
				stream.response = res
				stream.end()
			} else if err != nil {
				logHRS2.Debug("something error %v", err)
			} else {
				logHRS2.Debug("end response complete")
			}
		}
	}
}

func (stream HttpResponseStream) end() {
	logHRS.Debug("end response")
	stream.wg.Done()
}



