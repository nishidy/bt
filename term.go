package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"os/exec"
	"strings"
)

type Term struct {
	ypos int
	strs []string
}

func NewTerm(results []string) Term {
	term := Term{0, []string{}}
	for _, result := range results[1:] {
		term.strs = append(term.strs, result)
	}
	return term
}

func get_func_line(s string) (string, string) {
	for _, str := range strings.Split(s, " ") {
		if strings.Contains(str, "@L") {
			str_slice := strings.Split(str, "@L")
			file := str_slice[0]
			line := str_slice[1]
			return file, line
		}
	}
	return "", ""
}

func (t *Term) exec() {
	item := t.strs[t.ypos]
	file, line := get_func_line(item)
	if file != "" && line != "" {
		cmd := exec.Command("/usr/local/bin/vim", file, "+"+line)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func remove_ansi_code(str string) string {
	str_raw := str
	str_raw = strings.Replace(str_raw, "\x1b[34m", "|", -1)
	str_raw = strings.Replace(str_raw, "\x1b[31m", "~", -1)
	str_raw = strings.Replace(str_raw, "\x1b[0m", "^", -1)
	return str_raw
}

func (t *Term) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for y, str := range t.strs {
		bgcolor := termbox.ColorDefault
		if y == t.ypos {
			bgcolor = termbox.ColorWhite
		}

		color := termbox.ColorDefault
		str_raw := remove_ansi_code(str)
		i := 0
		for _, r := range str_raw {
			if r == '|' {
				color = termbox.ColorBlue
				continue
			}
			if r == '~' {
				color = termbox.ColorRed
				continue
			}
			if r == '^' {
				color = termbox.ColorDefault
				continue
			}
			termbox.SetCell(i, y, r, color, bgcolor)
			i += 1
		}
	}

	termbox.Flush()
}

func (t *Term) Run() {

	fmt.Println("# Finished search. You can open with vim to hit enter on next page.")
	fmt.Println("# Other available keys: ↓(C-j) ↑(C-k) Esc(C-q)")
	a := 0
	fmt.Scanln(&a)

	_ = termbox.Init()
	defer termbox.Close()

	t.draw()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc,
				termbox.KeyCtrlQ:
				return
			case termbox.KeyArrowDown,
				termbox.KeyCtrlJ:
				if t.ypos < len(t.strs)-1 {
					t.ypos += 1
				}
				t.draw()
			case termbox.KeyArrowUp,
				termbox.KeyCtrlK:
				if t.ypos > 0 {
					t.ypos -= 1
				}
				t.draw()
			case termbox.KeyEnter:
				t.exec()
				return
			}
		default:
			t.draw()
		}
	}
}