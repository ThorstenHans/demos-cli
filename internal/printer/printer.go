package printer

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Printer interface {
	Print(message string)
}

type PrefixedPrinter struct {
	prefix       string
	prefixStyle  lipgloss.Style
	contentStyle lipgloss.Style
}

func NewRemote() Printer {
	return &PrefixedPrinter{
		prefix:       "ssh",
		prefixStyle:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#3399cc")),
		contentStyle: lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("#fff")).PaddingLeft(2),
	}
}

func NewInfo() Printer {
	return &PrefixedPrinter{
		prefix:       "info",
		prefixStyle:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#1fd008")),
		contentStyle: lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("#fff")).PaddingLeft(1),
	}
}

func NewFailure() Printer {
	return &PrefixedPrinter{
		prefix:       "error",
		prefixStyle:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0000")),
		contentStyle: lipgloss.NewStyle().Bold(false).Foreground(lipgloss.NoColor{}).PaddingLeft(1),
	}
}

func (p *PrefixedPrinter) Print(message string) {
	prompt := p.prefixStyle.Render(fmt.Sprintf("%s>", p.prefix))
	if strings.Contains(message, "\n") {
		scanner := bufio.NewScanner(strings.NewReader(message))
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("%s %s\n", prompt, p.contentStyle.Render(line))
		}
		return
	}
	fmt.Printf("%s %s\n", prompt, p.contentStyle.Render(strings.TrimSpace(message)))
}
