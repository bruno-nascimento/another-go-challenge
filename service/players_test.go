package service

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/bruno-nascimento/another-go-challenge/config"
	mocks "github.com/bruno-nascimento/another-go-challenge/mocks/service"
	"github.com/bruno-nascimento/another-go-challenge/models"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/report.txt
var expectedReport string

func TestName(t *testing.T) {
	cacheDir := path.Join(os.TempDir(), fmt.Sprint(time.Now().UnixMilli()))
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "cache disabled",
			args: func() args {
				mock, err := config.NewMock(map[string]string{"CACHE_ENABLED": "false"})
				if err != nil {
					t.Fatal(err)
				}
				return args{cfg: mock}
			}(),
			want: expectedReport,
		},
		{
			name: "cache enabled - write to cache",
			args: func() args {
				mock, err := config.NewMock(map[string]string{"CACHE_ENABLED": "true", "CACHE_PATH": cacheDir})
				if err != nil {
					t.Fatal(err)
				}
				return args{cfg: mock}
			}(),
			want: expectedReport,
		},
		{
			name: "cache enabled - load from cache",
			args: func() args {
				mock, err := config.NewMock(map[string]string{"CACHE_ENABLED": "true", "CACHE_PATH": cacheDir})
				if err != nil {
					t.Fatal(err)
				}
				return args{cfg: mock}
			}(),
			want: expectedReport,
		},
	}
	teamsServiceMock := new(mocks.TeamService)
	teamsServiceMock.On("FindTeams").Return(loadTeams(t))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playersService := &Players{
				cfg:          tt.args.cfg,
				teamsService: teamsServiceMock,
				players:      make(map[string]*models.PlayerReport),
			}
			result := playersService.Report()
			assert.Equal(t, expectedReport, result, "something is wrong with the report")
		})

	}
}

func TestPlayers_Report(t *testing.T) {
	teamsServiceMock := new(mocks.TeamService)
	teamsServiceMock.On("FindTeams").Return(loadTeams(t))
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}
	playersService := &Players{
		cfg:          cfg,
		teamsService: teamsServiceMock,
		players:      make(map[string]*models.PlayerReport),
	}
	result := playersService.Report()
	assert.Equal(t, expectedReport, result, "something is wrong with the report")
}

func loadTeams(t *testing.T) []models.Team {
	var teams []models.Team
	dir, err := teamsJSONFiles.ReadDir("testdata/teams/filtered")
	if err != nil {
		t.Fatal(err)
		return nil
	}
	for _, file := range dir {
		team := models.Team{}
		content, err := teamsJSONFiles.ReadFile(path.Join("testdata/teams/filtered/", file.Name()))
		if err != nil {
			t.Fatal(err)
			return nil
		}
		err = json.Unmarshal(content, &team)
		if err != nil {
			t.Fatal(err)
			return nil
		}
		teams = append(teams, team)
	}
	return teams
}
