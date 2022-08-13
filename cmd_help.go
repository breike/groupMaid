package main

import (
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func maidGetHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *maidDB) (string, error) {
	msg_txt               := ""
	var err error          = nil
	var chat_id string     = strconv.Itoa(int(update.Message.Chat.ID))
	var key     string

	args_list := strings.Split(update.Message.Text, " ")

	if len(args_list) > 1 {
		key = args_list[1]
	} else {
		key = ""
	}

	if key != "" {
		if key[0:1] == "/" {
			key = key[1:]
		} 

		switch key {
		case db.Chats[chat_id].Config.BanCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.BanCmd + "`" +
				  " — бан пользователя по реплаю."
		case db.Chats[chat_id].Config.ConfigCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.ConfigCmd + "`" +
				" — конфигурация чата. Можно как получить информацию (по get), так и установить её (по set)." +
				"\n\nПример работы:" +
				"\n\n`/config set ban_command_on true`" +
			        "\n\n`/config get warns_action`" +
				"\n\nДля получения полного списка установленных параметров:" +
				"\n\n`/config get`" +
				"\n\nСписок параметров, которые можно установить:" +
			        "\n\n— `ban_command_on` — включает/отключает команду бана. Доступные значения: true или false." +
				"\n\n— `delete_last_welcome` — включает/отключает удаления прошлого приветствия. Доступные значения: true или false." +
				"\n\n— `disable_web_page_preview` — включает/отключает превью ссылок в сообщениях бота. Доступные значения: true или false." +
				"\n\n— `welcome_disable_web_page_preview` — включает/отключает превью ссылок в приветствии бота. Доступные значения: true или false." +
				"\n\n— `rules_disable_web_page_preview` — включает/отключает превью ссылок в правилах бота. Доступные значения: true или false." +
				"\n\n— `help_command_on` — включает/отключает команду помощи. Доступные значения: true или false." +
				"\n\n— `info_command_on` — включает/отключает команду паспорта. Доступные значения: true или false." +
				"\n\n— `mute_command_on` — включает/отключает команду мута. Доступные значения: true или false." +
				"\n\n— `rules_command_on` — включает/отключает команду правил. Доступные значения: true или false." +
				"\n\n— `welcome_on` — включает/отключает команду приветствия. Доступные значения: true или false." +
				"\n\n— `warns_limit` — Количество варнов до определённого действия (мут/бан/кик). Доступные значения: true или false." +
				"\n\n— `warns_action` — Действие бота при достижении определённого количества варнов у пользователя. Доступные значения: 0, 1, 2. Где 0 — бан, 1 — кик, 2 — мут." +
				"\n\n— `ban_command` — команда бана. На вход идёт строка без слеша." +
				"\n\n— `config_command` — команда настроек. На вход идёт строка без слеша." +
				"\n\n— `help_command` — команда помощи. На вход идёт строка без слеша." +
				"\n\n— `info_command.` — команда информации о пользователе. На вход идёт строка без слеша." +
				"\n\n— `kick_command` — команда кика. На вход идёт строка без слеша." +
				"\n\n— `mute_command` — команда мута. На вход идёт строка без слеша." +
				"\n\n— `remove_command` — команда удаления инфы о пользователе. На вход идёт строка без слеша." +
				"\n\n— `rules_command` — команда правил. На вход идёт строка без слеша." +
				"\n\n— `set_comand` — команда установления определённых настроек чата или бота. На вход идёт строка без слеша." +
				"\n\n— `setrules_command` — команда установления правил. На вход идёт строка без слеша." +
				"\n\n— `setwelcome_command` — команда установления текста приветствия. На вход идёт строка без слеша." +
				"\n\n— `status_command` — команда отображения статуса пользователя в чате (в чате/в бане и т.п.). На вход идёт строка без слеша." +
				"\n\n— `unmute_command` — команда анмута. На вход идёт строка без слеша." +
				"\n\n— `unset_command` — команда снятия определённых настроек. На вход идёт строка без слеша." +
				"\n\n— `update_command` — команда обновления базы данных. На вход идёт строка без слеша." +
				"\n\n— `warn_command` — команда варна. На вход идёт строка без слеша." +
				"\n\n— `welcome_command` — команда приветствия. На вход идёт строка без слеша."
		case db.Chats[chat_id].Config.HelpCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.HelpCmd + "`" +
				  " — команда помощи"
		case db.Chats[chat_id].Config.InfoCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.InfoCmd + "`" +
				  " — команда получения информации о пользователе. Работает при реплае тому пользователю, инфу о котором хочется получить."
		case db.Chats[chat_id].Config.KickCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.KickCmd + "`" +
				  " — кикает пользователя из чата при реплае на него. По факту, это бан на 30 секунд."
		case db.Chats[chat_id].Config.MuteCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.MuteCmd + "`" +
				  " — мутит пользователя при реплае. Можно указать время мута в команде (в минутах). Без указания мутит на 60 минут." +
				  "\n\nПример:" +
				  "\n\n`/mute 120`"
		case db.Chats[chat_id].Config.RemoveCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.RemoveCmd + "`" +
				  " — при реплае уменьшает какую-то характеристику пользователя на 1 (если эта характеристика исчисляема). По факту, используется лишь на ворнах." +
				  "\n\nПример:" +
				  "\n\n`/remove warn`"
		case db.Chats[chat_id].Config.RulesCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.RulesCmd + "`" +
				  " — отображает правила чата."
		case db.Chats[chat_id].Config.SetCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.SetCmd + "`" +
				  " — при реплае устанавливает какую-то характеристику пользователю." +
				  "\n\nПример:" +
				  "\n\n`/set notes омежка" +
			          "\n\nДоступные характеристики:" +
				  "\n\n— `adminnotes` — заметки, которые могут посмотреть лишь админы." +
				  "\n\n— `ban` — причина бана." +
				  "\n\n— `banfrom` — кем выдан бан. Обычно автоматически устанавливается first name того, кто забанил." +
				  "\n\n— `gender` — гендер." +
				  "\n\n— `pronouns` — местоимения." +
				  "\n\n— `privileges` — количество административного ресурса у пользователя. Устанавливается в числовом формате от 0 до 100. О привелегиях ниже." +
				  "\n\n— `name` — имя." +
				  "\n\n— `notes` — заметки." +
				  "\n\n— `warns` — количество ворнов. Устанавливается в числовом формате." +
				  "\n\nПривелегии разделяются по таким уровням:" +
				  "\n\n— 0 — обычный пользователь." +
				  "\n\n— 10 — может менять информацию о других пользователях." +
				  "\n\n— 50 — модератор, может мутить и банить." +
				  "\n\n— 70 — может устанавливать правила и приветствие." +
				  "\n\n— 100 — полные админские права."
		case db.Chats[chat_id].Config.SetrulesCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.SetrulesCmd + "`" +
				  " — Устанавливает правила чата. $name используется как first name вызвавшего пользователя."
		case db.Chats[chat_id].Config.SetwelcomeCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.SetwelcomeCmd + "`" +
				  " — Устанавливает правила чата. $name используется как first name вызвавшего пользователя."
		case db.Chats[chat_id].Config.StatusCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.StatusCmd + "`" +
				  " — Отображает статус пользователя при реплае. Обычно берутся значения из API телеграма." +
				  "\n\nЗначения:" +
				  "\n\n— Имя — берётся first name пользователя." +
				  "\n\n— В чате — берётся [IsMember()](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5#ChatMemberIsMember) из API телеграма." +
				  "\n\n— Кикнут — берётся [WasKicked()](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5#ChatMember.WasKicked) из API телеграма." +
				  "\n\n— Привелегии — количество административного ресурса у пользователя."
		case db.Chats[chat_id].Config.UnmuteCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.UnmuteCmd + "`" +
				  " — анмутит пользователя при реплае."
		case db.Chats[chat_id].Config.UpdateCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.UpdateCmd + "`" +
				  " — обновляет базу данных бота. Используется лишь при ручном изменении базы данных (т.е. не через бота)."
		case db.Chats[chat_id].Config.WarnCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.WarnCmd + "`" +
				  " — выписывает ворн пользователю при реплае."
		case db.Chats[chat_id].Config.WarnCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.WarnCmd + "`" +
				  " — выписывает ворн пользователю при реплае."
		case db.Chats[chat_id].Config.WelcomeCmd:
			msg_txt = "`/" + db.Chats[chat_id].Config.WelcomeCmd + "`" +
				  " — (тестовая фича) Отображает приветствие."
		default:
			msg_txt = "`" + key + "`" + " I don't understand."
		}
	} else {
		msg_txt = "Список доступных команд:" +
			  "\n\n — /" + db.Chats[chat_id].Config.BanCmd +
			  "— команда бана." +
			  "\n\n — /" + db.Chats[chat_id].Config.ConfigCmd +
			  "— команда конфигурации чата." +
			  "\n\n — /" + db.Chats[chat_id].Config.HelpCmd +
			  "— команда помощи." +
			  "\n\n — /" + db.Chats[chat_id].Config.InfoCmd +
			  "— команда получения информации о пользователе." +
			  "\n\n — /" + db.Chats[chat_id].Config.KickCmd +
			  "— команда кика пользователя." +
			  "\n\n — /" + db.Chats[chat_id].Config.MuteCmd +
			  "— команда мута пользователя." +
			  "\n\n — /" + db.Chats[chat_id].Config.RemoveCmd +
			  "— команда уменьшения характеристики пользователя на 1." +
			  "\n\n — /" + db.Chats[chat_id].Config.RulesCmd +
			  "— команда отображения правил чата." +
			  "\n\n — /" + db.Chats[chat_id].Config.SetCmd +
			  "— команда установления характеристики пользователя." +
			  "\n\n — /" + db.Chats[chat_id].Config.SetrulesCmd +
			  "— команда установления правил чата." +
			  "\n\n — /" + db.Chats[chat_id].Config.SetwelcomeCmd +
			  "— команда установления приветствия чата." +
			  "\n\n — /" + db.Chats[chat_id].Config.StatusCmd +
			  "— команда отображения статуса пользователя в чате." +
			  "\n\n — /" + db.Chats[chat_id].Config.UnmuteCmd +
			  "— команда отображения анмута пользователя." +
			  "\n\n — /" + db.Chats[chat_id].Config.UpdateCmd +
			  "— команда обновления базы данных." +
			  "\n\n — /" + db.Chats[chat_id].Config.WarnCmd +
			  "— команда ворна пользователя." +
			  "\n\n — /" + db.Chats[chat_id].Config.WelcomeCmd +
			  "— команда отображения приветствия чата." +
			  "\n\n Для отображения подробной информации об определённой команде:" +
			  "\n\n`/help %команда из предложенных выше%`"
	}

	return msg_txt, err
}
