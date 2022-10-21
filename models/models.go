package models

import (
	"strconv"
)

type Team struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   struct {
		Team struct {
			Id          int    `json:"id"`
			OptaId      int    `json:"optaId"`
			Country     string `json:"country"`
			CountryName string `json:"countryName"`
			Name        string `json:"name"`
			LogoUrls    []struct {
				Size string `json:"size"`
				Url  string `json:"url"`
			} `json:"logoUrls"`
			IsNational      bool `json:"isNational"`
			HasOfficialPage bool `json:"hasOfficialPage"`
			Competitions    []struct {
				CompetitionId   int    `json:"competitionId"`
				CompetitionName string `json:"competitionName"`
			} `json:"competitions"`
			Players   []Player `json:"players"`
			Officials []struct {
				CountryName  string `json:"countryName"`
				Id           string `json:"id"`
				FirstName    string `json:"firstName"`
				LastName     string `json:"lastName"`
				Country      string `json:"country"`
				Position     string `json:"position"`
				ThumbnailSrc string `json:"thumbnailSrc"`
				Affiliation  struct {
					Name         string `json:"name"`
					ThumbnailSrc string `json:"thumbnailSrc"`
				} `json:"affiliation"`
			} `json:"officials"`
			Colors struct {
				ShirtColorHome string `json:"shirtColorHome"`
				ShirtColorAway string `json:"shirtColorAway"`
				CrestMainColor string `json:"crestMainColor"`
				MainColor      string `json:"mainColor"`
			} `json:"colors"`
		} `json:"team"`
	} `json:"data"`
	Message string `json:"message"`
}

type Player struct {
	Id           string `json:"id"`
	Country      string `json:"country"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Name         string `json:"name"`
	Position     string `json:"position"`
	Number       int    `json:"number"`
	BirthDate    string `json:"birthDate"`
	Age          string `json:"age"`
	Height       int    `json:"height"`
	Weight       int    `json:"weight"`
	ThumbnailSrc string `json:"thumbnailSrc"`
	Affiliation  struct {
		Name         string `json:"name"`
		ThumbnailSrc string `json:"thumbnailSrc"`
	} `json:"affiliation"`
}

type PlayerReport struct {
	ID    int
	Name  string
	Age   string
	Teams []string
}

func NewPlayerReport(p Player, t Team) (*PlayerReport, error) {
	id, err := strconv.Atoi(p.Id)
	if err != nil {
		return nil, err
	}
	return &PlayerReport{
		ID:    id,
		Name:  p.Name,
		Age:   p.Age,
		Teams: []string{t.Data.Team.Name},
	}, nil
}
