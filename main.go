package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
)

func main() {
	// 1. Parse Flags
	// default is 5, usage: git-recent -n 10
	limitFlag := flag.Int("n", 5, "Number of recent branches to display")
	flag.Parse()

	// 2. Get valid local branches
	validBranches := getLocalBranches()

	// 3. Query git reflog
	// We fetch a bit more than requested to account for duplicates we filter out
	cmd := exec.Command("git", "log", "-g", "--grep-reflog=checkout: moving", "--pretty=%gs", "--max-count=100")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error: Not a git repository or git not found.")
		os.Exit(1)
	}

	// 4. Parse and Deduplicate
	var options []huh.Option[string]
	seen := make(map[string]bool)
	re := regexp.MustCompile(`moving from .*? to (.*)$`)

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			branchName := strings.TrimSpace(matches[1])

			if _, exists := validBranches[branchName]; exists {
				if !seen[branchName] {
					seen[branchName] = true
					options = append(options, huh.NewOption(branchName, branchName))
				}
			}
		}
	}

	if len(options) == 0 {
		fmt.Println("No recent branch history found.")
		return
	}

	// 5. Apply the limit (n)
	// We slice the list so the UI only ever receives 'n' items.
	limit := *limitFlag
	if limit > len(options) {
		limit = len(options)
	}
	options = options[:limit]

	// 6. Render the TUI
	var selectedBranch string

	// specific height ensures no internal scrolling unless n is huge
	form := huh.NewSelect[string]().
		Title(fmt.Sprintf("Last %d Branches (Ctrl+C to cancel)", limit)).
		Options(options...).
		Value(&selectedBranch).
		Height(limit + 4) 

	err = form.Run()
	if err != nil {
		return // User cancelled
	}

	if selectedBranch != "" {
		checkout(selectedBranch)
	}
}

func checkout(branch string) {
	fmt.Printf("Switching to %s...\n", branch)
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to switch: %v\n", err)
	}
}

func getLocalBranches() map[string]bool {
	branchMap := make(map[string]bool)
	out, err := exec.Command("git", "branch", "--format=%(refname:short)").Output()
	if err != nil {
		return branchMap
	}
	lines := strings.Split(string(out), "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			branchMap[trimmed] = true
		}
	}
	return branchMap
}