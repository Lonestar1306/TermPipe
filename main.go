package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type ToolType string

const (
	Grep     ToolType = "grep"
	Awk      ToolType = "awk"
	Sed      ToolType = "sed"
	Uniq     ToolType = "uniq"
	Wc       ToolType = "wc"
	Sort     ToolType = "sort"
	Head     ToolType = "head"
	Tail     ToolType = "tail"
	Redirect ToolType = "redirect"
	Jq       ToolType = "jq"
)


type Step struct {
	Tool      ToolType
	Pattern   string 
	Col       int    
	Delimiter string  
	Search    string 
	Replace   string 
	Flags     string 
}


type AppState int

const (
	StateNormal AppState = iota
	StateInputGrep
	StateInputAwk
	StateInputSed
	StateInputUniq
	StateInputWc
	StateInputSort
	StateInputHead
	StateInputTail
	StateInputRedirect
	StateInputJq
)


type model struct {
	inputText     string
	lines         []string
	steps         []Step
	state         AppState
	textInput     textinput.Model
	vp            viewport.Model
	ready         bool
	width         int
	height        int
	outputPreview string
	bashCmd       string
	err           error
}


var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#4af626")).MarginBottom(1)
	stepStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#e0e0e0")).Background(lipgloss.Color("#333333")).Padding(0, 1).MarginRight(1).MarginBottom(1)
	cmdStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#4af626")).Bold(true)
	paneStyle      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#555555")).Padding(0, 1)
	helpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Italic(true).MarginTop(1)
	highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Bold(true).Background(lipgloss.Color("#2a0000"))
)


