package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"tg-app/internal/reminder"
	"tg-app/model"
	"tg-app/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var StateHandler = map[model.State]func(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService){
	model.StateMainMenu:          handlerMainMenu,
	model.StateRegistredText:     handlerRegistredText,
	model.StateRegistredInterval: handlerRegistredInterval,
	model.StateRegistredGroup:    handlerRegistredGroup,
	model.StateRegistredFinal:    handlerRegistredFinal,
	model.StateRegistredError:    handlerRegistredError,
	model.StateIdle:              handlerIdle,
	model.StateEmptyLists:        handlerEmptyLists,
	model.StateAddREminder:       handlerAddReminder,
	model.StateAllLists:          handlerAllLists,
	model.StateList:              handlerList,
	model.StateAddGroup:          handlerAddGroup,
	model.StateWaitAddGroup:      handlerWaitGroup,
	model.StateFinalAddGroup:     handlerFinalAddGroup,
	model.StateAllGroup:          handlerAllGroup,
	model.StateRemoveGroup:       handlerRemoveGroup,
	model.StateErrorAddGroup:     handlerErrorStatusAddGroup,
}

func handlerMainMenu(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	msg := tgbotapi.NewMessage(chatID, "<b>Выберите функцию👇</b>")
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать напоминание📋", "create_reminder"),
			tgbotapi.NewInlineKeyboardButtonData("Список напоминаний", "all_lists"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить группу", "add_group"),
			tgbotapi.NewInlineKeyboardButtonData("Список групп", "all_group"),
		),
	)
	infoMSG, err := h.Bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}

	session.RemoveMSG = infoMSG.MessageID
	session.RemoveMSGChatID = infoMSG.Chat.ID
}

func handlerRegistredText(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	msg := tgbotapi.NewMessage(chatID, "<b>Введите текст напоминания✍️</b>")
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
	infoMSG, err := h.Bot.Send(msg)
	if err != nil {
		log.Println("ERROR - ", err)
		return
	}

	session.RemoveMSG = infoMSG.MessageID
	session.RemoveMSGChatID = infoMSG.Chat.ID

	session.State = "registred_interval"
}

func handlerRegistredInterval(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	if session.RemoveMSGChatID != 0 && session.RemoveMSG != 0 {
		deleteMsg := tgbotapi.NewDeleteMessage(session.RemoveMSGChatID, session.RemoveMSG)
		if _, err := h.Bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		} else {
			session.RemoveMSGChatID = 0
			session.RemoveMSG = 0
		}

	}

	if session.IntervalRetry {
		msg := tgbotapi.NewMessage(chatID, "<b>Введите интервал напоминания⏰</b>")
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Назад к тексту", "redirect_registred_text"),
				tgbotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
			),
		)

		infoMSG, err := h.Bot.Send(msg)
		if err != nil {
			log.Println("ERROR - ", err)
			return
		}

		session.RemoveMSG = infoMSG.MessageID
		session.RemoveMSGChatID = infoMSG.Chat.ID

		session.IntervalRetry = false

		session.State = model.StateRegistredGroup
		return
	}

	session.UserText = update.Message.Text
	msg := tgbotapi.NewMessage(chatID, "<b>Введите интервал напоминания⏰</b>")
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад к тексту", "redirect_registred_text"),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
		),
	)

	infoMSG, err := h.Bot.Send(msg)
	if err != nil {
		log.Println("ERROR - ", err)
		return
	}

	session.RemoveMSG = infoMSG.MessageID
	session.RemoveMSGChatID = infoMSG.Chat.ID

	session.State = model.StateRegistredGroup
}

