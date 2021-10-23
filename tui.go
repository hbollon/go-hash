package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
	"github.com/sirupsen/logrus"
)

// General stuff for styling the view
var (
	p *tea.Program

	errorColor = makeFgStyle("1")
	payload    string

	term    = termenv.ColorProfile()
	keyword = makeFgStyle("211")
	subtle  = makeFgStyle("241")
	dot     = colorFg(" • ", "236")

	color               = termenv.ColorProfile().Color
	focusedTextColor    = "205"
	focusedPrompt       = termenv.String("> ").Foreground(color("205")).String()
	blurredPrompt       = "> "
	focusedSubmitButton = "[ " + termenv.String("Submit").Foreground(color("205")).String() + " ]"
	blurredSubmitButton = "[ Submit ]"
)

func LaunchTui() {
	initialModel := model{
		0,
		0,
		false,
		30,
		0,
		progress.NewModel(progress.WithDefaultGradient()),
		false,
		views{
			getCreateTableInputs(),
			getUniqueInputView("Filename"),
			getUniqueInputView("Filename"),
			getCrackView(),
		},
		false,
		false,
		false,
		"",
	}
	p = tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

type tickMsg struct{}
type frameMsg struct{}
type errorMsg struct{}
type backToMenuMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

type model struct {
	Choice        int
	Index         int
	Chosen        bool
	Ticks         int
	Progress      float64
	ProgressBar   progress.Model
	InputProvided bool
	Views         views
	Loading       bool
	Loaded        bool
	Quitting      bool
	ErrorMsg      string
}

type views struct {
	createTable  inputs
	loadTable    inputs
	getTableInfo inputs
	crackHash    inputs
}

type inputs struct {
	index        int
	input        []textinput.Model
	submitButton string
}

func (m *model) returnToMenu() {
	m.Index = 0
	m.Choice = 0
	m.Chosen = false
	m.InputProvided = false
	m.Loaded = false
	m.Loading = false
	m.Ticks = 30
	m.ErrorMsg = ""
	payload = ""
	m.ProgressBar = progress.NewModel(progress.WithDefaultGradient())
	go func() {
		time.Sleep(1 * time.Second)
		p.Send(tickMsg{})
	}()
}

func getCreateTableInputs() inputs {
	i := inputs{
		input: []textinput.Model{
			textinput.NewModel(),
			textinput.NewModel(),
			textinput.NewModel(),
			textinput.NewModel(),
			textinput.NewModel(),
			textinput.NewModel(),
			textinput.NewModel(),
		},
		submitButton: blurredSubmitButton,
	}

	i.input[0].Placeholder = "Alphabet"
	i.input[0].Focus()
	i.input[0].Prompt = focusedPrompt
	i.input[0].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusedTextColor))

	i.input[1].Placeholder = "Minimum word size"
	i.input[1].Prompt = blurredPrompt

	i.input[2].Placeholder = "Maximum word size"
	i.input[2].Prompt = blurredPrompt

	i.input[3].Placeholder = "Table height"
	i.input[3].Prompt = blurredPrompt

	i.input[4].Placeholder = "Table width"
	i.input[4].Prompt = blurredPrompt

	i.input[5].Placeholder = "Hash method (SHA1/MD5)"
	i.input[5].Prompt = blurredPrompt

	i.input[6].Placeholder = "Filename"
	i.input[6].Prompt = blurredPrompt

	return i
}

func getCrackView() inputs {
	i := inputs{
		input: []textinput.Model{
			textinput.NewModel(),
			textinput.NewModel(),
		},
		submitButton: blurredSubmitButton,
	}

	i.input[0].Placeholder = "Text to hash and crack"
	i.input[0].Focus()
	i.input[0].Prompt = focusedPrompt
	i.input[0].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusedTextColor))

	i.input[1].Placeholder = "Hash method (MD5/SHA1)"
	i.input[1].Prompt = blurredPrompt

	return i
}

