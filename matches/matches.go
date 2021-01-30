package matches

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/pished/esport-notify/text"
)

const scheduleUrl = "https://lolesports.com/schedule?leagues=lcs,lec,lck"

var valuedTeams = []string{"T1", "DK", "GEN", "HLE", "C9", "TSM", "100", "TL", "FNC", "G2"}

func GetNextMatch() {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var team1, team2, league, matchHour, ampm string
	var dates []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(scheduleUrl),
		chromedp.WaitVisible(".Event"),
		chromedp.Text(".single.future.event .teams .team.team1 .team-info .tricode", &team1, chromedp.ByQuery),
		chromedp.Text(".single.future.event .teams .team.team2 .team-info .tricode", &team2, chromedp.ByQuery),
		chromedp.Text(".single.future.event .league .name", &league, chromedp.ByQuery),
		chromedp.Text(".single.future.event .EventTime .time .hour", &matchHour, chromedp.ByQuery),
		chromedp.Text(".single.future.event .EventTime .time .ampm", &ampm, chromedp.ByQuery),
		chromedp.Nodes(".weekday", &dates, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get the game time relative to current time in minutes
	timeI, _ := strconv.Atoi(matchHour)
	gameTime := getMinutesUntil(timeI, ampm)
	//log.Println(gameTime)

	if isMatchToday(dates) {
		//log.Printf("%s plays against %s", team1, team2)
		//log.Printf("They play in %d minutes", (gameTime / 60))
		if isAlmostTime(gameTime) {
			if isValued(team1, team2) {
				text.SendText(fmt.Sprintf("%s vs %s in %d minutes", team1, team2, (gameTime / 60)))
			} else {
				text.SendText("Boring match soon")
			}
		} else {
			text.SendText("No maches soon")
		}
	}

}

func getMinutesUntil(matchHour int, ampm string) int64 {
	if ampm == "PM" {
		matchHour += 12
	}
	currentStamp := time.Now()
	newYork, _ := time.LoadLocation("America/New_York")
	futureGame := time.Date(currentStamp.Year(), currentStamp.Month(), currentStamp.Day(), matchHour, 0, 0, 0, newYork).Unix()
	remainingTime := futureGame - currentStamp.Unix()
	return remainingTime
}

func isAlmostTime(within int64) bool {
	if within <= 4000 && within >= 0 {
		return true
	}
	return false
}

func isValued(team1, team2 string) bool {
	for _, team := range valuedTeams {
		if team == team1 || team == team2 {
			return true
		}
	}
	return false
}

func isMatchToday(dates []*cdp.Node) bool {
	for _, day := range dates {
		dump := day.Dump(" ", " ", false)
		if strings.Contains(dump, "Today") == true {
			return true
		}
	}
	return false
}
