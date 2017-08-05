package main

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/funny/link"
	"github.com/golang/protobuf/proto"
	"share/pb"
	"time"
)

type loginRes struct {
	Token string `json:"token"`
	User  struct {
		UID    int64  `json:"id"`
		FakeID string `json:"fakeID"`
	} `json:"user"`
}

type Robot struct {
	UID     int64
	Token   string
	FakeID  string
	Session *link.Session
	*robotConf
}

func NewRobot(c *robotConf) *Robot {
	return &Robot{
		robotConf: c,
	}
}

func (r *Robot) String() {
	log.Debugf("phone %s uid %d token %s fakeID %s", r.Mobile, r.UID, r.Token, r.FakeID)

}

func (r *Robot) Register(idx int) {
	resp, err := doApiHttpRequest(&requestOptions{
		Method: "POST",
		URL:    fmt.Sprintf("%s/%s", GetHttpDomain(), "user/UserExtService/Register"),
		Payload: fmt.Sprintf(`{"mobile":"%s", "password": "%s", "name": "%s", "code": "%s"}`,
			r.Mobile, r.Password, r.Mobile, "string"),
	})
	if err != nil {
		log.Errorf("register err %s\n", err.Error())
		return
	}
	res := new(loginRes)
	if err = json.Unmarshal([]byte(resp.Body), &res); err != nil {
		log.Errorf("unmarshal register err %s\n", err.Error())
		return
	}
	r.fillRobotWithToken(res, idx)
}

func (r *Robot) Login(idx int) {
	resp, err := doApiHttpRequest(&requestOptions{
		Method:  "POST",
		URL:     fmt.Sprintf("%s/%s", GetHttpDomain(), "user/UserExtService/MobileLogin"),
		Payload: fmt.Sprintf(`{"mobile":"%s", "password": "%s"}`, r.Mobile, r.Password),
	})
	if err != nil {
		log.Errorf("login err %s\n", err.Error())
		r.Register(idx)
		return
	}
	res := new(loginRes)
	if err = json.Unmarshal([]byte(resp.Body), &res); err != nil {
		log.Errorf("unmarshal login err %s\n", err.Error())
		return
	}
	if res.Token == "" {
		r.Register(idx)
		return
	}
	r.fillRobotWithToken(res, idx)
	r.tcpConnect()
}

func (r *Robot) fillRobotWithToken(res *loginRes, idx int) {
	r.UID = res.User.UID
	r.Token = res.Token
	r.FakeID = res.User.FakeID
	gRobots[idx] = *r
	r.String()
}

func (r *Robot) tcpConnect() {
	session, err := link.Dial("tcp", GetTcpDomain(), link.ProtocolFunc(Protocol), 0)
	if err != nil {
		log.Errorf("connector err %s", err.Error())
		return
	}
	r.Session = session
	r.sendPing()
	r.recvMsg()
	if err := r.SendMsg(pb.Cmdid_CmdLogin, &pb.LoginReq{Token: r.Token, UserID: r.UID}); err != nil {
		log.Errorf("login err %s", err.Error())
	} else {
		log.Debugf("login ok")
	}
}

func (r *Robot) SendMsg(cmdId pb.Cmdid, msg proto.Message) error {

	packet := &Packet{Cmdid: cmdId}
	if msg != nil {
		request := &pb.Request{ServiceName: "service"}
		b, _ := proto.Marshal(msg)
		request.Body = b
		data, _ := proto.Marshal(request)
		packet.BodyRaw = data
	}
	if err := r.Session.Send(packet); err != nil {
		log.Errorf("send msg err %s", err.Error())
		return err
	}
	return nil
}

func (r *Robot) sendPing() {
	ticker := time.NewTicker(time.Second * 45)
	go func() {
		for {
			select {
			case <-ticker.C:
				r.SendMsg(pb.Cmdid_CmdHeartbeat, nil)
				log.Debugf("%d send ping", r.Mobile)
			}
		}
	}()
}

func (r *Robot) recvMsg() {
	go func() {
		for {
			p, err := r.Session.Receive()
			if err != nil {
				log.Errorf("receive err %s", err.Error())
				continue
			}
			packet := p.(*Packet)
			log.Debugf("%d receive cmdId %d msg %s", r.Mobile, packet.Cmdid, packet.Body.String())
		}
	}()
}
