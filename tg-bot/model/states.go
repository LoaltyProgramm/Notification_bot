package model

type State string

const (
	StateMainMenu State = "main_menu"
	StateRegistredText State = "registred_text"
	StateRegistredInterval State = "registred_interval"
	StateRegistredFinal State = "registred_final"
	StateRegistredError State = "registred_error"
	StateIdle State = "idle"
	StateEmptyLists State = "empty_lists"
	StateAddREminder State = "add_reminder"
	StateAllLists State = "all_lists"
	StateList State = "get_list"
)