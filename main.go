package main

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/legolord208/stdutil"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"unicode"
)

var COLOR_ERROR = color.New(color.FgRed, color.Bold)

const RULES_FILE = ".ar_rules"

type Rule struct {
	Msg          string
	Exact        bool
	Reply        string
	From         []string
	NotFrom      []string
	InChannel    []string
	NotInChannel []string
}

var rules = make([]Rule, 0)

func loadRules() {
	content, err := ioutil.ReadFile(RULES_FILE)
	if err != nil {
		if !os.IsNotExist(err) {
			stdutil.PrintErr("Couldn't read rules file", err)
		}
	} else {
		err = json.Unmarshal(content, &rules)
		if err != nil {
			stdutil.PrintErr("Couldn't parse rules file", err)
			exit()
		}

		warn := false

		for i, rule := range rules {
			msg := ""
			for _, c := range rule.Msg {
				if unicode.IsUpper(c) {
					warn = true
					c = unicode.ToLower(c)
				} else if c == '\n' {
					warn = true
					c = ' '
				}
				msg += string(c)
			}
			rule.Msg = strings.TrimSpace(msg)
			rules[i] = rule
		}

		if warn {
			stdutil.PrintErr("Very funny. You can edit. We get it.", nil)
			return
		}
	}
}
func saveRules() bool {
	content, err := json.MarshalIndent(rules, "", "\t")
	if err != nil {
		stdutil.PrintErr("Couldn't generate rules file", err)
	} else {
		err = ioutil.WriteFile(RULES_FILE, content, 0666)
		if err != nil {
			stdutil.PrintErr("Couldn't write rules file", err)
		} else {
			return true
		}
	}
	return false
}

func main() {
	args := os.Args[1:]

	stdutil.EventPrePrintError = append(stdutil.EventPrePrintError, func(full string, msg string, err error) bool {
		color.Unset()
		COLOR_ERROR.Set()
		return false
	})
	stdutil.EventPostPrintError = append(stdutil.EventPostPrintError, func(full string, msg string, err error) {
		color.Unset()
	})

	if len(args) < 1 {
		doEdit()
	} else if strings.EqualFold(args[0], "run") {
		doRun(args[1:])
	} else {
		stdutil.PrintErr("Invalid mode '"+args[0]+"'.", nil)
	}
}

func exit() {
	runtime.Goexit()
}
