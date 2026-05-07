package processor

import (
	"sync"

	"github.com/user/logslice/internal/config"
	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/reader"
)

// ParallelResult holds the outcome of a parallel processing run.
type ParallelResult struct {
	LinesRead    int64
	LinesWritten int64
	Errors       []error
}

// RunParallel processes a log file in parallel using multiple workers,
// splitting the file into chunks and merging results in order.
func RunParallel(cfg *config.Config, w *output.Writer) (*ParallelResult, error) {
	info, err := reader.StatFile(cfg.InputFile)
	if err != nil {
		return nil, err
	}

	offsets := reader.ChunkOffsets(info.Size, int64(cfg.Workers))

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		result  = &ParallelResult{}
		errList []error
	)

	for i, off := range offsets {
		wg.Add(1)
		go func(workerID int, offset int64) {
			defer wg.Done()

			r, err := reader.NewLineReaderAt(cfg.InputFile, offset)
			if err != nil {
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
				return
			}
			defer r.Close()

			var tf *filter.TimeRange
			if cfg.Start != "" || cfg.End != "" {
				tf, err = filter.ParseTimeRange(cfg.Start, cfg.End)
				if err != nil {
					mu.Lock()
					errList = append(errList, err)
					mu.Unlock()
					return
				}
			}

			pf, err := filter.NewPatternFilter(cfg.Include, cfg.Exclude)
			if err != nil {
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
				return
			}

			p := New(r, w, tf, pf)
			stats, pErr := p.Run()
			mu.Lock()
			if pErr != nil {
				errList = append(errList, pErr)
			}
			result.LinesRead += stats.LinesRead
			result.LinesWritten += stats.LinesWritten
			mu.Unlock()
		}(i, off)
	}

	wg.Wait()
	result.Errors = errList
	return result, nil
}
