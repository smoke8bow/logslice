package processor

import (
	"sync"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
)

// Job represents a chunk of a log file to be processed by a worker.
type Job struct {
	FilePath string
	Offset   int64
	Limit    int64
}

// WorkerPool manages concurrent processing of log file chunks.
type WorkerPool struct {
	numWorkers int
	parser     *parser.Parser
	pattern    *filter.PatternFilter
	timeRange  *filter.TimeRange
	writer     *output.Writer
	format     output.Format
}

// NewWorkerPool creates a WorkerPool with the given concurrency and dependencies.
func NewWorkerPool(
	numWorkers int,
	p *parser.Parser,
	pf *filter.PatternFilter,
	tr *filter.TimeRange,
	w *output.Writer,
	fmt output.Format,
) *WorkerPool {
	if numWorkers < 1 {
		numWorkers = 1
	}
	return &WorkerPool{
		numWorkers: numWorkers,
		parser:     p,
		pattern:    pf,
		timeRange:  tr,
		writer:     w,
		format:     fmt,
	}
}

// Run distributes jobs across workers and waits for completion.
func (wp *WorkerPool) Run(jobs []Job) {
	ch := make(chan Job, len(jobs))
	for _, j := range jobs {
		ch <- j
	}
	close(ch)

	var wg sync.WaitGroup
	for i := 0; i < wp.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range ch {
				wp.processJob(job)
			}
		}()
	}
	wg.Wait()
}

// processJob reads lines from the given file chunk and applies filters.
func (wp *WorkerPool) processJob(job Job) {
	r, err := reader.NewLineReaderAt(job.FilePath, job.Offset)
	if err != nil {
		return
	}
	defer r.Close()

	pipeline := New(r, wp.writer, wp.parser, wp.pattern, wp.timeRange, wp.format)
	pipeline.Run()
}
