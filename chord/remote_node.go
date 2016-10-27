package chord

import (
	"fmt"
	"net"
	"time"

	"github.com/ltu-tmmoa/chord-sky/log"
)

const (
	schemeHTTP = "http"
)

// RemoteNode represents some Chord node available remotely.
type RemoteNode struct {
	tcpAddr  net.TCPAddr
	id       ID
	chAction chan func(*net.TCPConn) error
}

// NewRemoteNode creates a new remote node from given address.
//
// The node automatically connects to the provided address on creation.
func NewRemoteNode(tcpAddr *net.TCPAddr) *RemoteNode {
	return newRemoteNode(tcpAddr, identity(tcpAddr, HashBitsMax))
}

func newRemoteNode(tcpAddr *net.TCPAddr, id *ID) *RemoteNode {
	node := &RemoteNode{
		tcpAddr:  *tcpAddr,
		id:       *id,
		chAction: make(chan func(*net.TCPConn) error, 4),
	}
	var conn *net.TCPConn

	// Schedule exection of node actions.
	go func() {
		defer func() {
			if conn != nil {
				conn.Close()
			}
		}()

		var err error
		for action := range node.chAction {
			if err = action(conn); err != nil {
				close(node.chAction)
			}
		}
		if err != nil {
			log.Logger.Printf("Client connection failed: %s\n  Reason: %s\n", tcpAddr.String(), err.Error())
		} else {
			log.Logger.Printf("Client disconnected: %s\n", tcpAddr.String())
		}
		// TODO: Make sure remote node is removed from local node finger table.
		// onDisconnect(id) ?
	}()

	// Submit remote node connection action.
	node.chAction <- func(_ *net.TCPConn) error {
		conn0, err := net.DialTimeout("tcp", node.TCPAddr().String(), 20*time.Second)
		conn, _ := conn0.(*net.TCPConn)
		if err != nil {
			return err
		}
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Second * 25)
		return nil
	}

	return node
}

// ID returns node ID.
func (node *RemoteNode) ID() *ID {
	return &node.id
}

// TCPAddr provides node network address.
func (node *RemoteNode) TCPAddr() *net.TCPAddr {
	return &node.tcpAddr
}

// FingerStart resolves start ID of finger table entry i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *RemoteNode) FingerStart(i int) *ID {
	m := node.ID().Bits()
	verifyIndexOrPanic(m, i)
	return calcFingerStart(node.ID(), i)
}

// FingerNode resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *RemoteNode) FingerNode(i int) <-chan Node {
	return node.sendQueryForNode(messageTypeGetFingerNode, fmt.Sprint(i), "")
}

func (node *RemoteNode) sendQueryForNode(typ messageType, arg0, arg1 string) <-chan Node {
	ch := make(chan Node, 1)
	node.chAction <- func(conn *net.TCPConn) error {
		conn.SetDeadline(time.Now().Add(30 * time.Second))
		defer conn.SetDeadline(time.Time{})

		err := encodeMessage(conn, typ, arg0, arg1)
		if err != nil {
			ch <- nil
			return err
		}
		m, err := decodeMessage(conn, typ)
		if err != nil {
			ch <- nil
			return err
		}
		tcpAddr, err := net.ResolveTCPAddr("tcp", m.arg0)
		if err != nil {
			ch <- nil
			return err
		}
		ch <- NewRemoteNode(tcpAddr)
		return nil
	}
	return ch
}

// SetFingerNode attempts to set this node's ith finger to given node.
//
// The operation is only valid for i in [1,M], where M is the amount of
// bits set at node ring creation.
func (node *RemoteNode) SetFingerNode(i int, fing Node) {
	arg0 := fmt.Sprint(i)
	arg1 := fing.TCPAddr().String()
	node.commit(messageTypeSetFingerNode, arg0, arg1)
}

func (node *RemoteNode) commit(typ messageType, arg0, arg1 string) {
	node.chAction <- func(conn *net.TCPConn) error {
		conn.SetDeadline(time.Now().Add(30 * time.Second))
		defer conn.SetDeadline(time.Time{})

		return encodeMessage(conn, typ, arg0, arg1)
	}
}

// Successor yields the next node in this node's ring.
func (node *RemoteNode) Successor() <-chan Node {
	return node.sendQueryForNode(messageTypeGetSuccessor, "", "")
}

// Predecessor yields the previous node in this node's ring.
func (node *RemoteNode) Predecessor() <-chan Node {
	return node.sendQueryForNode(messageTypeGetPredecessor, "", "")
}

// FindSuccessor asks this node to find successor of given ID.
func (node *RemoteNode) FindSuccessor(id *ID) <-chan Node {
	return node.sendQueryForNode(messageTypeFindSuccessor, id.String(), "")
}

// FindPredecessor asks this node to find a predecessor of given ID.
func (node *RemoteNode) FindPredecessor(id *ID) <-chan Node {
	return node.sendQueryForNode(messageTypeFindPredecessor, id.String(), "")
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *RemoteNode) SetSuccessor(succ Node) {
	node.commit(messageTypeSetSuccessor, succ.TCPAddr().String(), "")
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *RemoteNode) SetPredecessor(pred Node) {
	node.commit(messageTypeSetPredecessor, pred.TCPAddr().String(), "")
}

// Disconnect removes this node from its ring, and terminates any live
// network connections.
//
// Using this node after calling this method yields undefined behavior.
func (node *RemoteNode) Disconnect() {
	close(node.chAction)
}

// String turns Node into its canonical string representation.
func (node *RemoteNode) String() string {
	return fmt.Sprintf("%s %s", node.id.String(), node.tcpAddr.String())
}
