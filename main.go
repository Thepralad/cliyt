package main

import (
	"log"
	"github.com/joho/godotenv"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}
var videoURLS []string
var maxResults int

// Struct for the YouTube API response
type YouTubeResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			ChannelTitle string `json:"channelTitle"`
		} `json:"snippet"`
	} `json:"items"`
}

type VideoDetails struct {
	Items []struct {
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

func apiRequest(query string) *YouTubeResponse {
	//API key
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	apiURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&q=%s&key=%s&maxResults=%d", query, apiKey, maxResults)

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: API returned status code %d\n", resp.StatusCode)
		return nil
	}

	var result YouTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return nil
	}

	return &result
}

func cmdSearch(parts []string) *YouTubeResponse {
	if len(parts) < 2 {
		fmt.Println("Error: Please provide a search query")
		fmt.Println("Usage: search <query>")
		return nil
	}
	query := strings.ReplaceAll(parts[1], " ", "+")
	fmt.Printf("\033[32mSearching for: %s\n\033[0m", parts[1])

	//Send an API request
	result := apiRequest(query)

	// Initialize the slice with the number of results
	videoURLS = make([]string, len(result.Items))

	for i, item := range result.Items {
		duration := getVideoDuration(item.ID.VideoID)
		formattedDuration := formatDuration(duration)
		fmt.Printf("\033[36m%d. %s\n   Channel: %s\n   Duration: %s\n\033[0m",
			i+1,
			item.Snippet.Title,
			item.Snippet.ChannelTitle,
			formattedDuration,
		)
		// Construct the full YouTube URL
		videoURLS[i] = fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.ID.VideoID)
		fmt.Println("----------------------------------")
	}
	return result
}

func getVideoDuration(videoID string) string {
	apiKey := "AIzaSyBZDcFHZluTvIqPCkX61cKJN7W3SV6KPzM"
	apiURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=contentDetails&id=%s&key=%s", videoID, apiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "Unknown duration"
	}
	defer resp.Body.Close()

	var details VideoDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return "Unknown duration"
	}

	if len(details.Items) > 0 {
		return details.Items[0].ContentDetails.Duration
	}
	return "Unknown duration"
}

func formatDuration(duration string) string {
	// Remove "PT" prefix and convert to readable format
	duration = strings.TrimPrefix(duration, "PT")

	hours := 0
	minutes := 0
	seconds := 0

	if strings.Contains(duration, "H") {
		parts := strings.Split(duration, "H")
		hours, _ = strconv.Atoi(parts[0])
		duration = parts[1]
	}
	if strings.Contains(duration, "M") {
		parts := strings.Split(duration, "M")
		minutes, _ = strconv.Atoi(parts[0])
		duration = parts[1]
	}
	if strings.Contains(duration, "S") {
		parts := strings.Split(duration, "S")
		seconds, _ = strconv.Atoi(parts[0])
	}

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func main() {
	var results *YouTubeResponse
	maxResults = 5 // Initialize maxResults here
	fmt.Println("\033[31mWelcome to YouTube CLI!\033[0m")
	fmt.Println("Commands:")
	fmt.Println("  search <query> - Search for videos")
	fmt.Println("  exit - Exit the program")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\033[31mcliyt > \033[0m")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		parts := strings.SplitN(input, " ", 2)
		command := parts[0]

		switch command {
		case "search":
			clearScreen()
			results = cmdSearch(parts)
		case "setMax":
			if len(parts) < 2 {
				fmt.Println("Error: Please provide a number")
				fmt.Println("Usage: setMax <number>")
				break
			}
			if num, err := strconv.Atoi(parts[1]); err == nil {
				if num > 0 && num <= 50 {
					maxResults = num
					fmt.Printf("\033[32mMaximum results set to: %d\033[0m\n", maxResults)
				} else {
					fmt.Println("Please enter a number between 1 and 50")
				}
			} else {
				fmt.Println("Please enter a valid number")
			}
		case "clear":
			clearScreen()
		default:
			// Check if input is a number
			if num, err := strconv.Atoi(command); err == nil {
				if num > 0 && num <= len(videoURLS) && results != nil {
					cmd := exec.Command("mpv", videoURLS[num-1])
					fmt.Printf("\033[32mNow playing, %s\033[0m\n", results.Items[num-1].Snippet.Title)
					if err := cmd.Start(); err != nil {
						fmt.Printf("Error playing video: %v\n", err)
					}
				} else {
					fmt.Println("Invalid video number or no search results. Please search first.")
				}
			} else {
				fmt.Printf("Unknown command: %s\n", command)
				fmt.Println("Available commands: search, exit")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
