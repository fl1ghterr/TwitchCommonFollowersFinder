package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
)

type Follow struct {
	Login string `json:"login"`
}

type LoginCount struct {
	Login string
	Count int
}

func getFollows(username string) ([]string, error) {
	url := fmt.Sprintf("https://tools.2807.eu/api/getfollows/%s", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request error for %s: %v", username, err)
	}
	defer resp.Body.Close()

	var follows []Follow
	if err := json.NewDecoder(resp.Body).Decode(&follows); err != nil {
		return nil, fmt.Errorf("Error parsing JSON for %s: %v", username, err)
	}

	logins := make([]string, 0, len(follows))
	for _, follow := range follows {
		logins = append(logins, follow.Login)
	}
	return logins, nil
}

func worker(jobs <-chan string, results chan<- []string, wg *sync.WaitGroup, errLog chan<- error) {
	defer wg.Done()
	for username := range jobs {
		logins, err := getFollows(username)
		if err != nil {
			errLog <- err
			continue
		}
		results <- logins
	}
}

func main() {
	fmt.Println("🔄 Process running...")

	rawUsers := `
	example
	example`

	userList := strings.Fields(rawUsers)

	loginCount := make(map[string]int)
	var mutex sync.Mutex

	jobs := make(chan string, len(userList))
	results := make(chan []string, len(userList))
	errLog := make(chan error, len(userList))

	const workerCount = 10
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg, errLog)
	}

	for _, user := range userList {
		user = strings.TrimSpace(user)
		if user != "" {
			jobs <- user
		}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
		close(errLog)
	}()

	for logins := range results {
		mutex.Lock()
		for _, login := range logins {
			loginCount[login]++
		}
		mutex.Unlock()
	}

	for err := range errLog {
		fmt.Println("⚠️", err)
	}

	counts := []LoginCount{}
	for login, count := range loginCount {
		if count > 1 {
			counts = append(counts, LoginCount{Login: login, Count: count})
		}
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("❌ Error on file creation:", err)
		return
	}
	defer file.Close()

	for _, item := range counts {
		line := fmt.Sprintf("%s — %d times\n", item.Login, item.Count)
		file.WriteString(line)
	}

	fmt.Println("✅ Final result in output.txt")
}
