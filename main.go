package main

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	_ "github.com/golang/protobuf/proto"
	"github.com/peterh/liner"
	"io/ioutil"
	_ "log"
	"net/http"
	"os"
	"strings"
	_ "sync"
)

type requestOptions struct {
	Header  map[string]string
	Query   map[string]string
	Payload string
	Method  string
	URL     string
}

type httpResponse struct {
	Status int
	Header http.Header
	Body   []byte
}

var gRobots = make(map[int]Robot, ROBOT_SIZE)
var gIp string
var gHttpPort, gTcpPort string

func doApiWithUserInput(r Robot, srv, payload string, withToken bool) (err error) {
	req := &requestOptions{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/%s", GetHttpDomain(), srv),
		Payload: payload,
	}
	if withToken {
		req.Header["X-Token"] = r.Token
	}
	resp, err := doApiHttpRequest(req)
	if err != nil {
		log.Errorf("do request %s err %s", req.URL, err.Error())
		return
	}
	log.Debugf("request %s resp %s", req.URL, string(resp.Body))
	return
}

func doApiHttpRequest(opts *requestOptions) (resp *httpResponse, err error) {
	var req *http.Request
	payload := strings.NewReader(opts.Payload)
	if req, err = http.NewRequest(opts.Method, opts.URL, payload); err != nil {
		return
	}

	for k, v := range opts.Header {
		req.Header.Add(k, v)
	}
	for k, v := range opts.Query {
		req.URL.Query().Add(k, v)
	}
	client := new(http.Client)
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return
	}
	resp = new(httpResponse)
	resp.Status = res.StatusCode
	resp.Header = res.Header

	if res.Body != nil {
		defer res.Body.Close()
		if resp.Body, err = ioutil.ReadAll(res.Body); err != nil {
			return
		}
	}
	return
}

func waitForUserInput(ch chan string, historyFile string) {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	line.SetCompleter(func(input string) (c []string) {
		//return line.getHistoryByPrefix(input)
		for _, n := range line.GetHistory() {
			if strings.HasPrefix(n, strings.ToLower(input)) {
				c = append(c, n)
			}
		}
		return
	})

	if f, err := os.Open(historyFile); err == nil {
		line.ReadHistory(f)
		f.Close()
	} else {
		fmt.Printf("read historyFile err %s", err.Error())
	}
	for {
		if cmd, err := line.Prompt("Input the cmd: "); err == nil {
			if _, _, _, err = parseMsg(cmd); err == nil {
				line.AppendHistory(cmd)
				ch <- cmd
			} else {
				fmt.Printf("parseMsg err cmd = \"%s\" err = \"%s\" \n", cmd, err.Error())
			}
		} else if err == liner.ErrPromptAborted { // CTRL+C (gracefully close all groutines)
			goto exit
		} else {
			log.Errorf("line err %s", err.Error())
		}
	}
exit:
	if f, err := os.Create(historyFile); err != nil {
		fmt.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
	line.Close()
}

// s=user/UserExtService/Regisetr&d={}&t=1
func parseMsg(msg string) (srv, payload string, withToken bool, err error) {
	slice := strings.Split(msg, "&")
	if len(slice) < 1 {
		err = errors.New("wrong params")
		return
	}
	for _, s := range slice {
		ss := strings.Split(s, "=")
		if len(ss) != 2 {
			err = errors.New("wrong params")
			return
		}
		switch ss[0] {
		case "s":
			srv = ss[1]
		case "d":
			payload = ss[1]
		case "t":
			if ss[1] == "1" {
				withToken = true
			}
		}
	}
	return
}

func loadAndLoginRobots() {
	for idx, r := range conf {
		robot := NewRobot(&r)
		robot.Login(idx)
	}
}

func handleUserInput(msg string) {
	srv, payload, withToken, err := parseMsg(msg)
	if err != nil {
		log.Errorf("parse msg %s err %s\n", msg, err.Error())
		return
	}
	for _, robot := range gRobots {
		err := doApiWithUserInput(robot, srv, payload, withToken)
		if err != nil {
			log.Errorf("handle msg %s err %s\n", msg, err.Error())
		}
	}
}

func main() {
	logName := flag.String("log", "log.xml", "file name")
	historyFile := flag.String("hf", ".history_file", "history file")
	gIp = *flag.String("ip", "127.0.0.1", "ip")
	gHttpPort = *flag.String("http_port", "18084", "http_port")
	gTcpPort = *flag.String("tcp_port", "15324", "tcp_port")
	flag.Parse()

	logger, err := log.LoggerFromConfigAsFile(*logName)
	if err != nil {
		fmt.Printf("Create logger error %s", err.Error())
		return
	}
	log.ReplaceLogger(logger)
	defer log.Flush()

	reqChan := make(chan string)
	go waitForUserInput(reqChan, *historyFile)
	loadAndLoginRobots()
	for {
		select {
		case msg := <-reqChan:
			handleUserInput(msg)
		}
	}
}