func handlerRegistredGroup(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	session.Interval = update.Message.Text

	if session.RemoveMSGChatID != 0 && session.RemoveMSG != 0 {
		deleteMsg := tgbotapi.NewDeleteMessage(session.RemoveMSGChatID, session.RemoveMSG)
		if _, err := h.Bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		} else {
			session.RemoveMSGChatID = 0
			session.RemoveMSG = 0
		}

	}

	var err error
	session.Reminder, err = utils.ParseIntervalData(session)
	if err != nil {
		log.Println(err)
		session.State = "registred_error"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не 1 правильный формат ввода интервала\nВведите интервал заново:")
		if _, err := h.Bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
		return
	}

	groups, err := service.ListGroups(context.Background(), session)
	if err != nil {
		log.Printf("Error handlerRegistredGroup in func ListGroups: %s", err)
		return
	}

	if len(groups) == 0 {
		session.State = model.StateMainMenu
		msg := tgbotapi.NewMessage(chatID, "У вас нет групп, чтобы прикрепить напоминие к ней.\nСначала добавьте группу, чтобы пользоваться ботом.")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Добавить группу", "add_group"),
			),
		)

		if _, err := h.Bot.Send(msg); err != nil {
			log.Printf("Error handlerRegistredGroup in send empty group: %s", err)
			return
		}
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, v := range groups {
		btn := tgbotapi.NewInlineKeyboardButtonData(v.TitleGroup, fmt.Sprintf("group_registred_%d", v.ID))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
	))

	msg := tgbotapi.NewMessage(chatID, "Выберите группу, куда присылать напоминание:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}

	log.Println("Sucsses")
}

func handlerRegistredFinal(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	if session.RemoveMSGChatID != 0 && session.RemoveMSG != 0 {
		deleteMsg := tgbotapi.NewDeleteMessage(session.RemoveMSGChatID, session.RemoveMSG)
		if _, err := h.Bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		} else {
			session.RemoveMSGChatID = 0
			session.RemoveMSG = 0
		}

	}

	var err error
	session.Reminder, err = utils.ParseIntervalData(session)
	if err != nil {
		log.Println(err)
		return
	}

	session.Reminder.GroupID = session.SendGroupIdint

	log.Println(session.Reminder.GroupID)

	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("<b>Подтверждаете напоминание?</b>\nТекст:\n%s\nИнтервал:\n%s\nОтправлять в группу:\n%v",
			session.UserText,
			session.Interval,
			session.SendGroupTitle))
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвержаю", "success_data"),
			tgbotapi.NewInlineKeyboardButtonData("Назад к интервалу", "back_interval"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "back"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func handlerRegistredError(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	session.Interval = update.Message.Text
	log.Println(session.Interval)
	var err error
	session.Reminder, err = utils.ParseIntervalData(session)
	log.Println(session.Reminder)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не 2 правильный формат ввода интервала\nВведите интервал заново:")
		if _, err := h.Bot.Send(msg); err != nil {
			log.Println(err)
			return
		}
		return
	}

	session.State = model.StateRegistredGroup
	handlerRegistredGroup(h, update, session, chatID, service)
}

func handlerIdle(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	return
}

func handlerEmptyLists(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	msg := tgbotapi.NewMessage(chatID, "Список пуст\nПополняй скорее его)")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func handlerAddReminder(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	err := service.Createreminder(context.Background(), session.Reminder)
	if err != nil {
		log.Fatal(err)
	}
	session.State = "main_menu"
	msg := tgbotapi.NewMessage(chatID, "Напоминание добавлено✅")
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}

	handlerMainMenu(h, update, session, chatID, service)
}

func handlerAllLists(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	reminders, err := service.ListRemindersForChatID(context.Background(), session)
	if err != nil {
		log.Println(err)
		return
	}

	if len(reminders) <= 0 {
		session.State = model.StateEmptyLists
		handlerEmptyLists(h, update, session, chatID, service)
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, v := range reminders {
		btn := tgbotapi.NewInlineKeyboardButtonData(v.Text, fmt.Sprintf("reminder_%d", v.ID))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
	))

	msg := tgbotapi.NewMessage(chatID, "Ваши напоминания:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	session.State = "idle"
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}

	handlerIdle(h, update, session, chatID, service)
}

