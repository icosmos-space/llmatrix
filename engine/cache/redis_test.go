package cache

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRedisCache(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go serveFakeRedis(c)
		}
	}()

	r := NewRedis(ln.Addr().String(), time.Minute, WithRedisPrefix("t:"))
	defer r.Close()
	r.Set("k", map[string]string{"v": "1"}, time.Minute)
	got, ok := r.Get("k")
	if !ok {
		t.Fatal("miss")
	}
	b := got.([]byte)
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil || m["v"] != "1" {
		t.Fatalf("%v %v", err, m)
	}
	r.Delete("k")
	if _, ok := r.Get("k"); ok {
		t.Fatal("expected deleted")
	}
}

func serveFakeRedis(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)
	store := map[string]string{}
	for {
		args, err := readRESPArray(r)
		if err != nil {
			return
		}
		if len(args) == 0 {
			continue
		}
		switch strings.ToUpper(args[0]) {
		case "SETEX":
			if len(args) >= 4 {
				store[args[1]] = args[3]
			}
			_, _ = c.Write([]byte("+OK\r\n"))
		case "GET":
			if len(args) >= 2 {
				if v, ok := store[args[1]]; ok {
					_, _ = c.Write([]byte("$" + strconv.Itoa(len(v)) + "\r\n" + v + "\r\n"))
					continue
				}
			}
			_, _ = c.Write([]byte("$-1\r\n"))
		case "DEL":
			if len(args) >= 2 {
				delete(store, args[1])
			}
			_, _ = c.Write([]byte(":1\r\n"))
		default:
			_, _ = c.Write([]byte("+OK\r\n"))
		}
	}
}

func readRESPArray(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "*") {
		return nil, err
	}
	n, _ := strconv.Atoi(line[1:])
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		hlen, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		hlen = strings.TrimSpace(hlen)
		if !strings.HasPrefix(hlen, "$") {
			return nil, err
		}
		size, _ := strconv.Atoi(hlen[1:])
		buf := make([]byte, size+2)
		if _, err := ioReadFull(r, buf); err != nil {
			return nil, err
		}
		out = append(out, string(buf[:size]))
	}
	return out, nil
}
