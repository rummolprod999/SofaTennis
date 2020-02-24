package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

type Tennis struct {
	homeTeam       string
	homeScore      []byte
	homeScoreMap   map[string]int
	awayTeam       string
	awayScore      []byte
	awayScoreMap   map[string]int
	statusType     string
	id             int64
	changeDate     string
	tournamentName string
	categoryName   string
	seasonName     string
}

func (m *Tennis) printMatch() {
	/*if m.statusType == "notstarted" || m.statusType == "finished" || m.statusType == "canceled" {
		return
	}*/
	fmt.Printf("Id game: %d\n", m.id)
	fmt.Printf("Status game: %s\n", m.statusType)
	fmt.Printf("Category name: %s\n", m.categoryName)
	fmt.Printf("Season name: %s\n", m.seasonName)
	fmt.Printf("Tournament name: %s\n", m.tournamentName)
	fmt.Printf("Date Change: %s\n", m.changeDate)
	fmt.Printf("Status game: %s\n", m.statusType)
	fmt.Printf("Home Team: %s\n", m.homeTeam)
	for k, v := range m.homeScoreMap {
		fmt.Printf("%s: %d\n", k, v)
	}
	fmt.Printf("Away Team: %s\n", m.awayTeam)
	for k, v := range m.awayScoreMap {
		fmt.Printf("%s: %d\n", k, v)
	}
	fmt.Printf("\n\n\n")
}

func (m *Tennis) sendMatch() {
	if m.statusType == "notstarted" || m.statusType == "finished" || m.statusType == "canceled" {
		return
	}
	for k, v := range m.awayScoreMap {
		if !strings.Contains(k, "period") {
			continue
		}
		if per, ok := m.homeScoreMap[k]; ok {
			if per == v {
				SendToTelegram(m)
			}
		}
	}
}
func (m *Tennis) CheckConditions() bool {
	_, okAwP1 := m.awayScoreMap["period1"]
	_, okAwP2 := m.awayScoreMap["period2"]
	_, okAHP := m.awayScoreMap["period2"]
	_, okHP2 := m.awayScoreMap["period2"]
	if !okAwP1 && !okAwP2 && !okAHP && !okHP2 {
		return false
	}
	return false
}
func SendToTelegram(m *Tennis) {
	if !CheckIfExist(fmt.Sprintf("%d", m.id)) {
		return
	}
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		Logging(err)
	}
	msg := tgbotapi.NewMessage(ChannelId, CreateMessage(m))
	msg.ParseMode = "html"
	_, err = bot.Send(msg)
	if err != nil {
		Logging(err)
	}
	Logging("send message")
}

func CreateMessage(m *Tennis) string {
	message := ""
	seasName, err := strconv.Unquote("\"" + m.seasonName + "\"")
	if err != nil {
		seasName = m.seasonName
	}
	seasName = strings.Replace(seasName, "\\", "", -1)
	tournamentName, err := strconv.Unquote("\"" + m.tournamentName + "\"")
	if err != nil {
		tournamentName = m.tournamentName
	}
	categoryName, err := strconv.Unquote("\"" + m.categoryName + "\"")
	if err != nil {
		categoryName = m.categoryName
	}
	message += fmt.Sprintf("<b>Category:</b> %s\n", categoryName)
	message += fmt.Sprintf("<b>Season:</b> %s\n", seasName)
	message += fmt.Sprintf("<b>Tournament:</b> %s\n", tournamentName)
	message += fmt.Sprintf("\n")
	message += fmt.Sprintf("<b>Date Change:</b> %s\n", m.changeDate)
	message += fmt.Sprintf("<b>Status Game:</b> %s\n", m.statusType)
	message += fmt.Sprintf("\n")
	homeTeam, err := strconv.Unquote("\"" + m.homeTeam + "\"")
	if err != nil {
		homeTeam = m.homeTeam
	}
	message += fmt.Sprintf("<b>Home Team:</b> %s\n", homeTeam)
	for k, v := range m.homeScoreMap {
		if !strings.Contains(k, "period") {
			continue
		}
		message += fmt.Sprintf("%s: %d\n", k, v)
	}
	message += fmt.Sprintf("\n")
	awayTeam, err := strconv.Unquote("\"" + m.awayTeam + "\"")
	if err != nil {
		awayTeam = m.awayTeam
	}
	message += fmt.Sprintf("<b>Away Team:</b> %s\n", awayTeam)
	for k, v := range m.awayScoreMap {
		if !strings.Contains(k, "period") {
			continue
		}
		message += fmt.Sprintf("%s: %d\n", k, v)
	}
	message += fmt.Sprintf("\n")
	return message
}

func CheckIfExist(id_game string) bool {
	db, err := DbConnection()
	if err != nil {
		Logging(err)
		return true
	}
	defer db.Close()
	rows, err := db.Query("SELECT id FROM sofa WHERE id_game=$1", id_game)
	if err != nil {
		Logging(err)
		return true
	}
	if rows.Next() {
		rows.Close()
		return false
	}
	rows.Close()
	_, err = db.Exec("INSERT INTO sofa (id, id_game) VALUES (NULL, $1)", id_game)
	if err != nil {
		Logging(err)
		return true
	}
	return true
}
