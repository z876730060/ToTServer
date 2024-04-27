package pkg

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestEngine_Start(t *testing.T) {
	t.Log(os.Getpid())
	ctx := context.Background()
	e := NewEngine(ctx, "测试子进程", &Options{
		Path: "cmd",
	})
	err := e.Start()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log("TestEngine_Start successfully")
}

func TestEngine_Stop(t *testing.T) {
	t.Log(os.Getpid())
	ctx := context.Background()
	e := NewEngine(ctx, "测试子进程", &Options{
		Path: "test.exe",
		Dir:  "testdata",
	})
	err := e.Start()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	time.Sleep(time.Second * 2)
	err = e.Stop()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log("TestEngine_Stop successfully")
}

func TestEngine_Release(t *testing.T) {
	t.Log(os.Getpid())
	ctx := context.Background()
	e := NewEngine(ctx, "测试子进程", &Options{
		Path: "test.exe",
		Dir:  "testdata",
	})
	err := e.Start()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	go func() {
		_ = e.Wait()
	}()

	time.Sleep(time.Second * 2)
	err = e.Release()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log("TestEngine_Release successfully")
}
