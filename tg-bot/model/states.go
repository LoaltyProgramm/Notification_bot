package model

type State string

const (
	StateMainMenu          State = "main_menu"
	StateRegistredText     State = "registred_text"
	StateRegistredInterval State = "registred_interval"
	StateRegistredGroup    State = "registred_group"
	StateRegistredFinal    State = "registred_final"
	StateRegistredError    State = "registred_error"
	StateIdle              State = "idle"
	StateEmptyLists        State = "empty_lists"
	StateAddREminder       State = "add_reminder"
	StateAllLists          State = "all_lists"
	StateList              State = "get_list"
	StateAddGroup          State = "add_state"
	StateWaitAddGroup      State = "wait_add_group"
	StateFinalAddGroup     State = "final_add_group"
	StateAllGroup          State = "all_group"
	StateRemoveGroup       State = "remove_group"
	StateErrorAddGroup     State = "error_status_add_group"
)