func getUniqueInputView(placeholder string) inputs {
	i := inputs{
		input: []textinput.Model{
			textinput.NewModel(),
		},
		submitButton: blurredSubmitButton,
	}

	i.input[0].Placeholder = placeholder
	i.input[0].Focus()
	i.input[0].Prompt = focusedPrompt
	i.input[0].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusedTextColor))

	return i
}

func (m model) Init() tea.Cmd {
	return tick()
}

// Main update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	switch message := msg.(type) {
	case tea.KeyMsg:
		k := message.String()
		if k == "esc" || k == "ctrl+c" {
			if m.Chosen {
				m.returnToMenu()
				return m, nil
			} else {
				m.Quitting = true
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.ProgressBar.Width = message.Width - 10*2 - 4
	}
	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.Choice += 1
			if m.Choice > 5 {
				m.Choice = 5
			}
		case "up":
			m.Choice -= 1
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "tab":
			m.Choice += 1
			if m.Choice > 5 {
				m.Choice = 0
			}
		case "enter":
			if (m.Choice == 3 || m.Choice == 4) && CurrentTable.table == nil {
				m.ErrorMsg = errorColor("No rainbow table loaded! Please create a new table or load an existing one first.\n\n")
				return m, nil
			}

			m.ErrorMsg = ""
			m.Chosen = true
			return m, frame()
		}

	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks -= 1
		return m, tick()
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch m.Choice {
	case 0:
		return updateCreateTable(msg, m)
	case 1:
		return updateLoadTable(msg, m)
	case 2:
		return updateGetTableInfo(msg, m)
	case 3:
		if payload == "" {
			payload = CurrentTable.Stats()
		}
		break
	case 4:
		return updateCrackHash(msg, m)
	default:
		break
	}

	return m, nil
}

func updateCreateTable(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if !m.Loading {
			switch msg.String() {
			case "down":
				m.Index += 1
				if m.Index > 7 {
					m.Index = 7
				}
			case "up":
				m.Index -= 1
				if m.Index < 0 {
					m.Index = 0
				}
			case "tab":
				m.Index += 1
				if m.Index > 7 {
					m.Index = 0
				}
			case "enter":
				if m.Index == 7 {
					minSize, _ := strconv.Atoi(m.Views.createTable.input[1].Value())
					maxSize, _ := strconv.Atoi(m.Views.createTable.input[2].Value())
					height, _ := strconv.Atoi(m.Views.createTable.input[3].Value())
					width, _ := strconv.Atoi(m.Views.createTable.input[4].Value())
					m.InputProvided = true
					m.Loading = true
					go func() {
						CurrentTable = CreateRaindowTable(
							height,
							width,
							GenerateAlphabet(
								m.Views.createTable.input[0].Value(),
								minSize,
								maxSize,
							),
							HashType(m.Views.createTable.input[5].Value()),
						)
						CurrentTable.Export(m.Views.createTable.input[6].Value())
						p.Send(backToMenuMsg{})
					}()
				}
				return m, frame()
			}
		}

	case frameMsg:
		if m.Loading {
			prog := m.ProgressBar.SetPercent(CurrentLoading.Percentage)
			return m, tea.Batch(frame(), prog)
		}

	case tickMsg:
		if m.Loaded && m.Chosen {
			if m.Ticks == 0 {
				m.returnToMenu()
				return m, nil
			}
			m.Ticks -= 1
			return m, tick()
		}

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.ProgressBar.Update(msg)
		m.ProgressBar = progressModel.(progress.Model)
		return m, cmd

	case backToMenuMsg:
		m.Loading = false
		m.Loaded = true
		m.Ticks = 10
		prog := m.ProgressBar.SetPercent(1.0)
		return m, tea.Batch(tick(), prog)
	}

	if !m.Loading && !m.Loaded {
		for i := 0; i < len(m.Views.createTable.input); i++ {
			if i == m.Index {
				// Set focused state
				m.Views.createTable.input[i].Focus()
				m.Views.createTable.input[i].Prompt = focusedPrompt
				m.Views.createTable.input[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusedTextColor))
				continue
			}
			// Remove focused state
			m.Views.createTable.input[i].Blur()
			m.Views.createTable.input[i].Prompt = blurredPrompt
			m.Views.createTable.input[i].TextStyle = lipgloss.Style{}
		}

		if m.Index == len(m.Views.createTable.input) {
			m.Views.createTable.submitButton = focusedSubmitButton
		} else {
			m.Views.createTable.submitButton = blurredSubmitButton
		}

		return updateInputs(msg, m, m.Views.createTable.input)
	}
	return m, nil
}

