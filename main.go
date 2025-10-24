package main

import(
	"encoding/json"
	"fmt"
	"os"
	"net/http"
)

type GitHubEvent struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		return
	}

	username := os.Args[1]
	url := "https://api.github.com/users/" + username + "/events"

	resp, err := http.Get(url)
	if err !=nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Println("User not found")
		return
	}

	var events []GitHubEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	if len(events) == 0 {
		fmt.Println("No recent activity found for user:", username)
		return
	}

	for _, events := range events {
		switch events.Type {
			case "PushEvent":
				fmt.Printf("- Pushed commits to: %s\n", events.Repo.Name)
			case "IssueEvents":
				fmt.Printf("- Worked an issue in: %s\n", events.Repo.Name)
			case "WatchEvent" :
				fmt.Printf("- Starred repository: %s\n", events.Repo.Name)
			default:
				fmt.Printf("- %s in repository: %s\n", events.Type, events.Repo.Name)
		}
	}
}