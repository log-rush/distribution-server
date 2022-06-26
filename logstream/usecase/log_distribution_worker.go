package usecase

import (
	"context"
	"strconv"
	"sync"

	"github.com/log-rush/distribution-server/domain"
	"github.com/log-rush/distribution-server/pkg/commons"
	"github.com/log-rush/distribution-server/pkg/lrp"
)

type logJob struct {
	logs   []domain.Log
	stream string
}

type logDistributionWorkerPool struct {
	workers          []*logDistributionWorker
	maxWorkers       int
	subscriptionRepo *domain.SubscriptionsRepository
	logPlugins       *[]domain.LogPlugin
	jobs             chan logJob
	results          chan error
	l                *domain.Logger
}

type logDistributionWorker struct {
	id         int
	stop       chan bool
	jobs       <-chan logJob
	results    chan<- error
	l          *domain.Logger
	repo       *domain.SubscriptionsRepository
	logPlugins *[]domain.LogPlugin
}

var (
	encoder = lrp.NewEncoder()
)

func NewPool(maxWorkers int, subscriptionRepo *domain.SubscriptionsRepository, logPlugins *[]domain.LogPlugin, logger domain.Logger) logDistributionWorkerPool {
	return logDistributionWorkerPool{
		maxWorkers:       maxWorkers,
		subscriptionRepo: subscriptionRepo,
		jobs:             make(chan logJob, 64),
		results:          make(chan error),
		l:                &logger,
		workers:          []*logDistributionWorker{},
		logPlugins:       logPlugins,
	}
}

func (p logDistributionWorkerPool) Start() {
	(*p.l).Debugf("started worker pool (%d instances)", p.maxWorkers)
	go func() {
		defer commons.RecoverRoutine(p.l)
		for err := range p.results {
			if err != nil {
				(*p.l).Warnf("error in worker: %s", err.Error())
			}
		}
	}()

	for i := 0; i < p.maxWorkers; i++ {
		worker := newWorker(i, p.jobs, p.results, p.subscriptionRepo, p.logPlugins, p.l)
		p.workers = append(p.workers, worker)
		go worker.work()
		(*p.l).Debugf("[%d] worker started", worker.id)
	}
}

func (p logDistributionWorkerPool) PostJob(logs []domain.Log, stream string) {
	p.jobs <- logJob{
		logs:   logs,
		stream: stream,
	}
}

func (p logDistributionWorkerPool) Stop() {
	for _, worker := range p.workers {
		worker.stop <- true
		(*p.l).Debugf("[%d] worker stopped", worker.id)
	}
}

func newWorker(id int, jobs <-chan logJob, result chan<- error, repo *domain.SubscriptionsRepository, logPlugins *[]domain.LogPlugin, logger *domain.Logger) *logDistributionWorker {
	return &logDistributionWorker{
		id:         id,
		jobs:       jobs,
		stop:       make(chan bool),
		results:    result,
		repo:       repo,
		l:          logger,
		logPlugins: logPlugins,
	}
}

func (w *logDistributionWorker) work() {
	defer commons.RecoverRoutine(w.l)
	for {
		select {
		case job := <-w.jobs:
			(*w.l).Debugf("[%d] worker received job", w.id)
			wg := sync.WaitGroup{}
			subscribers, err := (*w.repo).GetSubscribers(context.Background(), job.stream)
			if err != nil {
				w.results <- err
			}

			(*w.l).Debugf("[%d] sending to %d subscribers", w.id, len(subscribers))
			for _, client := range subscribers {
				wg.Add(1)
				go func(client domain.Client) {
					defer commons.RecoverRoutine(w.l)
					for _, log := range job.logs {
						client.Send <- encoder.Encode(lrp.NewMesssage(lrp.OprLog, []byte(job.stream+","+strconv.Itoa(log.TimeStamp)+","+log.Message)))
					}
					wg.Done()
				}(client)
			}

			(*w.l).Debugf("[%d] sending to %d plugins", w.id, len(*w.logPlugins))
			wg.Add(1)
			go func() {
				for _, plugin := range *w.logPlugins {
					for _, log := range job.logs {
						plugin.HandleLog(log)
					}
				}
				wg.Done()
			}()

			wg.Wait()
		case <-w.stop:
			return
		}
	}
}
