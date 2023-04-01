package goldbstorage

import (
	"context"
	"encoding/base64"
	"fmt"
	"helicopter/internal/config"
	"helicopter/internal/core"
	"strings"

	"github.com/lodthe/goldb/db"
	"go.uber.org/zap"
)

type storage struct {
	conn *db.Connection
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func decodeBase64(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}

func NewStorage(logger *zap.Logger, cfg config.Config) (*storage, error) {
	address := fmt.Sprintf("%s:%d", cfg.LseqDb.Host, cfg.LseqDb.Port)
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
	triplet, err := s.conn.Put(ctx, parent, encodeBase64(value))
	if err != nil {
		return core.Node{}, err
	}
	gotValue, err := decodeBase64(triplet.Value)
	if err != nil {
		return core.Node{}, fmt.Errorf("Error decoding database reply, maybe got data corruption: %w", err)
	}
	return core.Node{
		Lseq:   triplet.Version.String(),
		Parent: triplet.Key,
		Value:  gotValue,
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
		// db.IterFromVersion(db.NewVersion(fromLseqStr)), If uncommented we are loosing children of old parents
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
		gotValue, err := decodeBase64(item.Value)
		if err != nil {
			return fmt.Errorf("Error decoding database reply, maybe got data corruption: %w", err)
		}
		if fromLseqStr == "0" || strings.Compare(child, fromLseqStr) > 0 {
			*res = append(*res, core.Node{
				Lseq:   child,
				Parent: item.Key,
				Value:  gotValue,
			})
		}
		err = s.preorder(ctx, child, fromLseqStr, res)
		if err != nil {
			return err
		}
	}
	return nil
}
