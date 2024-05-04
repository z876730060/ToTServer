package pkg

import (
	"net/http"
	"testing"
	"time"
)

func TestServer_Start(t *testing.T) {
	t.Parallel()
	server := NewServer()
	go func() {
		err := server.Start()
		if err != nil {
			t.Fail()
		}
	}()
	<-time.After(2 * time.Second)
	server.Kill()
	server.Wait()
	t.Log("Server test successful")
}

func TestServer_InitLog(t *testing.T) {
	NewServer().InitLog("TestServer_InitLog")
	t.Log("LogEngine test successful")
}

func TestServer_Load(t *testing.T) {
	t.Parallel()
	server := NewServer()
	go func() {
		err := server.LoadWebService(&http.Server{
			Addr: "localhost:8080",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello World"))
			}),
		}).Start()
		if err != nil {
			t.Fail()
		}
	}()
	<-time.After(2 * time.Second)
	server.Kill()
	server.Wait()
	t.Log("WebService test successful")
}

func TestServer_Banner(t *testing.T) {
	err := NewServer().SetBanner("./testdata/banner.txt").printBanner()
	if err != nil {
		t.Fail()
	}
}

func TestServer_InitCron(t *testing.T) {
	NewServer().InitCron()
	t.Log("CronEngine test successful")
}

func TestServer_Cron(t *testing.T) {
	cron := NewServer().InitCron().Cron()
	if cron != nil {
		t.Fail()
	}
	t.Log("CronEngine test successful")
}
