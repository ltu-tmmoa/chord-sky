package net

import (
	"errors"
	"net"
	"sync"
	"time"
)

const (
	// TCPSchedulerTimeout determines connection and transaction timeouts.
	TCPSchedulerTimeout = 20 * time.Second
)

var (
	// ErrTCPSchedulerClosed signifies that TCPScheduler is closed.
	ErrTCPSchedulerClosed = errors.New("TCPScheduler closed")
)

// TCPScheduler manages scheduling of TCP reads/writes, making sure that no
// read/write transactions overlap despite being issued from different threads.
type TCPScheduler struct {
	addr net.TCPAddr
	conn *net.TCPConn
	mutx sync.Mutex
	term bool
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
// The address will not be dialed util the scheduler is provided with a
// transaction.
func NewTCPSchedulerLazy(addr *net.TCPAddr) *TCPScheduler {
	return &TCPScheduler{
		addr: *addr,
	}
}

// Addr returns TCP address held by scheduler.
func (scheduler *TCPScheduler) Addr() *net.TCPAddr {
	return scheduler.Addr()
}

// Schedule schedules provided transaction for eventual execution.
//
// The transaction is limited to a duration specified by `TCPSchedulerTimeout`.
func (scheduler *TCPScheduler) Schedule(trans func(*net.TCPConn) error, errh func(error)) {
	go func() {
		scheduler.mutx.Lock()
		defer scheduler.mutx.Unlock()

		if scheduler.term {
			errh(ErrTCPSchedulerClosed)
			return
		}

		if err := scheduler.dialTCP(); err != nil {
			errh(err)
			return
		}

		conn := scheduler.conn

		if err := conn.SetDeadline(time.Now().Add(TCPSchedulerTimeout)); err != nil {
			errh(err)
			return
		}
		if err := trans(conn); err != nil {
			errh(err)
		}
	}()
}

// Connects to TCP scheduler address, unless a current connection exists.
func (scheduler *TCPScheduler) dialTCP() error {
	if scheduler.conn == nil {
		conn0, err := net.DialTimeout("tcp", scheduler.Addr().String(), TCPSchedulerTimeout)
		if err != nil {
			scheduler.conn = nil
			return err
		}
		scheduler.conn, _ = conn0.(*net.TCPConn)
	}
	return nil
}

// Close terminates TCP scheduler, making further use of if always fail.
func (scheduler *TCPScheduler) Close() error {
	scheduler.mutx.Lock()
	defer scheduler.mutx.Unlock()

	scheduler.term = true

	var err error
	if scheduler.conn != nil {
		err = scheduler.conn.Close()
		scheduler.conn = nil
	}
	return err
}
