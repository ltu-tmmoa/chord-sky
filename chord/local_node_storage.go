package chord

func (node *localNode) downloadStorageOf(peer Node) error {
	return transferStorage(peer, node)
}

func (node *localNode) uploadStorageTo(peer Node) error {
	return transferStorage(node, peer)
}

func transferStorage(fromNode, toNode Node) error {
	if fromNode.ID().Eq(toNode.ID()) {
		return nil
	}
	return nil
}
