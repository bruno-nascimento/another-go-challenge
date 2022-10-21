package service

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/bruno-nascimento/another-go-challenge/config"
	"github.com/bruno-nascimento/another-go-challenge/models"
	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata
var teamsJSONFiles embed.FS

func TestTeams_FindTeams(t *testing.T) {
	cfg, _ := config.New()
	srv := teamsServer(t, cfg)
	cfg, _ = config.NewMock(map[string]string{"API_ENDPOINT": fmt.Sprintf("%s/acmeinc/{team_id}.json", srv.URL)})
	teamsService := NewTeams(cfg)
	teams := teamsService.FindTeams()
	assert.Equal(t, 10, len(teams), "10 teams were expected. got: %d", len(teams))
	for _, name := range cfg.APP.TeamsList {
		var found bool
		for _, team := range teams {
			if team.Data.Team.Name == name {
				found = true
			}
		}
		assert.True(t, found, "we were expecting find '%s' team but it wasn't found in the teams list", name)
	}
}

func teamsServer(t *testing.T, cfg *config.Config) *httptest.Server {
	m := make(map[string][]byte)
	for i := 0; i < len(cfg.APP.TeamsList); i++ {
		jsonBytes, err := teamsJSONFiles.ReadFile(fmt.Sprintf("testdata/teams/filtered/%s.json", strcase.ToSnake(cfg.APP.TeamsList[i])))
		if err != nil {
			t.Fatal(err)
		}
		var team models.Team
		err = json.Unmarshal(jsonBytes, &team)
		if err != nil {
			t.Fatal(err)
		}
		m[fmt.Sprint(team.Data.Team.Id)] = jsonBytes
	}
	re := *regexp.MustCompile(`^.*/(\d*).json$`)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := re.FindAllStringSubmatch(r.URL.Path, 1)
		if body, ok := m[res[0][1]]; ok {
			_, err := w.Write(body)
			if err != nil {
				t.Fatal(err)
			}
			return
		}
		body, err := json.Marshal(models.Team{})
		_, err = w.Write(body)
		if err != nil {
			t.Fatal(err)
		}
	}))
}
