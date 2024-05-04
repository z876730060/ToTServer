package pkg

import "testing"

type Demo struct{}

func (Demo) Run() {

}

func TestCron_NewCron(t *testing.T) {
	cron := NewCron()
	if cron == nil {
		t.Fail()
	}
}

func TestCron_Add(t *testing.T) {
	cron := NewCron()
	err := cron.AddJob("0 0/1 * * *", "one", Demo{})
	if err != nil {
		t.Fatal(err)
	}
	if cron.GetTaskNum() != 1 {
		t.Fail()
	}
	t.Log("add job success")
}

func TestCron_RemoveJob(t *testing.T) {
	cron := NewCron()
	err := cron.AddJob("0 0/1 * * *", "one", Demo{})
	if err != nil {
		t.Fatal(err)
	}
	if cron.GetTaskNum() != 1 {
		t.Fail()
	}
	err = cron.RemoveJob("one")
	if err != nil {
		t.Fatal(err)
	}
	if cron.GetTaskNum() != 0 {
		t.Fail()
	}
	t.Log("remove job success")
}
