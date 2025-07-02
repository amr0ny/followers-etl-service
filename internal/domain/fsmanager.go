package domain

import "context"

type FSManager interface {
	ReadFileBatched(ctx context.Context) (<-chan [][]string, error)
}
