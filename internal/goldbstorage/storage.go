package goldbstorage

import (
	"context"
	"helicopter/internal/core"

	"github.com/lodthe/goldb/db"
	"go.uber.org/zap"
)

type storage struct {
	conn *db.Connection
}

func NewStorage(logger *zap.Logger, address string) (*storage, error) {
	conn, err := db.Open(
		db.WithLogger(logger),
		db.WithServerAddress(address),
	)

	if err != nil {
		return nil, err
	}
	return &storage{conn: conn}, nil
}

func (s *storage) CreateNode(ctx context.Context, parent string, value []byte) (core.Node, error) {
	triplet, err := s.conn.Put(ctx, parent, string(value))
	if err != nil {
		return core.Node{}, err
	}
	return core.Node{
		Lseq:   triplet.Version.String(),
		Parent: triplet.Key,
		Value:  []byte(triplet.Value),
	}, nil
}

func (s *storage) GetSubTreeNodes(ctx context.Context, parentStr, fromLseqStr string) ([]core.Node, error) {
	res := make([]core.Node, 0)
	err := s.preorder(ctx, parentStr, fromLseqStr, &res)
	if err != nil {
		return []core.Node{}, err
	}
	return res, nil
}

func (s *storage) preorder(ctx context.Context, parentStr, fromLseqStr string, res *[]core.Node) error {
	options := []db.IterOption{
		db.IterKeyEquals(parentStr),
		db.IterFromVersion(db.NewVersion(fromLseqStr)),
	}
	iterator, err := s.conn.GetIterator(ctx, options...)
	if err != nil {
		return err
	}

	for iterator.HasNext() {
		item, err := iterator.GetNext()
		if err != nil {
			return err
		}
		child := item.Version.String()
		*res = append(*res, core.Node{
			Lseq:   child,
			Parent: item.Key,
			Value:  []byte(item.Value),
		})
		err = s.preorder(ctx, child, fromLseqStr, res)
		if err != nil {
			return err
		}
	}
	return nil
}