func updateLoadTable(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if !m.Loading {
			switch msg.String() {
			case "down":
				m.Index += 1
				if m.Index > 1 {
					m.Index = 1
				}
			case "up":
				m.Index -= 1
				if m.Index < 0 {
					m.Index = 0
				}
			case "tab":
				m.Index += 1
				if m.Index > 1 {
					m.Index = 0
				}
			case "enter":
				if m.Index == 1 {
					m.InputProvided = true
					m.Loading = true
					go func() {
						CurrentTable.Import(m.Views.loadTable.input[0].Value())
						p.Send(backToMenuMsg{})
					}()
				}
				return m, nil
			}
		}

	case tickMsg:
		if m.Loaded && m.Chosen {
			if m.Ticks == 0 {
				m.returnToMenu()
				return m, nil
			}
			m.Ticks -= 1
			return m, tick()
		}

	case backToMenuMsg:
		m.Loading = false
		m.Loaded = true
		m.Ticks = 10
		return m, tick()
	}

	if !m.Loading && !m.Loaded {
		for i := 0; i < len(m.Views.loadTable.input); i++ {
			if i == m.Index {
				// Set focused state
				m.Views.loadTable.input[i].Focus()
				m.Views.loadTable.input[i].Prompt = focusedPrompt
				m.Views.loadTable.input[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusedTextColor))
				continue
			}
			// Remove focused state
			m.Views.loadTable.input[i].Blur()
			m.Views.loadTable.input[i].Prompt = blurredPrompt
			m.Views.loadTable.input[i].TextStyle = lipgloss.Style{}
		}

		if m.Index == len(m.Views.loadTable.input) {
			m.Views.loadTable.submitButton = focusedSubmitButton
		} else {
			m.Views.loadTable.submitButton = blurredSubmitButton
		}

		return updateInputs(msg, m, m.Views.loadTable.input)
	}
	return m, nil
}

func updateGetTableInfo(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if !m.Loading {
			switch msg.String() {
			case "down":
				m.Index += 1
				if m.Index > 1 {
					m.Index = 1
				}
			case "up":
				m.Index -= 1
				if m.Index < 0 {
					m.Index = 0
				}
			case "tab":
				m.Index += 1
				if m.Index > 1 {
					m.Index = 0
				}
			case "enter":
				if m.Index == 1 {
					m.InputProvided = true
					m.Loading = true
					go func() {
						CurrentTable.Import(m.Views.getTableInfo.input[0].Value())
						payload = CurrentTable.Print()
						CurrentTable = RainbowTable{}
						p.Send(backToMenuMsg{})
					}()
				}
				return m, nil
			}
		}

	case backToMenuMsg:
		m.Loading = false
		m.Loaded = true
		return m, nil
	}

	if !m.Loading && !m.Loaded {
		for i := 0; i < len(m.Views.getTableInfo.input); i++ {
			if i == m.Index {
				// Set focused state
				m.Views.getTableInfo.input[i].Focus()
				m.Views.getTableInfo.input[i].Prompt = focusedPrompt
				m.Views.getTableInfo.input[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusedTextColor))
				continue
			}
			// Remove focused state
			m.Views.getTableInfo.input[i].Blur()
			m.Views.getTableInfo.input[i].Prompt = blurredPrompt
			m.Views.getTableInfo.input[i].TextStyle = lipgloss.Style{}
		}

		if m.Index == len(m.Views.getTableInfo.input) {
			m.Views.getTableInfo.submitButton = focusedSubmitButton
		} else {
			m.Views.getTableInfo.submitButton = blurredSubmitButton
		}

		return updateInputs(msg, m, m.Views.getTableInfo.input)
	}
	return m, nil
}

