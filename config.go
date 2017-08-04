package main

import "fmt"

const ROBOT_SIZE = 2
const (
	PHONE_NUM = 15000000000
	PASSWORD  = "123456"
)

const (
	//HTTP_DOMAIN = "http://127.0.0.1:18084"
	HTTP_DOMAIN = "http://192.168.1.206:8084"
	//TCP_DOMAIN  = "http://127.0.0.1:15324"

	TCP_DOMAIN = "192.168.1.206:5324"
)

type robotConf struct {
	Mobile   string
	Password string
	DeviceID string
}

var conf = map[int]robotConf{
	1: robotConf{Mobile: "15018409888", Password: "123456"},
}

func init() {
	conf = make(map[int]robotConf, ROBOT_SIZE)
	for i := 1; i <= ROBOT_SIZE; i++ {
		conf[i] = robotConf{
			Mobile:   fmt.Sprintf("%d", i+PHONE_NUM),
			Password: PASSWORD,
			DeviceID: "string",
		}
	}
}
