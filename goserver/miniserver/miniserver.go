package miniserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"net"
	"strconv"

	"../utils"
)

const ContentText string = "text/plain"
const ContentHtml string = "text/html"
const ContentJson string = "application/json"
const ContentJavascript = "text/javascript"
const ContentCss = "text/css"
const ContentXIcon = "image/x-icon"
const ContentSVG = "image/svg+xml"

type MiniServer struct {
	port     int
	listener net.Listener
}

func NewServer(port int) *MiniServer {
	return &MiniServer{
		port: port,
	}
}

func (miniserver *MiniServer) StartListening(callback func(client *Client) *Response) {
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		request.ParseForm()
		var client *Client = newClient(request)

		res := callback(client)
		if res == nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte("Not found"))
		} else {
			var content []byte = []byte(res.body)
			response.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", res.contentType))
			response.Header().Set("Server", res.serverDescription)

			if len(res.file) > 0 {
				if _, err := os.Stat(res.file); err == nil {
					buf, err := ioutil.ReadFile(res.file)
					if err == nil {
						content = buf
					}
				}
			}

			response.WriteHeader(res.statusCode)
			response.Write(content)
		}
	})

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(miniserver.port))
	utils.Panic(err)
	miniserver.listener = listener
	http.Serve(listener, nil)
}

func (miniserver *MiniServer) StopListening() {
	if miniserver.listener != nil {
		miniserver.listener.Close()
	}
}
