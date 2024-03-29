package mock

import (
	"adsb-api/internal/logger"
	"net"
)

type TcpStub interface {
	StartServer() (string, error)
	CloseConn() error
	CloseListener() error
	SetResponse(response []byte)
}

type StubImpl struct {
	ln       net.Listener
	conn     net.Conn
	response []byte
}

func InitStub(response []byte) *StubImpl {
	return &StubImpl{response: response}
}

func (stub *StubImpl) StartServer() (string, error) {
	logger.Info.Println("starting mock TCP server")

	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	stub.ln = ln

	go func() {
		defer func(ln net.Listener) {
			err := ln.Close()
			if err != nil {
				return
			}
		}(stub.ln)

		conn, err := stub.ln.Accept()
		if err != nil {
			return
		}

		stub.conn = conn

		defer func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				return
			}
		}(stub.conn)

		_, err = stub.conn.Write(stub.response)
		if err != nil {
			return
		}
	}()

	return stub.ln.Addr().String(), nil
}

func (stub *StubImpl) CloseConn() error {
	return stub.conn.Close()
}

func (stub *StubImpl) CloseListener() error {
	return stub.ln.Close()
}
