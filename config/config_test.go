package config

import (
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	t.Cleanup(os.Clearenv)
	err := os.Setenv("APP_TEAMS_LIST", "Gremio FBPA")
	if err != nil {
		t.Error(err)
		return
	}
	config, err := New()
	if err != nil {
		t.Fatalf("error creating new config: %s", err.Error())
	}
	if !reflect.DeepEqual(config.APP.TeamsList, []string{"Gremio FBPA"}) {
		t.Fatalf("we were expecting the APP.TeamsList prop to had an 'Gremio FBPA' value but we got : %s", config.APP.TeamsList)
	}
}

func TestDefault(t *testing.T) {
	t.Cleanup(os.Clearenv)
	config, err := New()
	if err != nil {
		t.Fatalf("error creating new config: %s", err.Error())
	}
	got := strings.Join(config.APP.TeamsList, ",")
	expected := "Germany,England,France,Spain,Manchester United,Arsenal,Chelsea,Barcelona,Real Madrid,Bayern Munich"
	if got != expected {
		t.Fatalf("we were expecting the APP.TeamsList prop to had an '%s' value but we got : %s", expected, got)
	}
}

func TestLoadMock(t *testing.T) {
	t.Cleanup(os.Clearenv)
	mapper := map[string]string{
		"APP_TEAMS_LIST": "Brazil",
	}
	config, err := NewMock(mapper)
	if err != nil {
		t.Fatalf("error creating mocked config: %s", err.Error())
	}
	if !reflect.DeepEqual(config.APP.TeamsList, []string{"Brazil"}) {
		t.Fatalf("we were expecting the APP.TeamsList prop to had an 'Brazil' value but we got : %s", config.APP.TeamsList)
	}
}

func TestParallelRequestProp(t *testing.T) {
	t.Cleanup(os.Clearenv)
	config, err := New()
	if err != nil {
		t.Fatalf("error creating new config: %s", err.Error())
	}
	if config.API.ParallelRequests == 0 {
		t.Fatalf("we were expecting the API.ParallelRequests prop to had an '%d' value but we got : %d", runtime.NumCPU()*4, config.API.ParallelRequests)
	}
	if config.API.ParallelRequests != runtime.NumCPU()*4 {
		t.Fatalf("we were expecting the API.ParallelRequests prop to had an '%d' value but we got : %d", runtime.NumCPU()*4, config.API.ParallelRequests)
	}
}
