package pkg

import (
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	c      *cron.Cron
	jobMap map[string]cron.EntryID
	count  atomic.Int32 // 任务数
	mutex  sync.Mutex
}

func NewCron() *Cron {
	c := &Cron{
		c:      cron.New(),
		jobMap: make(map[string]cron.EntryID),
		count:  atomic.Int32{},
		mutex:  sync.Mutex{},
	}
	c.c.Start()
	return c
}

func (c *Cron) AddJob(spec, key string, job cron.Job) error {
	if c == nil {
		return errors.New("cron not init")
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 判断key是否存在,存在证明已经运行
	_, ok := c.jobMap[key]
	if ok {
		return errors.New("job already running")
	}

	id, err := c.c.AddJob(spec, job)
	if err != nil {
		return err
	}

	c.jobMap[key] = id
	c.count.Add(1)
	slog.Info("cron add job successfully", "jobNum", c.GetTaskNum(), "jobKey", key)
	return nil
}

func (c *Cron) GetTaskNum() int {
	return int(c.count.Load())
}

func (c *Cron) RemoveJob(key string) error {
	if c == nil {
		return errors.New("cron not init")
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	id, ok := c.jobMap[key]
	if !ok {
		return errors.New("job not running")
	}
	c.c.Remove(id)
	delete(c.jobMap, key)
	c.count.Add(-1)
	slog.Info("cron remove job successfully", "jobNum", c.GetTaskNum(), "jobKey", key)
	return nil
}
