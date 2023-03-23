package mockstorage_test

import (
	"testing"

	"helicopter/internal/core"
	"helicopter/internal/mockstorage"

	"github.com/stretchr/testify/assert"
)

func TestCreateNode(t *testing.T) {
	replicaId := int64(1)
	s, err := mockstorage.NewStorage(replicaId)
	assert.NoError(t, err)

	node, err := s.CreateNode(0, []byte("root"))
	assert.NoError(t, err)

	assert.Equal(t, int64((1<<24)|replicaId), node.Lseq)
	assert.Equal(t, []byte("root"), node.Value)
	assert.Equal(t, int64(0), node.Parent)
}

func TestGetSubTreeNodes(t *testing.T) {
	s, err := mockstorage.NewStorage(1)
	assert.NoError(t, err)

	root, err := s.CreateNode(0, []byte("new-chat"))
	assert.NoError(t, err)

	node1, err := s.CreateNode(root.Lseq, []byte("message1"))
	assert.NoError(t, err)

	node2, err := s.CreateNode(root.Lseq, []byte("message2"))
	assert.NoError(t, err)

	node3, err := s.CreateNode(node1.Lseq, []byte("reply-to-message1"))
	assert.NoError(t, err)

	nodes, err := s.GetSubTreeNodes(root.Lseq, 0)
	assert.NoError(t, err)

	expectedNodes := []core.Node{
		{
			Lseq:   node1.Lseq,
			Value:  node1.Value,
			Parent: node1.Parent,
		},
		{
			Lseq:   node2.Lseq,
			Value:  node2.Value,
			Parent: node2.Parent,
		},
		{
			Lseq:   node3.Lseq,
			Value:  node3.Value,
			Parent: node3.Parent,
		},
	}

	assert.ElementsMatch(t, expectedNodes, nodes)
}
