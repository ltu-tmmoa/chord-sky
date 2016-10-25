package chord

import (
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"sync"
)

// LocalNode represents a potential member of a Chord ring.
//
// All public methods, except those accessing immutable properties, are locked.
// This means that calling another public method from within a public method
// may cause a deadlock.
type LocalNode struct {
	ipAddr      net.IPAddr
	id          Hash
	fingers     []*Finger
	predecessor Node
	mutex       sync.RWMutex
}

// NewLocalNode creates a new local node from given address, which ought to be the application's public-facing IP
// address.
func NewLocalNode(ipAddr *net.IPAddr) *LocalNode {
	return newLocalNode(ipAddr, hash(ipAddr, HashBitsMax))
}

func newLocalNode(ipAddr *net.IPAddr, id *Hash) *LocalNode {
	node := new(LocalNode)
	node.ipAddr = *ipAddr
	node.id = *id

	fingers := make([]*Finger, id.bits)
	for i := range fingers {
		fingers[i] = newFinger(id, i+1)
	}
	node.fingers = fingers
	node.predecessor = nil
	return node
}

// ID returns node ID.
func (node *LocalNode) ID() ID {
	return &node.id
}

// IPAddr provides node network IP address.
func (node *LocalNode) IPAddr() *net.IPAddr {
	return &node.ipAddr
}

// Finger resolves Chord node at given finger table offset i.
//
// The result is only defined for i in [1,M], where M is the amount of bits set
// at node ring creation.
func (node *LocalNode) Finger(i int) *Finger {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	if 1 > i || i > node.id.bits {
		panic(fmt.Sprintf("%d not in [1,%d]", i, node.id.bits))
	}
	return node.finger(i)
}

func (node *LocalNode) finger(i int) *Finger {
	return node.fingers[i-1]
}

func unsafeFinger(node Node, i int) *Finger {
	switch n := node.(type) {
	case *LocalNode:
		return n.finger(i)
	default:
		return n.Finger(i)
	}
}

// Successor yields the next node in this node's ring.
func (node *LocalNode) Successor() (Node, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return node.finger(1).Node(), nil
}

func successor(node Node) (Node, error) {
	switch n := node.(type) {
	case *LocalNode:
		return n.finger(1).Node(), nil
	default:
		return node.Finger(1).Node(), nil
	}
}

// Predecessor yields the previous node in this node's ring.
func (node *LocalNode) Predecessor() (Node, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return node.predecessor, nil
}

func predecessor(node Node) (Node, error) {
	switch n := node.(type) {
	case *LocalNode:
		return n.predecessor, nil
	default:
		return n.Predecessor()
	}
}

// FindSuccessor asks this node to find successor of given ID.
//
// See Chord paper figure 4.
func (node *LocalNode) FindSuccessor(id ID) (Node, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return findSuccessor(node, id)
}

func findSuccessor(node Node, id ID) (Node, error) {
	node0, err := findPredecessor(node, id)
	if err != nil {
		return nil, err
	}
	return successor(node0)
}

// FindPredecessor asks node to find id's predecessor.
//
// See Chord paper figure 4.
func (node *LocalNode) FindPredecessor(id ID) (Node, error) {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	return findPredecessor(node, id)
}

// Asks node to find id's predecessor.
//
// See Chord paper figure 4.
func findPredecessor(n Node, id ID) (Node, error) {
	n0 := n
	for {
		succ, err := successor(n0)
		if err != nil {
			return nil, err
		}
		if idIntervalContainsEI(n0.ID(), succ.ID(), id) {
			return n0, nil
		}
		n0, err = closestPrecedingFinger(n0, id)
		if err != nil {
			return nil, err
		}
	}
}

// Returns closest finger preceding ID.
//
// See Chord paper figure 4.
func closestPrecedingFinger(n Node, id ID) (Node, error) {
	for i := n.ID().Bits(); i > 0; i-- {
		if f := unsafeFinger(n, i).Node(); idIntervalContainsEE(n.ID(), id, f.ID()) {
			return f, nil
		}
	}
	return n, nil
}

// SetSuccessor attempts to set this node's successor to given node.
func (node *LocalNode) SetSuccessor(succ Node) error {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.finger(1).SetNode(succ)
	return nil
}

func setSuccessor(node, succ Node) error {
	switch n := node.(type) {
	case *LocalNode:
		n.finger(1).SetNode(succ)
		return nil
	default:
		return n.SetSuccessor(succ)
	}
}

// SetPredecessor attempts to set this node's predecessor to given node.
func (node *LocalNode) SetPredecessor(pred Node) error {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	node.predecessor = pred
	return nil
}

func setPredecessor(node, pred Node) error {
	switch n := node.(type) {
	case *LocalNode:
		n.predecessor = pred
		return nil
	default:
		return n.SetPredecessor(pred)
	}
}