func updateCrackHash(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if !m.Loading {
			switch msg.String() {
			case "down":
				m.Index += 1
				if m.Index > 2 {
					m.Index = 1
				}
			case "up":
				m.Index -= 1
				if m.Index < 0 {
					m.Index = 0
				}
			case "tab":
				m.Index += 1
				if m.Index > 2 {
					m.Index = 0
				}
			case "enter":
				if m.Index == 2 {
					m.InputProvided = true
					m.Loading = true
					go func(m *model) {
						hash, err := Hash(m.Views.crackHash.input[0].Value(), HashType(m.Views.crackHash.input[1].Value()))
						if err != nil {
							p.Send(errorMsg{})
							return
						}
						payload, err = CurrentTable.Invert(hash)
						if err != nil {
							p.Send(errorMsg{})
							return
						}
						p.Send(backToMenuMsg{})
					}(&m)
				}
				return m, frame()
			}
		}

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.ProgressBar.Update(msg)
		m.ProgressBar = progressModel.(progress.Model)
		return m, cmd

	case frameMsg:
		if m.Loading {
			prog := m.ProgressBar.SetPercent(CurrentLoading.Percentage)
			return m, tea.Batch(frame(), prog)
		}

	case tickMsg:
		if m.Loaded && m.Chosen {
			if m.Ticks == 0 {
				m.returnToMenu()
				return m, nil
			}
			m.Ticks -= 1
			return m, tick()
		}

	case errorMsg:
		m.Loading = false
		m.ErrorMsg = CurrentLoading.Error.Error()
		return m, nil

	case backToMenuMsg:
		m.Loading = false
		m.Loaded = true
		m.Ticks = 10
		prog := m.ProgressBar.SetPercent(1.0)
		return m, tea.Batch(tick(), prog)
	}

	if !m.Loading && !m.Loaded {
		for i := 0; i < len(m.Views.crackHash.input); i++ {
			if i == m.Index {
				// Set focused state
				m.Views.crackHash.input[i].Focus()
				m.Views.crackHash.input[i].Prompt = focusedPrompt
				m.Views.crackHash.input[i].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusedTextColor))
				continue
			}
			// Remove focused state
			m.Views.crackHash.input[i].Blur()
			m.Views.crackHash.input[i].Prompt = blurredPrompt
			m.Views.crackHash.input[i].TextStyle = lipgloss.Style{}
		}

		if m.Index == len(m.Views.crackHash.input) {
			m.Views.crackHash.submitButton = focusedSubmitButton
		} else {
			m.Views.crackHash.submitButton = blurredSubmitButton
		}

		return updateInputs(msg, m, m.Views.crackHash.input)
	}
	return m, nil
}

