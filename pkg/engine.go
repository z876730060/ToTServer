package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type Engine struct {
	cmd  exec.Cmd
	Pid  int
	name string
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
	return e.cmd.Process.Kill()
}

// Release 释放关联资源
func (e *Engine) Release() error {
	slog.Info(fmt.Sprintf("Engine release, name: %s ,pid: %d", e.name, e.Pid))
	return e.cmd.Process.Release()
}
