package Cassandra

import (
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

var (
	cassandraHost           = os.Getenv("VERLOOP_DSN")
	once sync.Once
)

func Instance() *Cassandra {
	var cql *Cassandra
	once.Do(func() {
		cql = NewCQL()
	})
	return cql
}

type Cassandra struct {
	cluster    *gocql.ClusterConfig
	Session    *gocql.Session
	FailedJobs chan *Job
	MaxRetry   int
}

type Job struct {
	CQL     *Cassandra
	Payload interface{}
	Method  string
	Retry   int
}

func (cql *Cassandra) NewJob(payload interface{}, method string) *Job {
	return &Job{Payload: payload, CQL: cql, Retry: 0, Method: method}
}

func (job *Job) Execute() bool {
	result := reflect.ValueOf(job.Payload).MethodByName(job.Method).Call([]reflect.Value{})
	err := result[0].Interface()

	if err != nil {
		job.Retry++
		job.CQL.FailedJobs <- job
		return false
	}
	return true
}

func (cql *Cassandra) Connect() bool {
	var err error
	connected := false
	maxRetry := 5
	retry := 0
	for connected != true {
		if cql.Session == nil || cql.Session.Closed() {
			cql.Session, err = cql.cluster.CreateSession()
			if err != nil {
				if retry == maxRetry {
					log.Panic("Cassandra connection failed")
				}
				log.Errorf("Couldn't Connect to Cassandra: %s\nRetrying in 5 seconds...attempt %d", err, retry+1)
				time.Sleep(5 * time.Second)
				retry++
			} else {
				connected = true
				log.Info("Cassandra session started...")
			}
		} else {
			connected = true
		}
	}
	return connected
}

func NewCQL() *Cassandra {
	hostString := strings.Split(cassandraHost, ",")
	cql := Cassandra{}
	cql.FailedJobs = make(chan *Job, 10)
	cql.MaxRetry = 2
	cql.cluster = gocql.NewCluster(hostString...)
	cql.cluster.ProtoVersion = 4
	cql.cluster.Keyspace = "stories"
	cql.cluster.Timeout = 2000 * time.Millisecond
	cql.cluster.ConnectTimeout = 2000 * time.Millisecond
	cql.cluster.Consistency = gocql.Quorum
	cql.cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		NumRetries: 3,
		Min:        500 * time.Microsecond,
		Max:        6 * time.Second}

	cql.Connect()
	return &cql
}

func (cql *Cassandra) ProcessFailedJobs() {
	for {
		select {
		case job := <-cql.FailedJobs:
			if cql.Connect() {
				if job.Retry <= cql.MaxRetry {
					log.Warnf("Executing Failed Job %s: Attempt: %d, Method: %s, Data: %+v",
						reflect.TypeOf(job.Payload).Name(), job.Retry, job.Method, job.Payload)
					job.Execute()
				} else {
					log.Errorf("Job permanently failed %s: Data: %v, Method: %s",
						reflect.TypeOf(job.Payload).Name(), job.Payload, job.Method)
					//todo do some shit with this shit
				}
			}
		}
	}
}
