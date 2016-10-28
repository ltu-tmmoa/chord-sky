package chord

import "net"

// Node represents some Chord node, available either locally or remotely.
type Node interface {
	// ID returns node ID.
	ID() *ID

	// TCPAddr provides node network address.
	TCPAddr() *net.TCPAddr

	// fingerStart resolves start ID of finger table entry i.
	//
	// The result is only defined for i in [1,M], where M is the amount of bits
	// set at node ring creation.
	FingerStart(i int) *ID

	// fingerNode resolves Chord node at given finger table offset i.
	//
	// The result is only defined for i in [1,M], where M is the amount of bits
	// set at node ring creation.
	FingerNode(i int) <-chan NodeErr

	// SetFingerNode attempts to set this node's ith finger to given node.
	//
	// The operation is only valid for i in [1,M], where M is the amount of
	// bits set at node ring creation.
	SetFingerNode(i int, fing Node) <-chan error

	// Successor yields the next node in this node's ring.
	Successor() <-chan NodeErr

	// Predecessor yields the previous node in this node's ring.
	Predecessor() <-chan NodeErr

	// FindSuccessor asks this node to find successor of given ID.
	FindSuccessor(id *ID) <-chan NodeErr

	// FindPredecessor asks this node to find a predecessor of given ID.
	FindPredecessor(id *ID) <-chan NodeErr

	// SetSuccessor attempts to set this node's successor to given node.
	SetSuccessor(succ Node) <-chan error

	// SetPredecessor attempts to set this node's predecessor to given node.
	SetPredecessor(pred Node) <-chan error

	// String turns Node into its canonical string representation.
	String() string
}

// NodeErr represents the result of some node fetch operation.
type NodeErr struct {
	Node Node
	Err  error
}

// Unwrap returns contained node and error.
func (ne NodeErr) Unwrap() (Node, error) {
	return ne.Node, ne.Err
}

func newChanNodeErr(f func() (Node, error)) <-chan NodeErr {
	ch := make(chan NodeErr, 1)
	go func() {
		node, err := f()
		ch <- NodeErr{
			Node: node,
			Err:  err,
		}
	}()
	return ch
}

func newChanErr(f func() error) <-chan error {
	ch := make(chan error, 1)
	go func() {
		err := f()
		ch <- err
	}()
	return ch
}
