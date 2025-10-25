package demo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/thorstenhans/demos-over-ssh/internal/printer"
	"golang.org/x/crypto/ssh"
)

type Runner struct {
	info    printer.Printer
	remote  printer.Printer
	failure printer.Printer
}

func NewDemoRunner(info printer.Printer, ssh printer.Printer, failure printer.Printer) *Runner {
	return &Runner{
		info:    info,
		remote:  ssh,
		failure: failure,
	}
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

func (r *Runner) printHeader(name string) {
	clearScreen()
	r.info.Print("Starting Demo ...")
	r.info.Print(name)
	r.info.Print(fmt.Sprintf("\n\n🛸 Jumping into Space 🛸\n"))
}

func (r *Runner) Run(demo DemoScript, cfg *Config) error {
	r.printHeader(demo.Name)
	config := &ssh.ClientConfig{
		User: cfg.JumpBoxUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.JumpBoxPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", cfg.GetJumpBoxEndpoint(), config)
	if err != nil {
		return err
	}
	defer client.Close()
	r.remote.Print("🤝 First Contact (Connection established)")
	for _, step := range demo.Steps {
		if step.Kind == Markdown {
			r.info.Print(step.Command)
		}
		if step.Kind == Code {
			r.info.Print(fmt.Sprintf("👽 executing: %s (over ssh)", step.Command))
			output, err := r.runCommandOverSsh(step.Command, client)
			if err != nil {
				return err
			}
			r.remote.Print(output)
		}
	}
	r.info.Print("🪐 Expedition was successful, approaching earth...")
	r.info.Print("🌎 Back to earth")
	return nil
}

func (r *Runner) runCommandOverSsh(command string, client *ssh.Client) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		r.failure.Print(fmt.Sprintf("Failed to run: " + err.Error()))
		return "", err
	}
	return b.String(), nil
}
