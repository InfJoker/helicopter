package core

import "context"

type Node struct {
	Lseq   string `json:"lseq,omitempty"`
	Parent string `json:"ref"`
	Value  []byte `json:"content"`
}

type Storage interface {
	CreateNode(ctx context.Context, parent string, value []byte) (Node, error)

	GetSubTreeNodes(ctx context.Context, parent, fromLseq string) ([]Node, error)
}
