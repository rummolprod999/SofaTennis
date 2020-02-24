package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"strconv"
	"time"
)

type SofaTennis struct {
	seasonName     string
	tournamentName string
	categoryName   string
}

func (t *SofaTennis) run() {
	curDate := time.Now().Format("2006-01-02")
	UrlSofa := fmt.Sprintf("https://www.sofascore.com/tennis//%s/json?_=%d", curDate, random(1000000000, 9999999999))
	UrlSofa = "https://www.sofascore.com/tennis/livescore/json"
	response := DownloadPage(UrlSofa)
	if response == "" {
		Logging("got empty string")
		return
	}
	t.workWithResponse(response)
}

func (t *SofaTennis) workWithResponse(s string) {
	defer SaveStack()
	_, err := jsonparser.ArrayEach([]byte(s), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			Logging(err, "callback workWithResponse")
			return
		}
		scoring, _, _, err := jsonparser.Get(value, "events")
		//fmt.Println(string(scoring))
		if err != nil {
			Logging(err, "scoring")
			return
		}
		seasonNameByte, _, _, err := jsonparser.Get(value, "season", "name")
		if err == nil {
			t.seasonName = string(seasonNameByte)
		} else {
			t.seasonName = ""
		}
		tournamentNameByte, _, _, err := jsonparser.Get(value, "tournament", "name")
		if err == nil {
			t.tournamentName = string(tournamentNameByte)
		} else {
			t.tournamentName = ""
		}
		categoryNameByte, _, _, err := jsonparser.Get(value, "category", "name")
		if err == nil {
			t.categoryName = string(categoryNameByte)
		} else {
			t.categoryName = ""
		}
		_, err = jsonparser.ArrayEach([]byte(scoring), t.TennisMatch)
		if err != nil {
			Logging(err, "TennisMatch")
			return
		}
	}, "sportItem", "tournaments")
	if err != nil {
		Logging(err, "sportItem tournaments")
		return
	}
}

func (t *SofaTennis) TennisMatch(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	homeTeam, _, _, err := jsonparser.Get(value, "homeTeam", "name")
	if err != nil {
		Logging(err, "homeTeam")
		return
	}
	homeScore, _, _, err := jsonparser.Get(value, "homeScore")
	if err != nil {
		Logging(err, "homeScore")
		return
	}
	awayTeam, _, _, err := jsonparser.Get(value, "awayTeam", "name")
	if err != nil {
		Logging(err)
		return
	}
	awayScore, _, _, err := jsonparser.Get(value, "awayScore")
	if err != nil {
		Logging(err, "awayScore")
		return
	}
	statusType, _, _, err := jsonparser.Get(value, "status", "type")
	if err != nil {
		Logging(err, "statusType")
		return
	}
	id, err := jsonparser.GetInt(value, "id")
	if err != nil {
		Logging(err, "id")
		return
	}
	changeDate, _, _, err := jsonparser.Get(value, "changes", "changeDate")
	if err != nil {
		//Logging(err, "changeDate")
		//return
	}
	homeScoreMap := make(map[string]int)
	err = jsonparser.ObjectEach(homeScore, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		valInt, _ := strconv.ParseInt(string(value), 10, 32)
		homeScoreMap[string(key)] = int(valInt)
		return nil
	})
	if err != nil {
		Logging(err, "homeScore map", fmt.Sprintf("%s", string(homeScore)))
		return
	}
	awayScoreMap := make(map[string]int)
	err = jsonparser.ObjectEach(awayScore, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		valInt, _ := strconv.ParseInt(string(value), 10, 32)
		awayScoreMap[string(key)] = int(valInt)
		return nil
	})
	if err != nil {
		Logging(err, "awayScore map", fmt.Sprintf("%s", string(awayScore)))
		return
	}
	volT := Tennis{homeTeam: string(homeTeam), homeScore: homeScore, awayTeam: string(awayTeam), awayScore: awayScore, statusType: string(statusType), id: id, changeDate: string(changeDate), homeScoreMap: homeScoreMap, awayScoreMap: awayScoreMap, seasonName: t.seasonName, tournamentName: t.tournamentName, categoryName: t.categoryName}
	volT.sendMatch()
}
