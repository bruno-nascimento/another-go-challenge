package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/bruno-nascimento/another-go-challenge/config"
	"github.com/bruno-nascimento/another-go-challenge/service"
	"github.com/gorilla/mux"
)

type testHandler struct {
	cfg *config.Config
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamName := vars["name"]

	t.cfg.APP.TeamsList = []string{teamName}
	teamService := service.Teams{Cfg: t.cfg, WantedTeamNames: map[string]struct{}{strings.ToLower(teamName): {}}}

	teams := teamService.FindTeams()
	p, err := json.Marshal(teams[0].Data.Team.Players)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(p)
}

func (t testHandler) name() {

}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("error loading app config: ", err)
	}
	playersService := service.NewPlayers(cfg)
	playersService.Report()

	//router := mux.NewRouter()
	//
	//router.Methods("GET").Path("/team/{name}").Handler(testHandler{cfg: cfg})
	//
	//server := http.Server{
	//	Addr:         cfg.HTTP.Port,
	//	Handler:      router,
	//	ReadTimeout:  cfg.HTTP.Timeout,
	//	WriteTimeout: cfg.HTTP.Timeout,
	//}
	//log.Fatal(server.ListenAndServe())
}
