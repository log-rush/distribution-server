package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/log-rush/simple-server/domain"
	"github.com/log-rush/simple-server/pkg/lrp"
)

type logJob struct {
	log    domain.Log
	stream string
}

type logDistributionWorkerPool struct {
	maxWorkers       int
	subscriptionRepo *domain.SubscriptionsRepository
	stop             chan bool
	jobs             chan logJob
	results          chan error
}

var (
	encoder = lrp.NewEncoder()
)

func NewPool(maxWorkers int, subscriptionRepo *domain.SubscriptionsRepository) logDistributionWorkerPool {
	return logDistributionWorkerPool{
		maxWorkers:       maxWorkers,
		subscriptionRepo: subscriptionRepo,
		stop:             make(chan bool),
		jobs:             make(chan logJob, 64),
		results:          make(chan error),
	}
}

func (p logDistributionWorkerPool) Start() {
	go func() {
		for err := range p.results {
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	for i := 0; i < p.maxWorkers; i++ {
		go worker(p.jobs, p.results, p.stop, p.subscriptionRepo)
	}
}

func (p logDistributionWorkerPool) PostJob(log domain.Log, stream string) {
	p.jobs <- logJob{
		log:    log,
		stream: stream,
	}
}

func worker(jobs <-chan logJob, result chan<- error, stop <-chan bool, repo *domain.SubscriptionsRepository) {
	for {
		select {
		case job := <-jobs:
			fmt.Printf("new job\n")
			subscribers, err := (*repo).GetSubscribers(context.Background(), job.stream)
			if err != nil {
				result <- err
			}
			wg := sync.WaitGroup{}
			for _, client := range subscribers {
				wg.Add(1)
				go func(client domain.Client) {
					client.Send <- encoder.Encode(lrp.NewMesssage(lrp.OprLog, []byte(job.log.Message)))
					wg.Done()
				}(client)
			}
			wg.Wait()
			result <- nil
		case <-stop:
			return
		}
	}
}
