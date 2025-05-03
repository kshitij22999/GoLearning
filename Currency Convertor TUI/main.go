package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const baseUrl string = "https://hexarate.paikama.co/api/rates/latest/"

type ApiResponse struct {
	StatusCode int  `json:"status_code"`
	Data       Data `json:"data"`
}

type Data struct {
	Base      string  `json:"base"`
	Target    string  `json:"target"`
	Mid       float64 `json:"mid"`
	Unit      int     `json:"unit"`
	Timestamp string  `json:"timestamp"`
}

type Model struct {
	choices   []string
	cursor    int
	stage     int // 0 = select base, 1 = select target, 2 = enter the amount
	base      string
	target    string
	textinput textinput.Model
	amount    float64
	style     *Style
	completed bool
}

type Style struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func getExchangeRate(base string, target string) float64 {
	response, err := http.Get(baseUrl + base + "?target=" + target)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", response.StatusCode)
	}
	body, errr := io.ReadAll(response.Body)
	if errr != nil {
		log.Fatal(err)
	}
	var result ApiResponse
	erro := json.Unmarshal([]byte(body), &result)
	if erro != nil {
		log.Fatal(erro)
	}

	return result.Data.Mid
}

func DefaultStyles() *Style {
	s := new(Style)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

func Newmodel() Model {
	styles := DefaultStyles()
	textinput := textinput.New()
	textinput.Placeholder = "Enter the amount"
	textinput.Focus()
	textinput.CharLimit = 20
	textinput.Width = 40
	return Model{
		choices:   []string{"INR", "USD", "EUR", "JPY", "AUD"},
		textinput: textinput,
		style:     styles,
	}
}

func main() {

	m := Newmodel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.stage == 0 {
				m.base = m.choices[m.cursor]
				m.cursor = 0 // reset cursor
				m.stage = 1  // move to target selection
			} else if m.stage == 1 {
				m.target = m.choices[m.cursor]
				m.cursor = 0
				m.stage = 2
			} else if m.stage == 2 {
				// Try to parse float amount
				amount, err := strconv.ParseFloat(m.textinput.Value(), 64)
				if err != nil {
					// optional: handle invalid input
					return m, nil
				}
				m.amount = amount
				m.completed = true
				return m, tea.Quit
			}
		}
	}
	if m.stage == 2 {
		var cmd tea.Cmd
		m.textinput, cmd = m.textinput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	var s string
	currencySymbols := map[string]string{
		"INR": "₹",
		"USD": "$",
		"EUR": "€",
		"JPY": "¥",
		"AUD": "A$",
	}
	if m.completed {
		s = ""
		res := m.amount * getExchangeRate(m.base, m.target)
		s += "\nSelection completed\n"
		s += "Base currency: " + m.base + "\n"
		s += "Target currency: " + m.target + "\n"
		s += "Amount: " + currencySymbols[m.base] + strconv.FormatFloat(m.amount, 'f', -1, 64) + "\n"
		s += "Converted amount is " + currencySymbols[m.target] + strconv.FormatFloat(res, 'f', -1, 64) + "\n"
		return (lipgloss.JoinVertical(
			lipgloss.Center,
			m.style.InputField.Render(s),
		))
	}

	if m.stage == 0 {
		s = "Select a base currency:\n\n"
	} else if m.stage == 1 {
		s = "Select a target currency to convert to:\n\n"
	} else {
		s = "Please enter the amount to convert\n\n"
	}

	if m.stage != 2 {
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	} else {
		m.textinput.Focus()
		s += m.textinput.View() + "\n"
	}
	s += "\nPress cntrl + c to quit\n"
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.style.InputField.Render(s),
	)
}
