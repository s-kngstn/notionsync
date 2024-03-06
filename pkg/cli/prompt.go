package cli

func Prompt(ui UserInput, prompt string) string {
	return ui.ReadString(prompt)
}