func initialModel(inputData string) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40


	if inputData == "" {
		inputData = `save-ts;event-ts;op;D1;D2
2026-01-26 03:10:59;2026-01-26 02:10:31;13;45.89;13.49
2026-01-26 03:10:59;2026-01-26 02:09:08;12;3.54;9.25
2026-03-30 02:01:03,896 - INFO - Validazione e Sincronizzazione completate
2026-03-30 02:01:03,897 - INFO - Chiusura delle connessioni del ciclo...
{"log_level": "ERROR", "message": "Simulazione JSON", "code": 500}`
	}

	m := model{
		inputText: inputData,
		lines:     strings.Split(inputData, "\n"),
		steps:     []Step{},
		state:     StateNormal,
		textInput: ti,
	}
	
	return runPipeline(m)
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		switch m.state {
		case StateNormal:
			switch msg.String() {
			case "q", "esc":
				return m, tea.Quit
			case "g":
				m.state = StateInputGrep
				m.textInput.Prompt = "Grep Regex: "
				m.textInput.SetValue("")
				m.textInput.Focus()
			case "a":
				m.state = StateInputAwk
				
				bestDelim := ""
				var sampleLines []string
				

				for _, l := range strings.Split(m.outputPreview, "\n") {
					if strings.TrimSpace(l) != "" {
						sampleLines = append(sampleLines, l)
						if len(sampleLines) >= 4 {
							break
						}
					}
				}
				
				if len(sampleLines) > 0 {
					candidates := []string{";", ",", "|", "\t"}
					for _, cand := range candidates {
						firstCount := strings.Count(sampleLines[0], cand)

						if firstCount == 0 {
							continue
						}
						
						isConsistent := true

						for i := 1; i < len(sampleLines); i++ {
							if strings.Count(sampleLines[i], cand) != firstCount {
								isConsistent = false
								break
							}
						}
						

						if isConsistent {
							bestDelim = cand
							break
						}
					}
				}
				
				if bestDelim != "" {
					m.textInput.SetValue(bestDelim + " ")
					m.textInput.SetCursor(len(bestDelim) + 1)
				} else {
					m.textInput.SetValue("")
				}

				m.textInput.Prompt = "Awk sep col (es. ; 3): "
				m.textInput.Focus()
			case "s":
				m.state = StateInputSed
				m.textInput.Prompt = "Sed s/[]/[]/g: "
				m.textInput.SetValue("")
				m.textInput.Focus()
			case "u":
				m.state = StateInputUniq
				m.textInput.Prompt = "Uniq flags: "
				m.textInput.SetValue("")
				m.textInput.Focus()
			case "w":
				m.state = StateInputWc
				m.textInput.Prompt = "Wc flags: "
				m.textInput.SetValue("")
				m.textInput.Focus()
			case "o":
				m.state = StateInputSort
				m.textInput.Prompt = "Sort flags: "
				m.textInput.SetValue("")
				m.textInput.Focus()
			case "h":
				m.state = StateInputHead
				m.textInput.Prompt = "Head -n: "
				m.textInput.SetValue("")
				m.textInput.Focus()
			case "t":
				m.state = StateInputTail
				m.textInput.Prompt = "Tail -n: "
				m.textInput.SetValue("")
				m.textInput.Focus()
			case "x":
				m.state = StateInputRedirect
				m.textInput.Prompt = "Salva in (>): "
				m.textInput.SetValue("")
				m.textInput.Focus()
			case "j":
				m.state = StateInputJq
				m.textInput.Prompt = "Jq filter: "
				m.textInput.SetValue("")
				m.textInput.Focus()
			case "backspace", "delete":
				if len(m.steps) > 0 {
					m.steps = m.steps[:len(m.steps)-1]
					m = runPipeline(m)
				}
			default:
				var vpCmd tea.Cmd
				m.vp, vpCmd = m.vp.Update(msg)
				return m, vpCmd
			}

		case StateInputGrep, StateInputAwk, StateInputSed, StateInputUniq, StateInputWc, StateInputSort, StateInputHead, StateInputTail, StateInputRedirect, StateInputJq:
			if msg.Type == tea.KeyEsc {
				m.state = StateNormal 
			} else if msg.Type == tea.KeyEnter {
				val := m.textInput.Value()
				if val != "" || m.state == StateInputUniq || m.state == StateInputWc || m.state == StateInputSort || m.state == StateInputHead || m.state == StateInputTail {
					
					if (m.state == StateInputHead || m.state == StateInputTail) && val == "" {
						val = "10"
					}

					switch m.state {
					case StateInputGrep:
						m.steps = append(m.steps, Step{Tool: Grep, Pattern: val})
					case StateInputAwk:
						delim, col := parseAwkInput(val)
						if col > 0 {
							m.steps = append(m.steps, Step{Tool: Awk, Col: col, Delimiter: delim})
						}
					case StateInputSed:
						parts := strings.SplitN(val, "/", 2)
						if len(parts) == 2 {
							m.steps = append(m.steps, Step{Tool: Sed, Search: parts[0], Replace: parts[1]})
						}
					case StateInputUniq:
						m.steps = append(m.steps, Step{Tool: Uniq, Flags: strings.TrimSpace(val)})
					case StateInputWc:
						m.steps = append(m.steps, Step{Tool: Wc, Flags: strings.TrimSpace(val)})
					case StateInputSort:
						m.steps = append(m.steps, Step{Tool: Sort, Flags: strings.TrimSpace(val)})
					case StateInputHead:
						m.steps = append(m.steps, Step{Tool: Head, Flags: strings.TrimSpace(val)})
					case StateInputTail:
						m.steps = append(m.steps, Step{Tool: Tail, Flags: strings.TrimSpace(val)})
					case StateInputRedirect:
						m.steps = append(m.steps, Step{Tool: Redirect, Flags: strings.TrimSpace(val)})
					case StateInputJq:
						m.steps = append(m.steps, Step{Tool: Jq, Flags: strings.TrimSpace(val)})
					}
					m = runPipeline(m) 
				}
				m.state = StateNormal
			} else {
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		}

	case tea.MouseMsg:
		var vpCmd tea.Cmd
		m.vp, vpCmd = m.vp.Update(msg)
		return m, vpCmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		headerHeight := 13 
		
		if !m.ready {
			m.vp = viewport.New(msg.Width-4, msg.Height-headerHeight)
			m.vp.SetContent(m.outputPreview)
			m.ready = true
		} else {
			m.vp.Width = msg.Width - 4
			m.vp.Height = msg.Height - headerHeight
		}
	}

	return m, cmd
}

