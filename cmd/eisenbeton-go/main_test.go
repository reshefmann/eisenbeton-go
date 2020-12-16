package main

import (
	"eisenbeton-go/flatbuff/eisenbeton/wire/request"
	"eisenbeton-go/wire"
	"fmt"
	"io/ioutil"
	"log"

	//"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
)

func TestFlatbuff(t *testing.T) {
	builder := flatbuffers.NewBuilder(1024)

	uri := builder.CreateString("/hhh/jjj")
	path := builder.CreateString("/jjj")
	method := builder.CreateString("POST")
	contentType := builder.CreateString("application/json")

	request.EisenRequestStartContentVector(builder, 20)
	for i := 0; i <= 19; i++ {
		builder.PrependByte(byte(i))
	}
	content := builder.EndVector(20)

	request.EisenRequestStart(builder)
	request.EisenRequestAddUri(builder, uri)
	request.EisenRequestAddPath(builder, path)
	request.EisenRequestAddMethod(builder, method)
	request.EisenRequestAddContentType(builder, contentType)
	request.EisenRequestAddContent(builder, content)
	req := request.EisenRequestEnd(builder)
	builder.Finish(req)

	log.Println(req)
	buf := builder.FinishedBytes()

	log.Println(buf)
}

func TestWire(t *testing.T) {
	p := wire.EisenRequest{
		Host:        "hhhh",
		Path:        "",
		Method:      "",
		ContentType: "",
		Content:     "",
	}
	log.Print(p.Host)

	//t.Errorf("Error %s", p.Host)
}

func TestConvertHttpToPb2(t *testing.T) {
	req := httptest.NewRequest("POST", "http://example.com/test", ioutil.NopCloser(strings.NewReader("hello world")))
	w := httptest.NewRecorder()

	makeHandler(nil)(w, req)
	resp := w.Result()
	fmt.Println(req.URL.Path)
	fmt.Println(resp)
}

func TestReadConfig(t *testing.T) {
	cfg := readConfig()
	log.Println(cfg)
	if cfg.HttpPort != "8500" {
		t.Errorf("Expected 8500, got %s", cfg.HttpPort)
	}

}

//func TestConvertHttpToPb(t *testing.T) {
//req := http.Request{
//Method: "POST",
//Header: map[string][]string{
//"Content-Type": {"application/json"},
//},
//Body:       ioutil.NopCloser(strings.NewReader("hello world")),
//Host:       "aaaaaa",
//RequestURI: "jjj",
//}

//pb := convertHttpToProtoStruct(&req)

//log.Print(pb)
//if pb.Content != "hello world" {
//t.Errorf("Error: got %s", pb.Content)
//}
//}
