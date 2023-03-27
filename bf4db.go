package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var IsDebug = false
var homeDir, _ = os.UserHomeDir()
var envPath = filepath.Join(homeDir, "/.BF4DB-Search-Tool")
var apiKey string

func main() {
	// Verify if API key is set
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	key, err := godotenv.Read(envPath)
	if err != nil {
		fmt.Println("No API key found.")
		setApiKey()
	}
	apiKey = key["BF4DB_API_KEY"] // Set API key

	if len(os.Args) < 2 {
		fmt.Println("Usage: bf4db <player name> or -h for help")
		return
	}
	player := os.Args[1]

	if len(os.Args) > 2 {
		if os.Args[2] == "dbg" {
			IsDebug = true
		}
	}
	if player == "-config" || player == "-c" {
		// Set a new API key
		setApiKey()
	}
	if player == "-help" || player == "-h" || player == "-?" || player == "?" {
		// Show help
		fmt.Println("Usage: bf4db <player name> to search for a player, add 'dbg' to show debug info")
		fmt.Println("Usage: bf4db -c to set a new API key")
		fmt.Println("Usage: bf4db -h to show this weird help message")
		return
	}
	// Check if is an ip:port, useful when CTRC+C players IP from Procon Layer
	ip, _, err := net.SplitHostPort(player)
	if err == nil {
		player = ip
	}

	fmt.Println("Searching for " + player + "\n")
	GlobalSearch(player)
}

func setApiKey() {

	// Prompt user for API key
	fmt.Println("Please enter your BF4DB Patreon API key. You can get one here: https://bf4db.com/patreon")
	reader := bufio.NewReader(os.Stdin)
	apiKey, _ = reader.ReadString('\n')
	apiKey = apiKey[:len(apiKey)-1] // Remove newline character from end of input

	// check if API Key is a 64-character string
	if len(apiKey) != 64 {
		fmt.Println("API key is invalid. Please try again.")
		setApiKey()
	}
	// Save the API as env. variable:
	toWirte := map[string]string{
		"BF4DB_API_KEY": apiKey,
	}
	err := godotenv.Write(toWirte, envPath)
	if err != nil {
		fmt.Println("Error saving API key to .env file", err)
	}
	os.Exit(0)
}

func GlobalSearch(player string) {
	myUrl := fmt.Sprint("https://bf4db.com/api/player/", player, "/search?api_token=", apiKey) // url with API key
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, myUrl, nil)

	if err != nil {
		fmt.Println("NewRequest Error")
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Do Error")
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("ReadAll Error")
		return
	}

	var bfdbApi BFDBAPI
	err = json.Unmarshal(body, &bfdbApi)
	if err != nil {
		fmt.Println("Your API key is invalid! Please set a new one with bf4db -c")
		if IsDebug { // if debug is enabled, print the response body
			fmt.Println(string(body))
		}
		return
	}
	if len(bfdbApi.Data) == 0 {
		if IsDebug {
			fmt.Println(bfdbApi.Data, "No player found") // For debug only
		}
		return
	}
	// print number of players founds when > 15 (as it is harder to read)
	if len(bfdbApi.Data) > 15 {
		fmt.Println("More than 15 players found! Total of", len(bfdbApi.Data), "\n")
	}
	for x := range bfdbApi.Data {
		if bfdbApi.Data[x].BanReason == "" {
			bfdbApi.Data[x].BanReason = "Under review"
		}

		if IsDebug == true { // For debug only
			fmt.Println("Received Data:", bfdbApi.Data[x])
			continue
		}
		// if is nil, do nothing
		if a := bfdbApi.Data[x].ID; a == 0 {
			continue
		}

		bfdbURL := fmt.Sprint("https://bf4db.com/player/", bfdbApi.Data[x].ID, "/")
		bf4crURL := fmt.Sprint("http://bf4cheatreport.com/?pid=", bfdbApi.Data[x].ID, "&uid=&cnt=200&startdate=", time.Now().Format("200601021504"))
		//bfAcp := fmt.Sprint("https://BFACP/players?player=", bfdbApi.Data[x].Name) TODO: Add custom BFACP link to config (?)
		fmt.Printf("%v | %v | Cheat score = %v | %v\n Cheat Report: %v\n\n", bfdbApi.Data[x].Name, bfdbApi.Data[x].BanReason, bfdbApi.Data[x].CheatScore, bfdbURL, bf4crURL)
	}
}

type BFDBAPI struct {
	Data []struct {
		Name       string    `json:"name"`
		IsBanned   int       `json:"is_banned"`
		BanReason  string    `json:"ban_reason"`
		CheatScore int       `json:"cheat_score"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		ID         int       `json:"id"`
	} `json:"data"`
	Links struct {
		First string      `json:"first"`
		Last  string      `json:"last"`
		Prev  interface{} `json:"prev"`
		Next  interface{} `json:"next"`
	} `json:"links"`
	Meta struct {
		CurrentPage int    `json:"current_page"`
		From        int    `json:"from"`
		LastPage    int    `json:"last_page"`
		Path        string `json:"path"`
		PerPage     int    `json:"per_page"`
		To          int    `json:"to"`
		Total       int    `json:"total"`
	} `json:"meta"`
}
