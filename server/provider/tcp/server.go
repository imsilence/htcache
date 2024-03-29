package tcp

import (
	"bufio"
	"fmt"
	"htcache/cache"
	"htcache/server"
	"htcache/server/cluster"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	cluster.Node
	cache.Cache
}

func New(n cluster.Node, c cache.Cache) (server.Server, error) {
	return &Server{n, c}, nil
}

func (s *Server) Listen(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
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
	return nil
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

				if addr, ok := s.IsProcess(key); !ok {
					s.response(nil, fmt.Errorf("redirect: %s", addr), conn)
				} else {
					bytes, err := s.Get(key)
					// log.Printf("Get Key: %s, Value: %v, Error: %v", key, bytes, err)
					s.response(bytes, err, conn)
				}
			case 'S':
				key, value, err := s.readKeyAndValue(reader)
				if err != nil {
					if err != io.EOF {
						s.response(nil, err, conn)
					}
					break END
				}

				if addr, ok := s.IsProcess(key); !ok {
					s.response(nil, fmt.Errorf("redirect: %s", addr), conn)
				} else {
					err = s.Set(key, value)
					// log.Printf("Set Key: %s, Value: %v, Error: %v", key, value, err)
					s.response(nil, err, conn)
				}
			case 'D':
				key, err := s.readKey(reader)
				if err != nil {
					if err != io.EOF {
						s.response(nil, err, conn)
					}
					break END
				}

				if addr, ok := s.IsProcess(key); !ok {
					s.response(nil, fmt.Errorf("redirect: %s", addr), conn)
				} else {
					err = s.Del(key)
					log.Printf("Del Key: %s, Error: %v", key, err)
					s.response(nil, s.Del(key), conn)
				}
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
		_, err := conn.Write([]byte(fmt.Sprintf("-%d %s", len(e), e)))
		return err
	} else {
		_, err := conn.Write(append([]byte(fmt.Sprintf("%d ", len(result))), result...))
		return err
	}
}
