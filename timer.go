package main

import (
	"bufio"
	"encoding/json"
	"fmt"

	"os"
	"time"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
)

type ProjectData struct {
	Project   string    `json:"project"`
	TotalTime int       `json:"total_time"`
	Date      time.Time `json:"date_completed"`
}

//timer flags

func timerFlag() {
	fmt.Printf(`OPTIONS:
  -h, --help                  - usage help
  -v, --view               - view clock ins and clock outs
  -o, --open                - opens timer application
  -r, --regex                 - search string is a regular expression
  -k <key>, --key=<key>       - note key
  -t <title>, --title=<title> - title of note for create (CLI mode)
  -c <file>, --config=<file>  - config file to read from (defaults to ~/.snclirc)
`)
}

// displayTimer updates the terminal with the countdown

func displayTimer(duration time.Duration, message string) {
	for remaining := duration; remaining > 0; remaining -= time.Second {
		fmt.Printf("\r%s: %02d:%02d", message, remaining/time.Minute, remaining%time.Minute/time.Second)
		time.Sleep(time.Second)
	}
	color.HiMagenta("\r%vTime's up!", emoji.Unicorn)
}

// saveProjectData saves the project details to a JSON file

func saveProjectData(data ProjectData) {
	file, err := os.OpenFile("project_data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.HiMagenta("Error opening file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		color.HiMagenta("Error writing to JSON file:", err)
	}
}

// viewProjectData reads and displays saved session data

func viewProjectData() {
	file, err := os.Open("project_data.json")
	if err != nil {
		color.HiMagenta("Could not find an existing JSON file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var projects []ProjectData

	for decoder.More() {
		var project ProjectData
		if err := decoder.Decode(&project); err != nil {
			color.HiMagenta("Error decoding JSON:", err)
			return
		}
		projects = append(projects, project)
	}

	color.HiMagenta("\n%vSaved Project Data:", emoji.Unicorn)
	for _, p := range projects {
		fmt.Printf("Project: %s | Total Time: %d minutes | Date: %s\n", p.Project, p.TotalTime, p.Date.Format(time.RFC1123))
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var workMinutes, breakMinutes int
	var project string

	color.HiMagenta("Enter project details: ")
	project, _ = reader.ReadString('\n')
	project = project[:len(project)-1] // Remove newline character

	color.HiMagenta("%vEnter work session duration (in minutes): ", emoji.Unicorn)
	fmt.Scan(&workMinutes)

	color.HiMagenta("Enter break session duration (in minutes): ")
	fmt.Scan(&breakMinutes)

	totalTime := 0
	workDuration := time.Duration(workMinutes) * time.Minute
	breakDuration := time.Duration(breakMinutes) * time.Minute

	for {
		totalTime += workMinutes + breakMinutes
		date := time.Now()

		// Save project details to JSON

		saveProjectData(ProjectData{
			Project:   project,
			TotalTime: totalTime,
			Date:      date,
		})

		color.Magenta("\nTime to focus!", "Project:", project, "Total:", totalTime, "Minutes")
		displayTimer(workDuration, "Current work session:")

		color.HiMagenta("\nTime for a break!")
		displayTimer(breakDuration, "Break session")

		var continueResponse string
		fmt.Print("\nDo you want to start another session? (yes/no): ")
		fmt.Scan(&continueResponse)

		if continueResponse != "yes" {
			color.HiMagenta("Thanks for using the Pomodoro Timer. Stay productive!")
			break
		}

		var fResponse string
		fmt.Print("Press F if you would like to view session data: ")
		fmt.Scan(&fResponse)

		if fResponse == "f" || fResponse == "F" {
			viewProjectData()
		} else {
			color.HiMagenta("Have a good day")
		}
	}
}
