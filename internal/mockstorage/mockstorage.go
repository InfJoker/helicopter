package mockstorage

import (
	"sync"

	"helicopter/internal/core"
)

type storage struct {
	mutex     sync.RWMutex
	replicaId int64
	seq       int64
	tree      map[int64][]int64
	values    map[int64][]byte
}

func genLseq(seq, replicaId int64) int64 {
	return (seq << 24) + replicaId
}

func NewStorage(replicaId int64) (*storage, error) {
	return &storage{
		replicaId: replicaId,
		seq:       0,
		tree:      make(map[int64][]int64),
		values:    make(map[int64][]byte),
	}, nil
}

func (s *storage) CreateNode(parent int64, value []byte) (core.Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.seq += 1
	lseq := genLseq(s.seq, s.replicaId)
	s.tree[parent] = append(s.tree[parent], lseq)
	s.values[lseq] = value
	return core.Node{
		Lseq:   lseq,
		Value:  value,
		Parent: parent,
	}, nil
}

func (s *storage) GetSubTreeNodes(parent, fromLseq int64) ([]core.Node, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	res := make([]core.Node, 0)
	s.getSubTreeNodesRec(parent, fromLseq, &res)
	return res, nil
}

func (s *storage) getSubTreeNodesRec(
	parent, fromLseq int64, res *[]core.Node,
) {
	children, is_parent := s.tree[parent]
	if !is_parent {
		return
	}
	for _, child := range children {
		if child > fromLseq {
			*res = append(*res, core.Node{
				Lseq:   child,
				Value:  s.values[child],
				Parent: parent,
			})
		}
		s.getSubTreeNodesRec(child, fromLseq, res)
	}
}
