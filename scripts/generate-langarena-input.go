//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// Generate ~6MB of realistic Apache-style log lines for LangArena LogParser benchmark.
// Mirrors the data generation from https://kostya.github.io/LangArena/ (Etc::LogParser).
//
// Each line:
//   192.168.1.X - - [DD/Oct/2023:HH:55:36 +0000] "METHOD /path HTTP/1.1" STATUS 2326 "http://domain" "User-Agent"
//
// Mix includes login URLs with email/password params, API endpoints with tokens,
// session IDs, bot user-agents, suspicious paths, various HTTP methods and status codes.

const targetSize = 6 * 1024 * 1024 // 6 MB

var (
	methods = []string{"GET", "POST", "PUT", "DELETE"}
	paths   = []string{
		"/index.html", "/api/users", "/admin",
		"/images/logo.png", "/etc/passwd", "/wp-admin/setup.php",
	}
	statuses = []int{200, 201, 301, 302, 400, 401, 403, 404, 500, 502, 503}
	agents   = []string{
		"Mozilla/5.0", "Googlebot/2.1", "curl/7.68.0", "scanner/2.0",
	}
	users = []string{
		"john", "jane", "alex", "sarah", "mike", "anna", "david", "elena",
	}
	domains = []string{
		"example.com", "gmail.com", "yahoo.com", "hotmail.com", "company.org", "mail.ru",
	}
)

func generateLogData() string {
	ips := make([]string, 255)
	for i := range ips {
		ips[i] = fmt.Sprintf("192.168.1.%d", i+1)
	}

	var b strings.Builder
	b.Grow(targetSize + 4096)

	for i := 0; b.Len() < targetSize; i++ {
		b.WriteString(ips[i%len(ips)])
		fmt.Fprintf(&b, " - - [%d/Oct/2023:%d:55:36 +0000] \"", i%31, i%24)
		b.WriteString(methods[i%len(methods)])
		b.WriteString(" ")

		switch {
		case i%3 == 0:
			// Login URL with email and password (~33% of lines)
			fmt.Fprintf(&b, "/login?email=%s%d@%s&password=secret%d",
				users[i%len(users)], i%100,
				domains[i%len(domains)], i%10000)
		case i%5 == 0:
			// API endpoint with token (~20% of lines)
			b.WriteString("/api/data?token=")
			for j := 0; j < (i%3)+1; j++ {
				b.WriteString("abcdef123456")
			}
		case i%7 == 0:
			// Session ID (~14% of lines)
			fmt.Fprintf(&b, "/user/profile?session_id=sess_%x", i*12345)
		default:
			b.WriteString(paths[i%len(paths)])
		}

		fmt.Fprintf(&b, " HTTP/1.1\" %d 2326 \"http://%s\" %q\n",
			statuses[i%len(statuses)],
			domains[i%len(domains)],
			agents[i%len(agents)])
	}

	return b.String()
}

func main() {
	rand.Seed(42) // Fixed seed for reproducibility

	data := generateLogData()

	if err := os.MkdirAll("input", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating dir: %v\n", err)
		os.Exit(1)
	}

	err := os.WriteFile("input/langarena-data.txt", []byte(data), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated input/langarena-data.txt (%.2f MB)\n", float64(len(data))/1024/1024)
}
