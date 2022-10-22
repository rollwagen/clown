package prompter

type Prompter interface {
	Select(msg string, defaultValue string, options []string) (result int, err error)
}
