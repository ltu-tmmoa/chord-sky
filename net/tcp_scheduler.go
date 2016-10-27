package net

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// TCPScheduler manages scheduling of TCP reads/writes, making sure that no
// read/write transactions overlap despite being issued from different threads.
type TCPScheduler struct {
	addr net.TCPAddr
	conn *net.TCPConn
	mutx sync.Mutex
}

// NewTCPScheduler creates new initialized TCP scheduler.
//
// The provided TCP connection must be valid, and could suitably have been
// acquired via a TCP listen operation.
func NewTCPScheduler(conn *net.TCPConn) *TCPScheduler {
	addr, _ := conn.RemoteAddr().(*net.TCPAddr)
	return &TCPScheduler{
		addr: *addr,
		conn: conn,
	}
}

// NewTCPSchedulerLazy creates new uninitialized TCP scheduler.
//
// The address will not be dialed util the scheduler is provided with a task.
func NewTCPSchedulerLazy(addr *net.TCPAddr) *TCPScheduler {
	return &TCPScheduler{
		addr: *addr,
	}
}

// Addr returns TCP address held by scheduler.
func (scheduler *TCPScheduler) Addr() *net.TCPAddr {
	return scheduler.Addr()
}

// Schedule schedules provided function for eventual execution.
func (scheduler *TCPScheduler) Schedule(f func(*net.TCPConn) error) <-chan error {
	ch := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- errors.New(fmt.Sprint(r))
			}
		}()

		scheduler.mutx.Lock()
		defer scheduler.mutx.Unlock()

		scheduler.dialTCP()
		ch <- f(scheduler.conn)
	}()
	return ch
}

// Connects to TCP scheduler address, unless a current connection exists.
func (scheduler *TCPScheduler) dialTCP() error {
	if scheduler.conn == nil {
		conn0, err := net.DialTimeout("tcp", scheduler.Addr().String(), 20*time.Second)
		if err != nil {
			scheduler.conn = nil
			return err
		}
		scheduler.conn, _ = conn0.(*net.TCPConn)
	}
	return nil
}
