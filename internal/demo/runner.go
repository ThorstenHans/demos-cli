package demo

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh"
)

var localPrompt = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#56990b"))
var localLog = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("#fff")).PaddingLeft(1)

var sshPrompt = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#3399cc"))
var sshOutput = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.NoColor{}).PaddingLeft(3)

var failurePrompt = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0000"))
var failure = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.NoColor{}).PaddingLeft(1)

func log(message string) {
	fmt.Printf("%s %s\n", localPrompt.Render("local>"), localLog.Render(strings.TrimSpace(message)))
}
func logFromSsh(message string) {
	if strings.Contains(message, "\n") {
		scanner := bufio.NewScanner(strings.NewReader(message))
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("%s %s\n", sshPrompt.Render("ssh>"), sshOutput.Render(line))
		}
		return
	}
	fmt.Printf("%s %s\n", sshPrompt.Render("ssh>"), sshOutput.Render(strings.TrimSpace(message)))
}
func logError(err error) {
	fmt.Printf("%s %s\n", failurePrompt.Render("error>"), failure.Render(err.Error()))
}

func (k Kind) String() string {
	switch k {
	case Markdown:
		return "Markdown"
	case Code:
		return "Code"
	default:
		return "Unknown"
	}
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func Run(demo DemoScript, cfg *Config) error {
	clearScreen()
	log(fmt.Sprintf("Running Demo: %s\n", demo.Name))
	log("Jumping into Space 🛸")
	config := &ssh.ClientConfig{
		User: cfg.JumpBoxUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.JumpBoxPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cfg.JumpBoxAddress, cfg.JumpBoxPort), config)
	if err != nil {
		return err
	}
	defer client.Close()
	logFromSsh("Connection established 🤝")
	for _, step := range demo.Steps {
		if step.Kind == Markdown {
			logFromSsh(step.Command)
		}
		if step.Kind == Code {
			logFromSsh(fmt.Sprintf("Running: %s", step.Command))
			output, err := runCommandOverSsh(step.Command, client)
			if err != nil {
				return err
			}
			logFromSsh(output)
		}

	}
	log("Landing on Earth again 🌎")
	return nil
}

func runCommandOverSsh(command string, client *ssh.Client) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		fmt.Println(failure.Render("Failed to run: " + err.Error()))
		return "", err
	}

	return b.String(), nil
}