// Join makes this node join the ring of given other node.
//
// If given node is nil, this node will form its own ring.
func (node *LocalNode) Join(node0 Node) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	if node0 != nil {
		node.initFingerTable(node0)
		updateOthers(node)
		// TODO: Move keys in (predecessor,n] from successor
	} else {
		m := node.id.Bits()
		for i := 1; i <= m; i++ {
			node.finger(i).SetNode(node)
		}
		node.predecessor = node
	}
}

// Update all nodes whose finger tables should refer to node.
//
// See Chord paper figure 6.
func updateOthers(n Node) error {
	m := n.ID().Bits()
	for i := 2; i <= m; i++ {
		var id ID
		{
			subt := big.Int{}
			subt.SetInt64(2)
			subt.Exp(&subt, big.NewInt(int64(i-1)), nil)
			id = n.ID().Diff(newHash(subt, m))
		}
		pred, err := findPredecessor(n, id)
		if err != nil {
			return err
		}
		err = unsafeUpdateFingerTable(pred, n, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// If s is the i:th finger of node, update node's finger table with s.
//
// See Chord paper figure 6.
func unsafeUpdateFingerTable(n, s Node, i int) error {
	fing := unsafeFinger(n, i)
	if idIntervalContainsIE(fing.Start(), fing.Node().ID(), s.ID()) {
		fing.SetNode(s)
		pred, err := predecessor(n)
		if err != nil {
			return err
		}
		unsafeUpdateFingerTable(pred, s, i)
	}
	return nil
}

// Initializes finger table of local node; node0 is an arbitrary node already in
// the network.
//
// See Chord paper figure 6.
func (node *LocalNode) initFingerTable(node0 Node) error {
	var err error
	// Add this node to node0 node's ring.
	{
		var succ Node
		var pred Node

		succ, err = findSuccessor(node0, node.finger(1).Start())
		if err != nil {
			return err
		}
		pred, err = predecessor(succ)
		if err != nil {
			return err
		}

		setSuccessor(node, succ)
		setPredecessor(node, pred)

		if err = setSuccessor(pred, node); err != nil {
			return err
		}
		if err = setPredecessor(succ, node); err != nil {
			return err
		}
	}
	// Update this node's finger table.
	{
		m := node.id.Bits()
		for i := 1; i < m; i++ {
			this := node.finger(i)
			next := node.finger(i + 1)
			var nextNode Node
			if idIntervalContainsIE(node.ID(), this.Node().ID(), next.Start()) {
				nextNode = this.Node()
			} else {
				nextNode, err = findSuccessor(node0, next.Start())
			}
			if err != nil {
				return err
			}
			next.SetNode(nextNode)
		}
	}
	return nil
}

// Stabilize attempts to fix any ring issues arising from joining or leaving Chord ring nodes.
//
// Recommended to be called periodically in order to ensure node data integrity.
func (node *LocalNode) Stabilize() error {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	succ := node.finger(1).Node()
	x, err := predecessor(succ)
	if err != nil {
		return err
	}
	if idIntervalContainsEE(node.ID(), succ.ID(), x.ID()) {
		if err = setSuccessor(node, x); err != nil {
			return err
		}
	}
	succ = node.finger(1).Node()
	return notify(succ, node)
}

func notify(n, n0 Node) error {
	pred, err := predecessor(n)
	if err != nil {
		return nil
	}
	if pred == nil || idIntervalContainsEE(pred.ID(), n.ID(), n0.ID()) {
		return setPredecessor(n, n0)
	}
	return nil
}

// FixRandomFinger refreshes this node's finger table entries in relation to Chord ring changes.
//
// Recommended to be called periodically in order to ensure finger table integrity.
func (node *LocalNode) FixRandomFinger() error {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	i := rand.Int() % len(node.fingers)
	finger := node.fingers[i]
	fingerNode, err := findSuccessor(node, finger.Start())
	finger.SetNode(fingerNode)
	return err
}

// FixAllFingers refreshes all of this node's finger table entries.
func (node *LocalNode) FixAllFingers() error {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	for _, finger := range node.fingers {
		fingerNode, err := findSuccessor(node, finger.Start())
		if err != nil {
			return err
		}
		finger.SetNode(fingerNode)
	}
	return nil
}

// PrintRing outputs this node's ring to console.
func (node *LocalNode) PrintRing() {
	node.mutex.RLock()
	defer node.mutex.RUnlock()

	fmt.Printf("Node %v ring:", node.id.String())
	succ, err := successor(node)
	if err != nil {
		fmt.Println(err)
		return
	}
	for !node.ID().Eq(succ.ID()) {
		fmt.Printf(" => %v", succ)
		succ, err = successor(node)
		if err != nil {
			fmt.Printf(" (%v)", err)
			return
		}
	}
	fmt.Println()
}

// String produces canonical string representation of this node.
func (node *LocalNode) String() string {
	return node.id.String()
}
