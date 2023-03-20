package core

type Node struct {
	Lseq   int64
	Parent int64
	Value  []byte
}

type Storage interface {
	CreateNode(parent int64, value []byte) (Node, error)

	GetSubTreeNodes(parent, fromLseq int64) ([]Node, error)
}
