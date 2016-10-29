package chord

import (
	"github.com/ltu-tmmoa/chord-sky/log"
)

func (node *localNode) downloadStorageOf(peer Node) error {
	log.Logger.Println("Downloading storage from", peer, "...")
	return transferStorage(peer, node)
}

func (node *localNode) uploadStorageTo(peer Node) error {
	log.Logger.Println("Uploading storage to", peer, "...")
	return transferStorage(node, peer)
}

func transferStorage(fromNode, toNode Node) error {
	if fromNode.ID().Eq(toNode.ID()) {
		return nil
	}
	fromStorage := fromNode.Storage()
	toStorage := toNode.Storage()

	keys, err := fromStorage.GetKeyRange(toNode.ID(), fromNode.ID())
	if err != nil {
		return err
	}
	for _, key := range keys {
		value, err := fromStorage.Get(key)
		if err != nil {
			return err
		}
		toStorage.Set(key, value)
	}
	return nil
}
