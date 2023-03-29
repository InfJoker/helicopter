package core

import "context"

type Node struct {
	Lseq   int64  `json:"lseq,omitempty"`
	Parent int64  `json:"ref"`
	Value  []byte `json:"content"`
}

type Storage interface {
	CreateNode(ctx context.Context, parent int64, value []byte) (Node, error)

	GetSubTreeNodes(ctx context.Context, parent, fromLseq int64) ([]Node, error)
}
