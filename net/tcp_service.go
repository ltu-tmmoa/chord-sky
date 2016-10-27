package net

import (
	"net"
	"time"
)

const (
	// TCPServiceTimeout determines connection timeouts.
	TCPServiceTimeout = 20 * time.Second
)

// TCPService handles listening for and accepting incoming TCP connections.
type TCPService struct {
	tcpAddr  net.TCPAddr
	listener *net.TCPListener
	chClose  chan int
}

// TCPAddr returns TCP address held by service.
func (service *TCPService) TCPAddr() *net.TCPAddr {
	return &service.tcpAddr
}

// Accept first blocks while starting to listen for incoming connections, and
// then starts accepting connections asynchronously.
//
// If listening fails, the method returns an error.
func (service *TCPService) Accept(callback func(*net.TCPConn, error)) error {
	var err error
	service.listener, err = net.ListenTCP("tcp", service.TCPAddr())
	if err != nil {
		return err
	}

	go func() {
		listener := service.listener
		for {
			callback(listener.AcceptTCP())
		}
	}()

	return nil
}

// Close terminates service, causing it to stop accepting additional
// connections.
func (service *TCPService) Close() error {
	return service.listener.Close()
}
