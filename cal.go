package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	MAX_MONTH_DAY = 35
	HOLYDAY_API = "https://calendrier.api.gouv.fr/jours-feries/metropole/%d.json"
)

var holydayDays map[string]string

func main(){
	today := time.Now()
	year := flag.Int("year", today.Year(), "calendar year")
	month := flag.Int("month", int(today.Month()), "calendar month")
	flag.Parse()
	buildCalendar(*year, *month, today)
	fetch_holydays(*year, *month)
}

func fetch_holydays(year int, month int){
	req, err := http.NewRequest("GET", fmt.Sprintf(HOLYDAY_API, year), nil)
	if err != nil{
		fmt.Printf("error creating the request: %s", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil{
		fmt.Printf("error doing the request: %s", err)
	}
	jsonDec := json.NewDecoder(res.Body)
	jsonDec.Decode(&holydayDays)
	for date, name := range holydayDays{
		t, err := time.Parse(time.DateOnly, date)
		if month != int(t.Month()){
			continue
		}
		if err != nil{
			log.Printf("error parsing the time")
		}
		fmt.Printf("%s: \033[31m%s\033[0m\n", date[8:], name)
	}
}

func buildCalendar(year int, month int, today time.Time){
	calendarLoop := time.Date(year, time.Month(month), 1, 0,0,0,0, time.UTC)
	//Center Month Name
	for i := 0; i < (21 - len(calendarLoop.Month().String())) / 2; i++{
		fmt.Print(" ")
	} 
	fmt.Printf("\033[31m%s\033[0m\n",calendarLoop.Month())
	fmt.Printf("\033[33mL  M  M  J  V  S  D\033[0m\n")

	//Print empty space for the first days if needed
	if calendarLoop.Weekday() == time.Sunday{
		fmt.Printf("                  ")
	}else{
		for i := 0; i < int(calendarLoop.Weekday() - 1); i++{
			fmt.Print("   ")
		}
	}
	Calendar:
	for day := range MAX_MONTH_DAY{
		dynamicCalendar := calendarLoop.Add((time.Hour * 24 * time.Duration(day)))
		var isToday bool   
		if dynamicCalendar.Month() == today.Month() && dynamicCalendar.Day() == today.Day() && dynamicCalendar.Year() == today.Year(){
			isToday = true
		}
		if dynamicCalendar.Month() != calendarLoop.Month(){
			break Calendar
		}
		//Space Number
		if dynamicCalendar.Day() < 10{
			if isToday{
				fmt.Printf("\033[36m%d\033[0m  \n", dynamicCalendar.Day())
			}else if dynamicCalendar.Weekday() == time.Sunday{
				fmt.Printf("%d\n", dynamicCalendar.Day())
			}else{
				fmt.Printf("%d  ", dynamicCalendar.Day())
			}
			continue Calendar
		}
		if dynamicCalendar.Weekday() == time.Sunday && dynamicCalendar.Day() != 1{
			if isToday{
				fmt.Printf("\033[36m%d\033[0m\n", dynamicCalendar.Day())
			}else{
				fmt.Printf("%d\n", dynamicCalendar.Day())
			}
			continue Calendar
		}

		if isToday{
			fmt.Printf("\033[46;30m%d\033[0m ", dynamicCalendar.Day())
		}else{
			fmt.Printf("%d ", dynamicCalendar.Day())
		}

	}
	fmt.Printf("\n\n")
}
