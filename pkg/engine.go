package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type Engine struct {
	cmd      exec.Cmd
	Pid      int
	name     string
	runState bool
}

type Options struct {
	Env  []string
	Path string
	Dir  string
	Args []string
}

func NewEngine(ctx context.Context, name string, options *Options) *Engine {
	cmd := exec.CommandContext(ctx, filepath.Join(options.Dir, options.Path), options.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = append(cmd.Env, options.Env...)
	cmd.Cancel = func() error {
		slog.Info("Engine cancel: name:" + name)
		return nil
	}
	e := new(Engine)
	e.cmd = *cmd
	e.name = name
	return e
}

// Start 子进程启动
func (e *Engine) Start() error {
	if err := e.cmd.Start(); err != nil {
		return err
	}
	e.runState = true
	e.Pid = e.cmd.Process.Pid
	slog.Info(fmt.Sprintf("Engine started, name: %s ,pid: %d", e.name, e.Pid))
	return nil
}

func (e *Engine) Wait() error {
	slog.Info(fmt.Sprintf("Engine waiting, name: %s ,pid: %d", e.name, e.Pid))
	wait, err := e.cmd.Process.Wait()
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Engine waited, name: %s ,pid: %d ,processState: %s", e.name, e.Pid, wait.String()))
	return nil
}

// Stop 终止进程
func (e *Engine) Stop() error {
	slog.Info(fmt.Sprintf("Engine stop, name: %s ,pid: %d", e.name, e.Pid))
	e.runState = false
	if err := e.cmd.Process.Kill(); err != nil {
		e.runState = true
		return err
	}
	return nil
}

type EngineManage struct {
	count    atomic.Int32
	mutex    sync.Mutex
	engines  map[string]*Engine
	runCount atomic.Int32
}

func NewEngineManage() *EngineManage {
	e := new(EngineManage)
	e.engines = make(map[string]*Engine)
	e.count = atomic.Int32{}
	e.mutex = sync.Mutex{}
	e.runCount = atomic.Int32{}
	return e
}

func (e *EngineManage) Add(ctx context.Context, name string, options *Options) error {
	_, err := e.add(ctx, name, options)
	if err != nil {
		return err
	}
	return nil
}

func (e *EngineManage) add(ctx context.Context, name string, options *Options) (*Engine, error) {
	e.mutex.Lock()
	if _, ok := e.engines[name]; ok {
		e.mutex.Unlock()
		return nil, fmt.Errorf("engine %s already exists", name)
	}
	e.count.Add(1)
	engine := NewEngine(ctx, name, options)
	e.engines[name] = engine
	e.mutex.Unlock()
	return engine, nil
}

func (e *EngineManage) Run(ctx context.Context, name string, options *Options) error {
	engine, err := e.add(ctx, name, options)
	if err != nil {
		return err
	}

	e.runCount.Add(1)
	if err := engine.Wait(); err != nil {
		e.runCount.Add(-1)
		return err
	}
	return nil
}

func (e *EngineManage) Start(name string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	engine, ok := e.engines[name]
	if !ok {
		return fmt.Errorf("engine %s not exists", name)
	}

	e.runCount.Add(1)
	if err := engine.Start(); err != nil {
		e.runCount.Add(-1)
		return err
	}
	return nil
}

func (e *EngineManage) Stop(name string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	engine, ok := e.engines[name]
	if !ok {
		return fmt.Errorf("engine %s not exists", name)
	}
	if !engine.runState {
		return fmt.Errorf("engine %s not running", name)
	}
	return engine.Stop()
}

func (e *EngineManage) GetCount() int32 {
	return e.count.Load()
}

func (e *EngineManage) GetRunCount() int32 {
	return e.runCount.Load()
}

func (e *EngineManage) Remove(name string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	engine, ok := e.engines[name]
	if !ok {
		return fmt.Errorf("engine %s not exists", name)
	}

	if engine.runState {
		return fmt.Errorf("engine %s already running", name)
	}

	delete(e.engines, name)
	e.count.Add(-1)
	return nil
}
