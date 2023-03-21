package core

type Node struct {
	Lseq   int64  `json:"lseq"`
	Parent int64  `json:"ref"`
	Value  []byte `json:"content"`
}

type Storage interface {
	CreateNode(parent int64, value []byte) (Node, error)

	GetSubTreeNodes(parent, fromLseq int64) ([]Node, error)
}