func updateInputs(msg tea.Msg, m model, inputs []textinput.Model) (model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	for i := 0; i < len(inputs); i++ {
		inputs[i], cmd = inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Sub-views

// The first view, where you're choosing a task
func choicesView(m model) string {
	c := m.Choice

	tpl := fmt.Sprintf("Welcome to %s project! What do you want to do?\n\n", keyword("go-hash"))
	tpl += m.ErrorMsg
	tpl += "%s\n\n"
	tpl += "Program quits in %s seconds\n\n"
	tpl += subtle("j/k, up/down, tab: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s",
		checkbox("Create and export rainbow table", c == 0),
		checkbox("Load existing rainbow table", c == 1),
		checkbox("Get table informations from file", c == 2),
		checkbox("Compute current table statistics", c == 3),
		checkbox("Crack hash using current rainbow table", c == 4),
		checkbox("Execute specific demonstration test", c == 5),
	)

	return fmt.Sprintf(tpl, choices, colorFg(strconv.Itoa(m.Ticks), "79"))
}

// The second view, after a task has been chosen
func chosenView(m model) string {
	var msg, errorField, content string

	switch m.Choice {
	case 0:
		if !m.InputProvided {
			msg = fmt.Sprintf("Okay, so we need you to define an associeted %s, a table %s, a table %s, a %s and a %s...",
				keyword("alphabet"), keyword("height"), keyword("width"), keyword("hash method"), keyword("destination filename"))
			if errorField != "" {
				msg += fmt.Sprintf("\n\n%s", errorColor(errorField))
			}
			for i := 0; i < len(m.Views.createTable.input); i++ {
				content += m.Views.createTable.input[i].View()
				if i < len(m.Views.createTable.input)-1 {
					content += "\n"
				}
			}
			content += "\n\n" + m.Views.createTable.submitButton
		} else if m.InputProvided {
			msg = fmt.Sprintf("Generating your %s, please wait...", keyword("rainbow table"))
			content = m.ProgressBar.View()
			if m.Loaded {
				content += fmt.Sprintf("\nTable successfully generated!\nYou will be redirected to the main menu in %ss (press %s to go back now)...", colorFg(strconv.Itoa(m.Ticks), "79"), keyword("esc"))
			}
		}

	case 1:
		if !m.InputProvided {
			msg = fmt.Sprintf("Which %s file do you want to %s?", keyword("rainbow table"), keyword("load"))
			if errorField != "" {
				msg += fmt.Sprintf("\n\n%s", errorColor(errorField))
			}
			content += m.Views.loadTable.input[0].View()
			content += "\n\n" + m.Views.loadTable.submitButton
		} else if m.InputProvided {
			msg = fmt.Sprintf("Loading your %s, please wait...", keyword("rainbow table"))
			if m.Loaded {
				content += fmt.Sprintf("Table successfully loaded!\nYou will be redirected to the main menu in %ss (press %s to go back now)...", colorFg(strconv.Itoa(m.Ticks), "79"), keyword("esc"))
			}
		}

	case 2:
		if !m.InputProvided {
			msg = fmt.Sprintf("There is all requested %s table informations:", keyword("rainbow"))
			if errorField != "" {
				msg += fmt.Sprintf("\n\n%s", errorColor(errorField))
			}
			content += m.Views.getTableInfo.input[0].View()
			content += "\n\n" + m.Views.getTableInfo.submitButton
		} else if m.InputProvided {
			if !m.Loaded {
				msg = fmt.Sprintf("Loading table and processing %s, please wait...", keyword("data"))
			} else {
				msg = fmt.Sprintf("Table %s! There is its data:", keyword("loaded"))
				content += payload
			}
		}

	case 3:
		msg = fmt.Sprintf("Current table %s:", keyword("statistics"))
		content = payload

	case 4:
		if !m.InputProvided {
			msg = fmt.Sprintf("Which %s do you want to try to crack?", keyword("hash"))
			if errorField != "" {
				msg += fmt.Sprintf("\n\n%s", errorColor(errorField))
			}
			for i := 0; i < len(m.Views.crackHash.input); i++ {
				content += m.Views.crackHash.input[i].View()
				if i < len(m.Views.crackHash.input)-1 {
					content += "\n"
				}
			}
			content += "\n\n" + m.Views.crackHash.submitButton
		} else if m.InputProvided {
			msg = fmt.Sprintf("Cracking provided %s, please wait...", keyword("hash"))
			content = m.ProgressBar.View()
			if !m.Loading && m.ErrorMsg != "" {
				content += "\n" + errorColor(m.ErrorMsg)
			}
			if m.Loaded {
				if CurrentLoading.Error != nil {
					content += errorColor(CurrentLoading.Error.Error())
				} else {
					content += fmt.Sprintf("\nHash successfully cracked! Clair text: %s\nYou will be redirected to the main menu in %ss (press %s to go back now)...", keyword(CurrentLoading.Res), colorFg(strconv.Itoa(m.Ticks), "79"), keyword("esc"))
				}
			}
		}

	case 5:
		msg = fmt.Sprintf("There is a bunch of %s comming with this project, which one do you want to %s?", keyword("tests"), keyword("execute"))
	default:
		logrus.Fatal("Invalid input!")
	}

	return msg + "\n\n" + content + "\n\n" + subtle("up/down, tab: select") + dot + subtle("enter: choose") + dot + subtle("esc: quit")
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("> "+label, "212")
	}
	return fmt.Sprintf("  %s", label)
}

// Utils

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
