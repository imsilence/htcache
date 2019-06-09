package tcp

import (
	"bufio"
	"fmt"
	"htcache/cache"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	cache.Cache
}

func New(c cache.Cache) *Server {
	return &Server{c}
}

func (s *Server) Listen(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	log.Printf("Started TCP Server, Listen On: %s", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Client Connect is error: %s", err)
		}
		log.Printf("Client is Conncted: %s", conn.RemoteAddr())
		go s.Process(conn)
	}
}

func (s *Server) Process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
END:
	for {
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		op, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break END
			} else if operr, ok := err.(*net.OpError); ok && operr.Timeout() {
				// Timeout
			} else {
				log.Printf("Client Reader is error: %s", err)
				break END
			}
		} else {
			switch op {
			case 'G':
				key, err := s.readKey(reader)
				if err != nil {
					if err != io.EOF {
						s.response(nil, err, conn)
					}
					break END
				}
				bytes, err := s.Get(key)
				s.response(bytes, err, conn)
			case 'S':
				key, value, err := s.readKeyAndValue(reader)
				if err != nil {
					if err != io.EOF {
						s.response(nil, err, conn)
					}
					break END
				}
				err = s.Set(key, value)
				s.response(nil, err, conn)
			case 'D':
				key, err := s.readKey(reader)
				if err != nil {
					if err != io.EOF {
						s.response(nil, err, conn)
					}
					break END
				}
				s.response(nil, s.Del(key), conn)
			}
		}
	}
}

func (s *Server) readLen(reader *bufio.Reader) (int, error) {
	cxt, err := reader.ReadString(' ')
	if err != nil {
		return 0, err
	}
	len, err := strconv.Atoi(strings.TrimSpace(cxt))
	if err != nil {
		return 0, err
	}
	return len, nil
}

func (s *Server) readKey(reader *bufio.Reader) (string, error) {
	len, err := s.readLen(reader)
	if err != nil {
		return "", err
	}
	bytes := make([]byte, len)
	_, err = io.ReadFull(reader, bytes)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *Server) readKeyAndValue(reader *bufio.Reader) (string, []byte, error) {
	klen, err := s.readLen(reader)
	if err != nil {
		return "", nil, err
	}
	vlen, err := s.readLen(reader)
	if err != nil {
		return "", nil, err
	}
	kbytes := make([]byte, klen)
	_, err = io.ReadFull(reader, kbytes)
	if err != nil {
		return "", nil, err
	}
	vbytes := make([]byte, vlen)
	_, err = io.ReadFull(reader, vbytes)
	if err != nil {
		return "", nil, err
	}
	return string(kbytes), vbytes, nil
}

func (s *Server) response(result []byte, err error, conn net.Conn) error {
	if err != nil {
		e := err.Error()
		_, err := conn.Write([]byte(fmt.Sprintf("-%d%s", len(e), e)))
		return err
	} else {
		_, err := conn.Write(append([]byte(fmt.Sprintf("%d", len(result))), result...))
		return err
	}
}
