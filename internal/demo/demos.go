package demo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

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

func getDemosPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".demo")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "demos.json"), nil
}

func LoadAll() []DemoScript {
	path, err := getDemosPath()
	if err != nil {
		fmt.Printf("Warn: Could not locate demos file, using default demos instead\n")
		return GetDefaultDemos()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Warn: Could not read demos file, using default demos instead\n")
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
	path, err := getDemosPath()
	if err != nil {
		fmt.Printf("Warn: Could not locate demos file, using default demos instead")
		return err
	}
	return os.WriteFile(path, jsonData, 0600)
}

func WannaGetDemosFileGenerated() bool {
	prompt := promptui.Select{
		Label: "Should I generate a demo file for you?",
		Items: []string{"Yes", "No"},
	}
	_, res, err := prompt.Run()
	if err != nil {
		fmt.Printf("Error while processing input: %s", err)
		return false
	}
	if strings.ToLower(res) == "yes" {
		return true
	}
	return false
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
			ShortDescription: "Run a Sample Load Test",
			Steps: []DemoStep{
				{Kind: Markdown, Command: "We'll now sent 100 requests Google"},
				{Kind: Code, Command: "hey -c 10 -n 100 https://www.google.com"},
				{Kind: Markdown, Command: "100 requests sent!"},
			},
		},
	}
}
