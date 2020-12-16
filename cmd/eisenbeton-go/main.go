package main

import (
	"eisenbeton-go/flatbuff/eisenbeton/wire/request"
	"eisenbeton-go/flatbuff/eisenbeton/wire/response"
	"eisenbeton-go/wire"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/ilyakaznacheev/cleanenv"
	nats "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func convertHttpToProtoStruct(req *http.Request) []byte {

	body, _ := ioutil.ReadAll(req.Body)
	content := string(body)

	er := wire.EisenRequest{
		Host:        req.Host,
		Path:        req.URL.Path,
		Method:      req.Method,
		ContentType: req.Header.Get("Content-Type"),
		Content:     content,
	}

	msg, err := proto.Marshal(&er)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return msg
}

func convertHttpToFlatbuffBytes(req *http.Request) []byte {

	builder := flatbuffers.NewBuilder(1024)

	uri := builder.CreateString(req.RequestURI)
	path := builder.CreateString(req.URL.Path)
	method := builder.CreateString(req.Method)
	contentType := builder.CreateString(req.Header.Get("Content-Type"))

	body := make([]byte, req.ContentLength)
	n, err := req.Body.Read(body)
	log.Println(n, err)
	request.EisenRequestStartContentVector(builder, int(req.ContentLength))

	for i := req.ContentLength - 1; i >= 0; i-- {
		builder.PrependByte(body[i])
	}
	content := builder.EndVector(int(req.ContentLength))

	request.EisenRequestStart(builder)
	request.EisenRequestAddUri(builder, uri)
	request.EisenRequestAddPath(builder, path)
	request.EisenRequestAddMethod(builder, method)
	request.EisenRequestAddContentType(builder, contentType)
	request.EisenRequestAddContent(builder, content)
	reqBuf := request.EisenRequestEnd(builder)
	builder.Finish(reqBuf)

	buf := builder.FinishedBytes()

	return buf
}

func sendToNatsPubOnly(nc *nats.Conn, w http.ResponseWriter, req *http.Request) {

	//msg := convertHttpToProtoStruct(req)
	msg := convertHttpToFlatbuffBytes(req)

	go func() {
		err := nc.Publish(req.URL.Path, msg)
		if err != nil {
			log.Print(err)
		}
	}()

	w.WriteHeader(200)
	//io.WriteString(w, "Hi nats")

}

func sendToNatsReqRep(nc *nats.Conn, timeout time.Duration, w http.ResponseWriter, req *http.Request) {

	msg := convertHttpToFlatbuffBytes(req)

	resp, err := nc.Request(req.URL.Path, msg, timeout)
	if err != nil {
		// TODO: Provide a way to customize a response. Maybe scripting? Maybe static response?
		w.WriteHeader(200)
		return
	}
	resp.
}

func makeHandler(nc *nats.Conn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		sendToNatsPubOnly(nc, w, req)
	}
}

func handleHttp(port string, nc *nats.Conn) {

	http.HandleFunc("/", makeHandler(nc))
	log.Println("Listening at :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type ConfigDatabase struct {
	HttpPort    string `edn:"http-port" env:"HTTP_PORT"`
	PubOnly     bool   `edn:"pub-only" env:"PUB_ONLY"`
	NatsUri     string `edn:"nats-uri" env:"NATS_URI"`
	NatsTimeout int    `edn:"nats-timeout" env:"NATS_TIMEOUT"`
}

func readConfig() *ConfigDatabase {
	cfg := &ConfigDatabase{}
	cleanenv.ReadConfig("config.edn", cfg)
	return cfg
}

func main() {
	cfg := readConfig()
	nc, err := nats.Connect(cfg.NatsUri)
	if err != nil {
		panic(err)
	}
	log.Println("Connected to nats.io @ " + cfg.NatsUri)
	handleHttp(cfg.HttpPort, nc)

}
