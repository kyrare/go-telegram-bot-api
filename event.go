package bot

type Event struct {
	checker func(update Update) bool
	action  func(update Update)
}

func (event Event) ExecuteChecker(update Update) bool {
	return event.checker(update)
}

func (event Event) ExecuteAction(update Update) {
	event.action(update)
}
