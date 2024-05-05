package pkg

import (
	"context"
	"fmt"
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

func TestNewEngineManage(t *testing.T) {
	manage := NewEngineManage()
	fmt.Println(manage)
}

func TestEngineManage_Add(t *testing.T) {
	manage := NewEngineManage()
	err := manage.Add(context.Background(), "test", &Options{
		Path: "test.exe",
		Dir:  "testdata",
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log("TestEngineManage_Add successfully")
	fmt.Println(manage)
}

func TestEngineManage_GetCount(t *testing.T) {
	manage := NewEngineManage()
	if manage.GetCount() != 0 {
		t.Fail()
	}
	err := manage.Add(context.Background(), "test", &Options{
		Path: "test.exe",
		Dir:  "testdata",
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if manage.GetCount() != 1 {
		t.Fail()
	}
}

func TestEngineManage_Start(t *testing.T) {
	manage := NewEngineManage()
	if 0 != manage.GetRunCount() {
		t.Fail()
	}
	err := manage.Add(context.Background(), "test", &Options{
		Path: "test.exe",
		Dir:  "testdata",
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if err = manage.Start("test"); err != nil {
		t.Log(err)
		t.Fail()
	}
	if manage.GetRunCount() != 1 {
		t.Fail()
	}
	manage.Stop("test")
	t.Log("TestEngineManage_Start successfully")
}

func TestEngineManage_Stop(t *testing.T) {
	manage := NewEngineManage()
	if 0 != manage.GetRunCount() {
		t.Fail()
	}
	err := manage.Add(context.Background(), "test", &Options{
		Path: "test.exe",
		Dir:  "testdata",
	})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if err = manage.Stop("test"); err == nil {
		t.Fail()
	}
	t.Log("TestEngineManage_Stop error successfully")
}
