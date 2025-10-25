package demo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const appDirectory = ".demos"

type DemoScript struct {
	Name             string     `json:"name"`
	Command          string     `json:"cliCommand"`
	Alias            string     `json:"alias"`
	ShortDescription string     `json:"description"`
	Steps            []DemoStep `json:"steps"`
}

type DemoStep struct {
	Command string `json:"command"`
	Kind    Kind   `json:"kind"`
}

type Kind int

const (
	Markdown Kind = iota
	Code
)

func getDemosFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, appDirectory)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "demos.json"), nil
}

func LoadAll() []DemoScript {
	path, err := getDemosFilePath()
	if err != nil {
		return GetDefaultDemos()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return GetDefaultDemos()
	}
	var demos []DemoScript
	if err := json.Unmarshal(data, &demos); err != nil {
		fmt.Printf("Warn: Could not unmarshal demos file, using default demos instead\n")
		return GetDefaultDemos()
	}
	return demos
}

func GenerateSampleDemosFile() error {
	defaultDemos := GetDefaultDemos()
	jsonData, err := json.Marshal(defaultDemos)
	if err != nil {
		return err
	}
	path, err := getDemosFilePath()
	if err != nil {
		fmt.Printf("Err: Could create app folder in user home")
		return err
	}
	return os.WriteFile(path, jsonData, 0600)
}

func HasDemosFile() bool {
	path, err := getDemosFilePath()
	if err != nil {
		return false
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func ValidateDemos(demos []DemoScript) error {
	commands := make(map[string]bool)
	aliases := make(map[string]bool)

	for _, d := range demos {
		if commands[d.Command] {
			return fmt.Errorf("duplicate command found: %q", d.Command)
		}
		commands[d.Command] = true

		if aliases[d.Alias] {
			return fmt.Errorf("duplicate alias found: %q", d.Alias)
		}
		aliases[d.Alias] = true
	}

	return nil
}

func GetDefaultDemos() []DemoScript {
	return []DemoScript{
		{
			Name:             "Sample Load Test",
			Command:          "load-test",
			Alias:            "lt",
			ShortDescription: "Run the sample load test",
			Steps: []DemoStep{
				{Kind: Markdown, Command: "We'll sent 100 requests to Google now"},
				{Kind: Code, Command: "which hey"},
				{Kind: Code, Command: "hey -c 10 -n 100 https://www.google.com"},
				{Kind: Markdown, Command: "100 requests sent!"},
			},
		},
	}
}
