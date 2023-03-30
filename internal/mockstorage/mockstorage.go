package mockstorage

import (
	"context"
	"strconv"
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

func toString(key int64) string {
	return strconv.FormatInt(key, 10)
}

func fromString(key string) (int64, error) {
	return strconv.ParseInt(key, 10, 64)
}

func NewStorage(replicaId int64) (*storage, error) {
	return &storage{
		replicaId: replicaId,
		seq:       0,
		tree:      make(map[int64][]int64),
		values:    make(map[int64][]byte),
	}, nil
}

func (s *storage) CreateNode(ctx context.Context, parentStr string, value []byte) (core.Node, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	parent, err := fromString(parentStr)
	if err != nil {
		return core.Node{}, err
	}
	s.seq += 1
	lseq := genLseq(s.seq, s.replicaId)
	s.tree[parent] = append(s.tree[parent], lseq)
	s.values[lseq] = value
	return core.Node{
		Lseq:   toString(lseq),
		Value:  value,
		Parent: parentStr,
	}, nil
}

func (s *storage) GetSubTreeNodes(ctx context.Context, parentStr, fromLseqStr string) ([]core.Node, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	parent, err := fromString(parentStr)
	if err != nil {
		return []core.Node{}, err
	}
	fromLseq, err := fromString(fromLseqStr)
	if err != nil {
		return []core.Node{}, err
	}
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
				Lseq:   toString(child),
				Value:  s.values[child],
				Parent: toString(parent),
			})
		}
		s.getSubTreeNodesRec(child, fromLseq, res)
	}
}