func runPipeline(m model) model {
	currentLines := make([]string, len(m.lines))
	copy(currentLines, m.lines)
	cmdParts := []string{"cat input.txt"} 
	
	var grepPatterns []string

	for _, step := range m.steps {
		switch step.Tool {
		case Grep:
			var nextLines []string
			re, err := regexp.Compile("(?i)" + step.Pattern) 
			if err == nil {
				grepPatterns = append(grepPatterns, step.Pattern) 
				for _, line := range currentLines {
					if re.MatchString(line) {
						nextLines = append(nextLines, line)
					}
				}
				currentLines = nextLines
			}
			cmdParts = append(cmdParts, fmt.Sprintf("grep -E -i '%s'", step.Pattern))

		case Awk:
			var nextLines []string
			for _, line := range currentLines {
				var parts []string
				if step.Delimiter == " " || step.Delimiter == "" {
					parts = strings.Fields(line)
				} else {
					parts = strings.Split(line, step.Delimiter)
				}
				if step.Col > 0 && step.Col <= len(parts) {
					nextLines = append(nextLines, parts[step.Col-1])
				} else {
					nextLines = append(nextLines, "")
				}
			}
			currentLines = nextLines
			if step.Delimiter == " " || step.Delimiter == "" {
				cmdParts = append(cmdParts, fmt.Sprintf("awk '{print $%d}'", step.Col))
			} else {
				cmdParts = append(cmdParts, fmt.Sprintf("awk -F'%s' '{print $%d}'", step.Delimiter, step.Col))
			}

		case Sed:
			var nextLines []string
			re, err := regexp.Compile(step.Search)
			for _, line := range currentLines {
				if err == nil {
					nextLines = append(nextLines, re.ReplaceAllString(line, step.Replace))
				} else {
					nextLines = append(nextLines, line)
				}
			}
			currentLines = nextLines
			cmdParts = append(cmdParts, fmt.Sprintf("sed -E 's/%s/%s/g'", step.Search, step.Replace))

		case Jq:
			var nextLines []string
			filter := strings.TrimSpace(step.Flags)
			for _, line := range currentLines {
				var data map[string]interface{}
				err := json.Unmarshal([]byte(line), &data)
				if err == nil && strings.HasPrefix(filter, ".") {
					key := filter[1:] 
					if val, ok := data[key]; ok {
						nextLines = append(nextLines, fmt.Sprintf("%v", val))
						continue
					}
				}
				nextLines = append(nextLines, line)
			}
			currentLines = nextLines
			cmdParts = append(cmdParts, fmt.Sprintf("jq '%s'", step.Flags))

		case Sort:
			if strings.Contains(step.Flags, "r") {
				sort.Sort(sort.Reverse(sort.StringSlice(currentLines)))
			} else {
				sort.Strings(currentLines)
			}
			cmdParts = append(cmdParts, fmt.Sprintf("sort %s", step.Flags))

		case Uniq:
			var nextLines []string
			if len(currentLines) > 0 {
				count := 1
				prev := currentLines[0]
				for i := 1; i <= len(currentLines); i++ {
					var curr string
					if i < len(currentLines) {
						curr = currentLines[i]
					}
					if i < len(currentLines) && curr == prev {
						count++
					} else {
						if strings.Contains(step.Flags, "c") {
							nextLines = append(nextLines, fmt.Sprintf("%7d %s", count, prev))
						} else {
							nextLines = append(nextLines, prev)
						}
						if i < len(currentLines) {
							prev = curr
							count = 1
						}
					}
				}
			}
			currentLines = nextLines
			cmdParts = append(cmdParts, fmt.Sprintf("uniq %s", step.Flags))

		case Wc:
			lines, words, chars := len(currentLines), 0, 0
			for _, l := range currentLines {
				words += len(strings.Fields(l))
				chars += len(l) + 1 
			}
			if strings.Contains(step.Flags, "l") {
				currentLines = []string{fmt.Sprintf("%d", lines)}
			} else if strings.Contains(step.Flags, "w") {
				currentLines = []string{fmt.Sprintf("%d", words)}
			} else if strings.Contains(step.Flags, "c") {
				currentLines = []string{fmt.Sprintf("%d", chars)}
			} else {
				currentLines = []string{fmt.Sprintf("%7d %7d %7d", lines, words, chars)}
			}
			cmdParts = append(cmdParts, fmt.Sprintf("wc %s", step.Flags))
			
		case Head:
			n, err := strconv.Atoi(step.Flags)
			if err == nil && n > 0 && n < len(currentLines) {
				currentLines = currentLines[:n]
			}
			cmdParts = append(cmdParts, fmt.Sprintf("head -n %s", step.Flags))
			
		case Tail:
			n, err := strconv.Atoi(step.Flags)
			if err == nil && n > 0 && n < len(currentLines) {
				currentLines = currentLines[len(currentLines)-n:]
			}
			cmdParts = append(cmdParts, fmt.Sprintf("tail -n %s", step.Flags))
			
		case Redirect:
			cmdParts = append(cmdParts, fmt.Sprintf("> %s", step.Flags))
		}
	}

	finalOutput := strings.Join(currentLines, "\n")
	
	for _, pattern := range grepPatterns {
		re, err := regexp.Compile("(?i)(" + pattern + ")")
		if err == nil {
			finalOutput = re.ReplaceAllStringFunc(finalOutput, func(match string) string {
				return highlightStyle.Render(match)
			})
		}
	}

	m.outputPreview = finalOutput

	bashCmd := ""
	if len(cmdParts) > 0 {
		bashCmd = cmdParts[0]
		for i := 1; i < len(cmdParts); i++ {
			if strings.HasPrefix(cmdParts[i], ">") {
				bashCmd += " " + cmdParts[i]
			} else {
				bashCmd += " | " + cmdParts[i]
			}
		}
	}
	m.bashCmd = bashCmd

	if m.ready {
		m.vp.SetContent(m.outputPreview)
		m.vp.GotoTop()
	}

	return m
}

