package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/bruno-nascimento/another-go-challenge/config"
	"github.com/bruno-nascimento/another-go-challenge/models"
	"github.com/jedib0t/go-pretty/v6/table"
)

// Players Service
type Players struct {
	cfg          *config.Config
	teamsService TeamService
	teams        []models.Team
	players      map[string]*models.PlayerReport
	playerList   []*models.PlayerReport
}

// NewPlayers returns an instance of the Players service
func NewPlayers(cfg *config.Config) Players {
	return Players{cfg: cfg, teamsService: NewTeams(cfg), players: make(map[string]*models.PlayerReport)}
}

func (p *Players) Report() string {
	players := p.getPlayers()
	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.AppendHeader(table.Row{"#", "ID", "FULL NAME", "AGE", "LIST OF TEAMS"})
	for idx, player := range players {
		tbl.AppendRow([]interface{}{idx + 1, player.ID, player.Name, player.Age, strings.Join(player.Teams, ", ")})
	}
	tbl.SetStyle(table.StyleColoredBright)
	return tbl.Render()
}

func (p *Players) getPlayers() []*models.PlayerReport {
	if p.loadFromCache() {
		return p.playerList
	}
	p.teams = p.teamsService.FindTeams()
	p.retrieveFromTeams()
	p.playerList = make([]*models.PlayerReport, 0)
	for _, v := range p.players {
		p.playerList = append(p.playerList, v)
	}
	sort.Slice(p.playerList, func(i, j int) bool {
		return p.playerList[i].ID < p.playerList[j].ID
	})
	p.saveIntoCache()
	return p.playerList
}

func (p *Players) loadFromCache() bool {
	if !p.cfg.APP.CacheEnabled {
		return false
	}
	cachePath := path.Join(p.cfg.APP.CachePath, "players.json")
	stat, err := os.Stat(cachePath)
	if err != nil {
		log.Printf("error checking cache file: %s", err)
		return false
	}
	if stat.ModTime().Add(p.cfg.APP.CacheTTL).Before(time.Now()) {
		return false
	}
	cacheContent, err := os.ReadFile(cachePath)
	if err != nil {
		log.Printf("error reading cache file: %s", err)
		return false
	}
	err = json.Unmarshal(cacheContent, &p.playerList)
	if err != nil {
		log.Printf("error reading cache file: %s", err)
		return false
	}
	return true
}

func (p *Players) saveIntoCache() {
	if !p.cfg.APP.CacheEnabled {
		return
	}
	err := os.MkdirAll(p.cfg.APP.CachePath, 0777)
	if err != nil {
		log.Printf("error creating cache dir: %s", err)
		return
	}
	cachePath := path.Join(p.cfg.APP.CachePath, "players.json")
	destination, err := os.Create(cachePath)
	defer func() {
		_ = destination.Close()
	}()
	if err != nil {
		log.Printf("error creating players cache file: %s", err)
		return
	}
	playersBytes, err := json.Marshal(p.playerList)
	if err != nil {
		log.Printf("error serializing players list: %s", err)
		return
	}
	_, err = fmt.Fprintf(destination, "%s", playersBytes)
	if err != nil {
		log.Printf("error writing to players cache file: %s", err)
	}
}

func (p *Players) retrieveFromTeams() {
	for _, t := range p.teams {
		for _, pl := range t.Data.Team.Players {
			if player, ok := p.players[pl.Id]; ok {
				player.Teams = append(player.Teams, t.Data.Team.Name)
			} else {
				pr, err := models.NewPlayerReport(pl, t)
				if err != nil {
					log.Fatal(err.Error())
				}
				p.players[pl.Id] = pr
			}
		}
	}
}
