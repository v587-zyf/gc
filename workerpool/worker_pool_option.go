package workerpool

type WorkerPoolOption struct {
	maxCount int
}

type Option func(o *WorkerPoolOption)

func NewWorkerPoolOption() *WorkerPoolOption {
	return &WorkerPoolOption{}
}

func WithMaxCount(maxCount int) Option {
	return func(o *WorkerPoolOption) {
		o.maxCount = maxCount
	}
}