func (m model) View() string {
	if m.width == 0 {
		return "Inizializzazione..."
	}

	header := titleStyle.Render("⚡ TermPipe Builder")

	var stepsView string
	if len(m.steps) == 0 {
		stepsView = "Nessun filtro. I dati passano inalterati.\n"
	} else {
		var parts []string
		for i, s := range m.steps {
			lbl := ""
			switch s.Tool {
			case Grep: lbl = fmt.Sprintf("[%d] grep '%s'", i+1, s.Pattern)
			case Awk:  
				if s.Delimiter == " " || s.Delimiter == "" {
					lbl = fmt.Sprintf("[%d] awk col %d", i+1, s.Col)
				} else {
					lbl = fmt.Sprintf("[%d] awk -F'%s' col %d", i+1, s.Delimiter, s.Col)
				}
			case Sed:      lbl = fmt.Sprintf("[%d] sed s/%s/%s/", i+1, s.Search, s.Replace)
			case Sort:     lbl = fmt.Sprintf("[%d] sort %s", i+1, s.Flags)
			case Uniq:     lbl = fmt.Sprintf("[%d] uniq %s", i+1, s.Flags)
			case Wc:       lbl = fmt.Sprintf("[%d] wc %s", i+1, s.Flags)
			case Head:     lbl = fmt.Sprintf("[%d] head -n %s", i+1, s.Flags)
			case Tail:     lbl = fmt.Sprintf("[%d] tail -n %s", i+1, s.Flags)
			case Jq:       lbl = fmt.Sprintf("[%d] jq '%s'", i+1, s.Flags)
			case Redirect: lbl = fmt.Sprintf("[%d] > %s", i+1, s.Flags)
			}
			parts = append(parts, stepStyle.Render(lbl))
		}
		stepsView = lipgloss.JoinHorizontal(lipgloss.Top, parts...) + "\n"
	}

	outputBlock := paneStyle.Width(m.width - 2).Render(m.vp.View())

	cmdView := "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).Render("Comando Bash generato:") + "\n"
	cmdView += cmdStyle.Render("$ " + m.bashCmd) + "\n"


	var awkLegend string
	if m.state == StateInputAwk {
		delim, _ := parseAwkInput(m.textInput.Value())
		firstLine := ""
		for _, l := range strings.Split(m.outputPreview, "\n") {
			if strings.TrimSpace(l) != "" {
				firstLine = l
				break
			}
		}
		
		if firstLine != "" {
			var columns []string
			if delim == " " || delim == "" {
				columns = strings.Fields(firstLine)
			} else {
				columns = strings.Split(firstLine, delim)
			}
			
			currentRow := ""
			for i, colText := range columns {
				colText = strings.TrimSpace(colText)
				if colText == "" {
					colText = "<vuoto>"
				} else if len(colText) > 15 {
					colText = colText[:12] + "..."
				}
				
				block := lipgloss.NewStyle().
					Background(lipgloss.Color("#222")).
					Foreground(lipgloss.Color("#d16969")).
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("#555")).
					Padding(0, 1).
					Render(fmt.Sprintf("%d: %s", i+1, colText))
				
				if lipgloss.Width(currentRow) + lipgloss.Width(block) + 5 > m.width - 4 {
					currentRow += lipgloss.NewStyle().Foreground(lipgloss.Color("#555")).Render(" ...")
					break
				}
				currentRow += block + " "
			}
			

			displayDelim := delim
			if delim == " " || delim == "" {
				displayDelim = "SPAZIO (Default)"
			} else if delim == "\t" {
				displayDelim = "TAB"
			}
			
			awkLegend = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")).Render(fmt.Sprintf("👀 Anteprima Colonne (separatore: %s):", displayDelim)) + "\n" + currentRow + "\n"
		}
	}

	var footerView string
	if m.state == StateNormal {
		footerView = helpStyle.Render("Cmd: [g]rep [a]wk [s]ed [j]q [o]sort [u]niq [w]c [h]ead [t]ail [x]Export [del]Undo [q]Esci\nScroll: Frecce / Mouse Wheel")
	} else {
		footerView = awkLegend + m.textInput.View() + helpStyle.Render("\n(Premi Invio per confermare, Esc per annullare)")
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", header, stepsView, outputBlock, cmdView, footerView)
}

