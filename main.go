package main

import (
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"time"

	"offerlist/colly"
)

func spider(c *net.Conn, input string) {
	params := strings.Split(input, " ")
	if len(params) != 2 {
		Reply(c, "err")
		return
	}
	platform := params[0]
	shortId := params[1]
	switch platform {
	case "UA":
		data := colly.AmazonA(shortId)
		log.Println(data)
		Reply(c, data)
	default:
		Reply(c, "err")
	}
}

const SockAddr = "/tmp/offerlist.sock"

func Listen() (net.Listener, error) {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}
	oldMask := syscall.Umask(0000)
	ln, err := net.Listen("unix", SockAddr)
	syscall.Umask(oldMask)
	return ln, err
}

func Reply(c *net.Conn, txt string) {
	msg := txt + "\n"
	if c != nil {
		(*c).Write([]byte(msg))
	}
}

func handler(c net.Conn) {
	defer c.Close()
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			break
		}
		input := strings.TrimSuffix(string(buf[0:nr]), "\n")
		log.Println(input)
		switch input {
		case "ping":
			Reply(&c, "pong")
		default:
			if strings.HasPrefix(input, "colly") {
				spider(&c, strings.TrimPrefix(input, "colly "))
			}
		}
	}
}

func main() {
	ln, err := Listen()
	if err != nil {
		log.Printf("Listen error: %v", err)
	}

	for {
		c, err := ln.Accept()
		if err == nil {
			go handler(c)
		} else {
			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				break
			}
			time.Sleep(time.Second)
		}
	}
}
