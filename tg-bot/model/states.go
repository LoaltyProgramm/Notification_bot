package model

type State string

const (
	StateMainMenu State = "main_menu"
	StateRegistredText State = "registred_text"
	StateRegistredInterval State = "registred_interval"
	StateRegistredFinal State = "registred_final"
	StateRegistredError State = "registred_error"
	StateIdle State = "idle"
	StateErrorInterval State = "error_interval"
)