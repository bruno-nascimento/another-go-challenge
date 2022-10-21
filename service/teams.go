package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/bruno-nascimento/another-go-challenge/config"
	"github.com/bruno-nascimento/another-go-challenge/models"
)

type TeamService interface {
	FindTeams() []models.Team
}

// Teams Service
type Teams struct {
	Cfg              *config.Config
	WantedTeamNames  map[string]struct{}
	waitGroup        sync.WaitGroup
	teamIDChannel    chan uint32
	stop             chan struct{}
	teamFoundChannel chan models.Team
	teams            []models.Team
	client           http.Client
}

// NewTeams returns an instance of the Teams service
func NewTeams(cfg *config.Config) TeamService {
	return &Teams{Cfg: cfg, WantedTeamNames: prepareWantedTeamsMap(cfg), client: http.Client{
		Timeout: cfg.API.RequestTimeout,
	}}
}

// FindTeams configured at 'APP_TEAMS_LIST' environment variable
func (t *Teams) FindTeams() []models.Team {
	t.waitGroup.Add(len(t.Cfg.APP.TeamsList))
	t.setupChannels()
	go t.teamIDIterator()
	go t.parallelTeamRequests()
	go t.storeFilteredTeams()
	t.wait()
	return t.teams
}

func (t *Teams) parallelTeamRequests() {
	for i := 0; i < t.Cfg.API.ParallelRequests; i++ {
		go t.teamRequest()
	}
}

func (t *Teams) teamRequest() {
	for teamID := range t.teamIDChannel {
		func() {
			resp, err := t.client.Get(strings.Replace(t.Cfg.API.Endpoint, "{team_id}", fmt.Sprint(teamID), 1))
			defer resp.Body.Close()
			if err != nil {
				log.Printf(err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf(err.Error())
					return
				}
				var team *models.Team
				err = json.Unmarshal(bodyBytes, &team)
				if err != nil {
					log.Printf(err.Error())
					return
				}
				if _, ok := t.WantedTeamNames[strings.ToLower(team.Data.Team.Name)]; ok {
					t.teamFoundChannel <- *team
				}
			}
		}()
		// else if resp.StatusCode != 200? retry with exponential backoff?
		// retry on errors as well
	}
}

func (t *Teams) teamIDIterator() {
	var teamID uint32
	for {
		select {
		case <-t.stop:
			close(t.teamIDChannel)
			close(t.stop)
			return
		default:
			teamID++
			t.teamIDChannel <- teamID
		}
	}
}

func (t *Teams) setupChannels() {
	t.teamIDChannel = make(chan uint32, t.Cfg.API.ParallelRequests)
	t.stop = make(chan struct{})
	t.teamFoundChannel = make(chan models.Team)
}

func (t *Teams) storeFilteredTeams() {
	for team := range t.teamFoundChannel {
		t.teams = append(t.teams, team)
		t.waitGroup.Done()
	}
}

func (t *Teams) wait() {
	t.waitGroup.Wait()
	t.stop <- struct{}{}
	close(t.teamFoundChannel)
}

func prepareWantedTeamsMap(cfg *config.Config) map[string]struct{} {
	teams := make(map[string]struct{})
	for _, t := range cfg.APP.TeamsList {
		teams[strings.ToLower(t)] = struct{}{}
	}
	return teams
}
