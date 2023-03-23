/*
Sequential: An interactive task queue for single-cored organisms

Usage:

    sequential

The flags are:

    -c
        Skips interactive mode and prints the current task
        directly to the command line.

Sequential opens an interactive CLI that allows the user to
see a single task without distraction, queue additional tasks,
and randomly choose a different task.
*/
package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Enumeration of actions available to the user.
type TaskOption int8

const (
	MarkTaskComplete TaskOption = iota
	AddNewTask
	ShuffleTask
)

// String representation of TaskOption enum values.
var TaskOptionLabels = map[TaskOption]string{
	MarkTaskComplete: "Mark task complete",
	AddNewTask:       "Add new task",
	ShuffleTask:      "Shuffle task",
}

// Top-level model passed to the Bubbletea TUI module.
type model struct {
	config    Config
	tasks     *Logfile
	completed *Logfile
	textInput textinput.Model
	cursor    TaskOption
	selected  TaskOption
	appStyles AppStyles
}

type AppStyles struct {
	Selected  lipgloss.Style
	Completed lipgloss.Style
	Disabled  lipgloss.Style
	Title     lipgloss.Style
}

// Initializes the Bubbletea model
func createModel(config Config, tasks *Logfile, completed *Logfile) model {
	ti := textinput.New()
	ti.Placeholder = "Enter your task description"
	ti.Focus()

	appStyles := AppStyles{
		Selected: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(config.Styles.SelectedColor)),
		Completed: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(config.Styles.CompletedColor)),
		Title: lipgloss.NewStyle().
			Bold(true),
		Disabled: lipgloss.NewStyle().
			Foreground(lipgloss.Color(config.Styles.DisabledColor)),
	}

	return model{
		config:    config,
		tasks:     tasks,
		completed: completed,
		textInput: ti,
		selected:  -1,
		appStyles: appStyles,
	}
}

// Initializes the Bubbletea model. Required to satisfy
// the tea.Model interface.
func (m model) Init() tea.Cmd {
	return nil
}

// Called during the Bubbletea update phase. Required to
// satisfy the tea.Model interface.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil

	// User input has to be handled separately depending on whether
	// the user is on the selection menu or the new task prompt.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.selected == -1 {
			cmd = m.handleSelection(msg)
		} else if m.selected == AddNewTask {
			cmd = m.handleNewTask(msg)
		}
	}
	return m, cmd
}

// Handles user input on the new task prompt.
func (m *model) handleNewTask(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return tea.Quit
	case tea.KeyEnter:
		m.tasks.Lines = append(m.tasks.Lines, m.textInput.Value())
		m.tasks.Save()
		return tea.Quit
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return cmd
}

// Handles user input on the selection menu.
func (m *model) handleSelection(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "up", "k":
		min := MarkTaskComplete
		if len(m.tasks.Lines) == 0 {
			min = AddNewTask
		}
		if m.cursor > min {
			m.cursor--
		}
	case "down", "j":
		max := ShuffleTask
		if len(m.tasks.Lines) < 2 {
			max = AddNewTask
		}
		if m.cursor < max {
			m.cursor++
		}
	case "enter", " ":
		if m.cursor == ShuffleTask {
			if len(m.tasks.Lines) < 2 {
				return tea.Quit
			}
			swapIndex := rand.Intn(len(m.tasks.Lines)-1) + 1
			swap := m.tasks.Lines[swapIndex]
			m.tasks.Lines[swapIndex] = m.tasks.Lines[0]
			m.tasks.Lines[0] = swap
			m.tasks.Save()
		} else if m.cursor == 0 {
			m.completed.Lines = append(m.completed.Lines, m.tasks.Lines[0])
			m.completed.Save()
			m.tasks.Lines = m.tasks.Lines[1:]
			m.tasks.Save()
		} else {
			m.selected = m.cursor
			if m.selected == AddNewTask {
				return textinput.Blink
			}
		}
	}
	return nil
}

// Called during the Bubbletea view phase. Required to
// satisfy the tea.Model interface.
func (m model) View() string {
	var s string
	if len(m.tasks.Lines) == 0 {
		s = fmt.Sprintf(
			"%s\n\n",
			m.appStyles.Completed.Render("All tasks are complete!"),
		)
	} else {
		s = fmt.Sprintf(
			"%s\n%s\n\n",
			m.appStyles.Title.Render("Current Task"),
			m.tasks.Lines[0],
		)
	}

	if m.selected == -1 {
		s += m.renderSelectionView()
	} else if m.selected == 1 {
		s += m.textInput.View() + "\n"
	}

	return s
}

// Handle specific rendering tasks for the selection menu view.
func (m *model) renderSelectionView() string {
	var s string

	// Reposition the cursor to handle disabled cases.
	cursorPos := m.cursor
	lineLength := len(m.tasks.Lines)
	isEmpty := lineLength == 0
	isShuffable := lineLength >= 2
	if m.cursor == MarkTaskComplete && isEmpty {
		cursorPos = AddNewTask
	}
	if m.cursor == ShuffleTask && !isShuffable {
		cursorPos = AddNewTask
	}

	for i := 0; i < len(TaskOptionLabels); i++ {
		j := TaskOption(i)
		cursor := " "
		choice := TaskOptionLabels[j]
		if (!isShuffable && j == ShuffleTask) || (isEmpty && j == MarkTaskComplete) {
			choice = m.appStyles.Disabled.Render(choice)
		}
		if cursorPos == j {
			cursor = m.appStyles.Selected.Render(">")
			choice = m.appStyles.Selected.Render(choice)
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return s
}

func main() {
	tasks := LoadLogfile("current.txt")
	completed := LoadLogfile("completed.txt")
	config := LoadConfig()

	defer tasks.Close()
	defer completed.Close()

	flags := GetFlagBitmask()
	if flags&BIT_HELP != 0 {
		fmt.Println(HelpText)
		os.Exit(0)
	} else if flags&BIT_INLINE != 0 {
		if len(tasks.Lines) > 0 {
			fmt.Print(tasks.Lines[0])
		}
		os.Exit(0)
	}

	m := createModel(config, tasks, completed)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Unexpected error: %v", err)
		os.Exit(1)
	}
}
