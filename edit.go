package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/legolord208/stdutil"
	"io"
	"strconv"
	"strings"
)

func doEdit() {
	loadRules()

	COLOR_BACK := color.New(color.BgBlack, color.FgWhite)
	READLINE, err := readline.New("")
	if err != nil {
		stdutil.PrintErr("Could not init readline library.", nil)
		return
	}

	COLOR_BACK.Println("Edit mode.")
	COLOR_BACK.Println("To start normal mode, do 'DiscordAR run <token>'")
	fmt.Println()
	fmt.Println("Get started! Create rule with 'new'. List rules with 'rules'. Edit rules with 'edit'")

	for i := 0; i < 3; i++ {
		fmt.Println()
	}
	for {
		READLINE.SetPrompt("")
		line := rlWrapper(READLINE.Readline())

		switch line {
		case "exit":
			return
		case "new":
			fmt.Println()
			fmt.Println("What will somebody type to trigger this reply?")
			READLINE.SetPrompt(COLOR_BACK.Sprint("Message:"))
			msg := rlWrapper(READLINE.Readline())

			fmt.Println()
			fmt.Println("Will this be exact? (say 'true' or 'false')")
			var exact bool
			for {
				READLINE.SetPrompt(COLOR_BACK.Sprint("Exact:"))
				exactStr := rlWrapper(READLINE.Readline())

				if exactStr == "true" {
					exact = true
				} else if exactStr == "false" {
					exact = false
				} else {
					fmt.Println("Invalid response.")
					continue
				}
				break
			}

			fmt.Println()
			fmt.Println("What is the reply?")
			READLINE.SetPrompt(COLOR_BACK.Sprint("Reply:"))
			reply := rlWrapperRaw(READLINE.Readline())

			rules = append(rules, Rule{
				Msg:   msg,
				Exact: exact,
				Reply: reply,
			})
			if saveRules() {
				fmt.Println("Saved rule!")
			}
		case "rules":
			for i, rule := range rules {
				fmt.Println("Rule #" + strconv.Itoa(i+1))
				fmt.Println("\tMessage: " + rule.Msg)
				fmt.Println("\tExact: " + strconv.FormatBool(rule.Exact))
				fmt.Println("\tReply: " + rule.Reply)
				fmt.Println("\tOnly from filter: " + strings.Join(rule.From, ", "))
				fmt.Println("\tNot from filter: " + strings.Join(rule.NotFrom, ", "))
				fmt.Println("\tOnly in channel filter: " + strings.Join(rule.InChannel, ", "))
				fmt.Println("\tNot in channel filter: " + strings.Join(rule.NotInChannel, ", "))
			}
		case "edit":
			fmt.Println("Select a rule.")
			fmt.Println()
			for i, rule := range rules {
				fmt.Println(strconv.Itoa(i+1) + ". " + rule.Msg)
			}

			fmt.Println()
			READLINE.SetPrompt(COLOR_BACK.Sprint("Rule:"))
			ruleStr := rlWrapper(READLINE.Readline())

			ruleNr, err := strconv.Atoi(ruleStr)
			if err != nil {
				stdutil.PrintErr("Not a number", nil)
				continue
			}

			ruleNr--
			if ruleNr < 0 || ruleNr >= len(rules) {
				stdutil.PrintErr("Rule does not exist", nil)
				continue
			}

			rule := rules[ruleNr]
			deleted := false

		property_loop:
			for {
				fmt.Println()
				fmt.Println("Select property to edit")
				fmt.Println()
				fmt.Println("1. Message")
				fmt.Println("2. Exact")
				fmt.Println("3. Reply")
				fmt.Println("4. Only from filter")
				fmt.Println("5. Not from filter")
				fmt.Println("6. Only in channel filter")
				fmt.Println("7. Not in channel filter")
				fmt.Println("8. Cancel")
				fmt.Println("9. DELETE!")

				fmt.Println()
				READLINE.SetPrompt(COLOR_BACK.Sprint("Property:"))
				prop := rlWrapper(READLINE.Readline())

				switch prop {
				case "1":
					READLINE.SetPrompt(COLOR_BACK.Sprint("Message:"))
					rule.Msg = rlWrapper(READLINE.Readline())
				case "2":
					rule.Exact = !rule.Exact
					fmt.Println("Exact toggled. Now: " + strconv.FormatBool(rule.Exact))
				case "3":
					READLINE.SetPrompt(COLOR_BACK.Sprint("Reply:"))
					rule.Reply = rlWrapperRaw(READLINE.Readline())
				case "4":
					arrEdit(&rule.From, READLINE, COLOR_BACK)
				case "5":
					arrEdit(&rule.NotFrom, READLINE, COLOR_BACK)
				case "6":
					arrEdit(&rule.InChannel, READLINE, COLOR_BACK)
				case "7":
					arrEdit(&rule.NotInChannel, READLINE, COLOR_BACK)
				case "8":
					break property_loop
				case "9":
					rules = append(rules[:ruleNr], rules[ruleNr+1:]...)
					deleted = true
					break property_loop
				}
			}

			if !deleted {
				rules[ruleNr] = rule
			}
			if saveRules() {
				fmt.Println("Saved!")
			}
		default:
			fmt.Println("Unknown command.")
			continue
		}
		for i := 0; i < 3; i++ {
			fmt.Println()
		}
	}
}

func rlWrapper(line string, err error) string {
	return strings.ToLower(strings.TrimSpace(rlWrapperRaw(line, err)))
}
func rlWrapperRaw(line string, err error) string {
	switch err {
	case nil:
		return line
	default:
		stdutil.PrintErr("Couldn't read line", err)
		fallthrough
	case io.EOF:
		fallthrough
	case readline.ErrInterrupt:
		exit()
		return ""
	}
}

func arrEdit(arr *[]string, READLINE *readline.Instance, COLOR_BACK *color.Color) {
	fmt.Println("Existing:")
	for _, item := range *arr {
		fmt.Println("\t" + item)
	}
	fmt.Println()
	fmt.Println("Now, write IDs like this:")
	fmt.Println("+ID to add.")
	fmt.Println("-ID to remove.")
	fmt.Println("Example: +123454321")

	fmt.Println()
	READLINE.SetPrompt(COLOR_BACK.Sprint("ID:"))
	id := rlWrapper(READLINE.Readline())

	if strings.HasPrefix(id, "+") {
		*arr = append(*arr, id[1:])
	} else if strings.HasPrefix(id, "-") {
		remove := id[1:]

		var arr2 []string
		for _, item := range *arr {
			if item != remove {
				arr2 = append(arr2, item)
			}
		}
		*arr = arr2
	} else {
		stdutil.PrintErr("Invalid.", nil)
	}
}
