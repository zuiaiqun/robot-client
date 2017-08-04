package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"share/pb"
	"sync"
	"time"

	"strings"

	"strconv"

	"fmt"

	"github.com/funny/link"
	"github.com/golang/protobuf/proto"
)

// Packet 头部
type Packet struct {
	headerLength  int32
	ClientVersion int32
	Cmdid         pb.Cmdid
	Seq           int32
	bodyLength    int32
	Body          proto.Message
	BodyRaw       []byte
	Header        map[string]string
}

// Protocol 数据协议
func Protocol(rw io.ReadWriter) (link.Codec, error) {
	return &codec{
		rw:      rw,
		bodyBuf: make([]byte, 1024),
		sendBuf: bytes.NewBuffer([]byte{}),
		reqPool: &sync.Pool{
			New: func() interface{} { return new(pb.Request) },
		},
		readDeadlineInterval: CommonHeartbeatInterval,
	}, nil
}

type codec struct {
	locker sync.Mutex

	rw                   io.ReadWriter
	bodyBuf              []byte
	sendBuf              *bytes.Buffer
	reqPool              *sync.Pool
	readDeadlineInterval int64
	heatbeatTimeoutTimes int32
}

func (c *codec) Receive() (interface{}, error) {
	p := new(Packet)

	conn := c.rw.(net.Conn)
	conn.SetReadDeadline(time.Now().Add(time.Duration(c.readDeadlineInterval) * time.Second))
	// 读取头部
	if err := binary.Read(c.rw, binary.BigEndian, &p.headerLength); err != nil {
		log.Println("read headerLength err:", p.headerLength, ":", err)
		// 如果是心跳超时
		if c.readDeadlineInterval != CommonHeartbeatInterval && strings.HasSuffix(err.Error(), "i/o timeout") {
			c.heatbeatTimeoutTimes++
			if c.heatbeatTimeoutTimes == 3 {
				return nil, errors.New("心跳连续3次超时")
			}
			return nil, nil
		}
		return nil, err
	}
	if c.heatbeatTimeoutTimes > 0 {
		c.heatbeatTimeoutTimes = 0
	}

	if p.headerLength != 20 {
		log.Println(p.headerLength)
		return nil, errors.New("错误的头部长度")
	}

	if err := binary.Read(c.rw, binary.BigEndian, &p.ClientVersion); err != nil {
		log.Println("read ClientVersion err:")
		return nil, err
	}
	if err := binary.Read(c.rw, binary.BigEndian, &p.Cmdid); err != nil {
		log.Println("read Cmdid err:")
		return nil, err
	}
	if err := binary.Read(c.rw, binary.BigEndian, &p.Seq); err != nil {
		log.Println("read Seq err:")
		return nil, err
	}
	if err := binary.Read(c.rw, binary.BigEndian, &p.bodyLength); err != nil {
		log.Println("read bodyLength err:")
		return nil, err
	}

	// 心跳特殊处理
	if p.Cmdid == pb.Cmdid_CmdHeartbeat {
		p.Body = &pb.Heartbeat{}
		return p, nil
	}

	if p.bodyLength > 0 {
		if p.bodyLength > 1024*1024 {
			return nil, fmt.Errorf("too long body length:%d", p.bodyLength)
		}
		if int32(cap(c.bodyBuf)) < p.bodyLength {
			c.bodyBuf = make([]byte, p.bodyLength, p.bodyLength+128)
		}
		buff := c.bodyBuf[:p.bodyLength]
		if _, err := io.ReadFull(c.rw, buff); err != nil {
			log.Println("read body err11111:", err)
			return nil, err
		}

		fn, ok := ProtocolHandlers[p.Cmdid]
		if !ok {
			return nil, errors.New("未知的协议:" + strconv.Itoa(int(p.Cmdid)))
		}
		body := fn()
		if err := proto.Unmarshal(buff, body); err != nil {
			log.Println("Unmarshal p.Body err22222:", err)
			return nil, err
		}
		p.Body = body

	}

	return p, nil
}

func (c *codec) Send(message interface{}) error {
	p, ok := message.(*Packet)
	if !ok {
		return errors.New("非Packet, 无法发送")
	}

	p.headerLength = 20
	p.bodyLength = int32(len(p.BodyRaw))

	c.sendBuf.Reset()
	if err := binary.Write(c.sendBuf, binary.BigEndian, p.headerLength); err != nil {
		return err
	}
	if err := binary.Write(c.sendBuf, binary.BigEndian, p.ClientVersion); err != nil {
		return err
	}
	if err := binary.Write(c.sendBuf, binary.BigEndian, p.Cmdid); err != nil {
		return err
	}
	if err := binary.Write(c.sendBuf, binary.BigEndian, p.Seq); err != nil {
		return err
	}
	if err := binary.Write(c.sendBuf, binary.BigEndian, p.bodyLength); err != nil {
		return err
	}
	if p.bodyLength > 0 {
		if err := binary.Write(c.sendBuf, binary.BigEndian, p.BodyRaw); err != nil {
			return err
		}
	}

	conn := c.rw.(net.Conn)
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	n, err := c.rw.Write(c.sendBuf.Bytes())
	if err != nil {
		return err
	}

	if n != c.sendBuf.Len() {
		return errors.New("没有成功写入数据")
	}

	return nil
}

func (c codec) Close() error {
	if closer, ok := c.rw.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