func handlerList(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Текст:\n%s\n\nИнтервал:\n%s", session.Reminder.Text, session.Reminder.FullTime))
	idStr := strconv.Itoa(session.Reminder.ID)
	log.Println(session.Reminder)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить", fmt.Sprintf("delete_reminder_%s", idStr)), //error
			tgbotapi.NewInlineKeyboardButtonData("Назад", "all_lists"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func handlerAddGroup(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	msg := tgbotapi.NewMessage(chatID, "Добавьте бота в группу.\nОжидаем получение информации...")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "back_add_group"),
		),
	)
	infoMSG, err := h.Bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}

	msgID := infoMSG.MessageID
	removeMsgChatID := infoMSG.Chat.ID

	session.RemoveMSG = msgID
	session.RemoveMSGChatID = removeMsgChatID

	session.State = model.StateWaitAddGroup

}

func handlerWaitGroup(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	if update.MyChatMember == nil {
		return
	}

	if session.RemoveMSG == 0 {
		log.Println("No msg id")
	}
	deleteMsg := tgbotapi.NewDeleteMessage(session.RemoveMSGChatID, session.RemoveMSG)
	if _, err := h.Bot.Request(deleteMsg); err != nil {
		log.Println(err)
		return
	} else {
		session.RemoveMSGChatID = 0
		session.RemoveMSG = 0
	}

	session.Group.TitleGroup = update.MyChatMember.Chat.Title
	session.Group.UserID = update.MyChatMember.From.ID
	session.Group.GroupID = update.MyChatMember.Chat.ID

	log.Println(session.Group.UserID)

	msg := tgbotapi.NewMessage(session.Group.UserID, fmt.Sprintf("Вы хотите добавить бота в группу?\nНазвание группы - %s", session.Group.TitleGroup))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить", "add"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отмена", "remove"),
		),
	)

	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}

}
func handlerFinalAddGroup(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	err := service.CreateGroup(context.Background(), session)
	if err != nil {
		log.Println(err)
		return
	}

	msg := tgbotapi.NewMessage(session.Group.UserID, fmt.Sprintf("Группа - %s, добавлена!", session.Group.TitleGroup))
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}

	session.State = model.StateMainMenu

	handlerMainMenu(h, update, session, chatID, service)
}

func handlerAllGroup(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	log.Println("Start handler all group")
	groups, err := service.ListGroups(context.Background(), session)
	if err != nil {
		log.Println(err)
		return
	}

	if len(groups) == 0 {
		session.State = model.StateEmptyLists
		handlerEmptyLists(h, update, session, chatID, service)
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, gr := range groups {
		btn := tgbotapi.NewInlineKeyboardButtonData(gr.TitleGroup, fmt.Sprintf("group_%d", gr.ID))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
	))

	log.Println("Start send message")
	session.State = model.StateIdle
	msg := tgbotapi.NewMessage(chatID, "Добавленные группы:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}

	handlerIdle(h, update, session, chatID, service)
}

func handlerRemoveGroup(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	titleGroup, err := h.ServiceReminder.ListTitleGroup(context.Background(), session.RemoveGroup)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.ServiceReminder.RemoveGroup(context.Background(), session.RemoveGroup)
	if err != nil {
		log.Println(err)
		return
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Вы убрали бота из группы - %s\nМы автоматически удалили группу из вашего списка", titleGroup))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Готово", "redirect_main_menu"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}

func handlerErrorStatusAddGroup(h *Handler, update tgbotapi.Update, session *model.UserSession, chatID int64, service *reminder.ReminderService) {
	if session.RemoveMSG != 0 && session.RemoveMSGChatID != 0 {
		deleteMsg := tgbotapi.NewDeleteMessage(session.RemoveMSGChatID, session.RemoveMSG)
		if _, err := h.Bot.Request(deleteMsg); err != nil {
			log.Println(err)
			return
		} else {
			session.RemoveMSGChatID = 0
			session.RemoveMSG = 0
		}
	}

	msg := tgbotapi.NewMessage(chatID, "Для корректного добавления бота в группу, удалите его из группы куда добавили.\nДалее в главном меню нажмите кнопку - добавить бота. После чего добавьте обратно бота в нужную группу\nТак бот корректно добавиться в список!")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "redirect_main_menu"),
		),
	)
	if _, err := h.Bot.Send(msg); err != nil {
		log.Println(err)
		return
	}
}
