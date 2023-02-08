package httptool

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	ContentTypeJson    = "application/json"
	ContentTypeFrom    = "application/x-www-form-urlencoded"
	ContentTypeTextXml = "text/xml"
)

var (
	ErrNoParm      = errors.New("NoParm")
	ErrContentType = errors.New("ErrorContentType")
)

type Httptool struct {
	params map[string]interface{}
	header map[string]string
	body   string
	client *http.Client
}

// NewHttptool 新建对象
func NewHttptool() *Httptool {
	return &Httptool{
		params: make(map[string]interface{}),
		client: http.DefaultClient,
		header: make(map[string]string),
	}
}

func parseUrl(url string, hs ...struct{}) string {
	if !strings.Contains(url, "http") {
		pre := "http://"
		if len(hs) > 0 {
			pre = "https://"
		}
		url = pre + url
	}
	return url
}

// DefaultGet 默认的get方法
func (h *Httptool) DefaultGet(url string, hs ...struct{}) (string, error) {
	url = parseUrl(url, hs...)
	r, err := http.Get(url)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	return string(b), nil
}

// DefaultPost 不带任何参数
func (h *Httptool) DefaultPost(url string, hs ...struct{}) (string, error) {
	url = parseUrl(url, hs...)
	r, err := http.Post(url, "", nil)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	return string(b), nil
}

// // Post 携带参数和格式
// func (h *Httptool) PostWithParm(url string, contentType string, hs ...struct{}) (string, error) {
// 	url = parseUrl(url, hs...)
// 	bodyReader := strings.NewReader(h.body)
// 	if contentType == "" {
// 		contentType = ContentTypeJson
// 	}
// 	r, err := http.Post(url, contentType, bodyReader)
// 	if err != nil {
// 		return "", err
// 	}
// 	b, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(b), nil
// }

func (h *Httptool) AddParm(key string, value interface{}) *Httptool {
	h.params[key] = value
	return h
}

func (h *Httptool) ClearParm() {
	h.params = nil
}
func (h *Httptool) AddHeader(key string, value string) *Httptool {
	h.header[key] = value
	return h
}

func (h *Httptool) ClearHeader() {
	h.header = nil
}
func (h *Httptool) InitClient(client *http.Client) *Httptool {
	h.client = client
	return h
}

// BuildBodyJson 构建json格式的字符串
func (h *Httptool) BuildBodyJson() (body string, err error) {
	b, err := json.Marshal(h.params)
	if err != nil {
		return "", err
	}
	body = string(b)
	h.body = body
	return
}

// BuildBodyForm 构建form格式的字符串
func (h *Httptool) BuildBodyForm() (body string, err error) {
	if h.params == nil {
		return "", ErrNoParm
	}
	for k, v := range h.params {
		body = body + "&" + fmt.Sprintf("%s=%v", k, v)
	}
	h.body = body[1:]
	return body[1:], nil
}

// PostWithHeader 默认json格式请求
func (h *Httptool) Post(url, contentType string, hs ...struct{}) (string, error) {
	url = parseUrl(url, hs...)
	body, err := h.BuildBodyJson()
	if contentType == ContentTypeFrom {
		body, err = h.BuildBodyForm()
	}
	if err != nil {
		return "", err
	}
	r := strings.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, url, r)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", ContentTypeJson)
	if contentType == ContentTypeFrom {
		req.Header.Set("Content-Type", ContentTypeFrom)
	}
	for k, v := range h.header {
		req.Header.Set(k, v)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(b), nil
}

// GetWithHeader hs 传了参数为https,不传为http  contentType 可不传 默认json
func (h *Httptool) Get(url, contentType string, hs ...struct{}) (string, error) {
	var (
		body string
		err  error
		r    = new(strings.Reader)
	)
	url = parseUrl(url, hs...)
	switch contentType {
	case "":
		body, err = h.BuildBodyForm()
		if err != nil {
			return "", err
		}
		url = url + "?" + body
	case ContentTypeFrom:
		return "", ErrContentType
	default:
		body, err = h.BuildBodyJson()
		if err != nil {
			return "", err
		}
		r = strings.NewReader(body)
	}
	req, err := http.NewRequest(http.MethodGet, url, r)
	if err != nil {
		return "", err
	}

	for k, v := range h.header {
		req.Header.Set(k, v)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(b), nil
}