func parseAwkInput(val string) (delim string, col int) {
	val = strings.TrimSpace(val)
	if val == "" {
		return " ", 0
	}
	val = strings.ReplaceAll(val, "'", "")
	val = strings.ReplaceAll(val, "\"", "")
	val = strings.TrimPrefix(val, "-F")

	fields := strings.Fields(val)
	if len(fields) == 1 {
		if c, err := strconv.Atoi(fields[0]); err == nil {
			return " ", c 
		}
		return fields[0], 0 
	}
	if c, err := strconv.Atoi(fields[0]); err == nil {
		return fields[1], c 
	}
	if c, err := strconv.Atoi(fields[1]); err == nil {
		return fields[0], c 
	}
	return " ", 0
}

func main() {
	var inputData string
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := io.ReadAll(os.Stdin)
		if err == nil {
			inputData = string(bytes)
		}
		
		tty, err := os.Open("/dev/tty")
		if err != nil {
			fmt.Println("Errore nell'apertura del terminale:", err)
			os.Exit(1)
		}
		defer tty.Close()
		
		p := tea.NewProgram(initialModel(inputData), tea.WithInput(tty), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Errore TUI: %v", err)
			os.Exit(1)
		}
	} else {
		p := tea.NewProgram(initialModel(""), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Errore TUI: %v", err)
			os.Exit(1)
		}
	}
}
