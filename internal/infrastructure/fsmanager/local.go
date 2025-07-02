package fsmanager

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/amr0ny/followers-etl-service/internal/domain"
	"io"
	"os"
)

type BatchSize int
type LocalFSManager struct {
	logger    domain.Logger
	fileName  string
	batchSize BatchSize
}

func NewLocalFSManager(logger domain.Logger, fileName string, batchSize BatchSize) domain.FSManager {
	return &LocalFSManager{logger: logger, fileName: fileName, batchSize: batchSize}
}

func (manager *LocalFSManager) ReadFileBatched(ctx context.Context) (<-chan [][]string, error) {
	ch := make(chan [][]string)

	file, err := os.Open(manager.fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	go func() {
		defer close(ch)
		defer file.Close()

		reader := csv.NewReader(file)
		var batch [][]string
		for {
			record, err := reader.Read()
			if err != nil {
				if err != io.EOF {
					manager.logger.Error("read error: %w", err)
					return
				}
				break
			}
			select {
			case <-ctx.Done():
				return
			default:
				batch = append(batch, record)
				if BatchSize(len(batch)) == manager.batchSize {
					batchCopy := make([][]string, len(batch))
					copy(batchCopy, batch)
					ch <- batchCopy
					batch = [][]string{}
				}
			}
		}
		if len(batch) > 0 {
			ch <- batch
		}
	}()

	return ch, nil
}
