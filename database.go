/* database functions
 *
 * Bot's database scheme:
 *
 *   %bot_directory/
 *   |  chats/                 - there is chats' preferences and data 
 *   |  |
 *   |  |  %chat_id%/          - naming of these directories
 *   |  |  |                     should be ONLY by chats' ID
 *   |  |  |
 *   |  |  |  chat_config.toml - chat's preferences 
 *
 */
package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/naoina/toml"
)

type maidDB struct {
	Chats     map[string]*chat
}

type chat struct {
	Config    chatConfig
	Users     map[string]*user
}

type chatConfig struct {
	BanCommandOn                    bool
	DeleteLastWelcome               bool
	DisableWebPagePreview           bool
	WelcomeDisableWebPagePreview    bool
	RulesDisableWebPagePreview      bool
	HelpCommandOn                   bool
	InfoCommandOn                   bool
	MuteCommandOn                   bool
	RulesCommandOn                  bool
	WelcomeOn                       bool

	RulesMsg                        string
	WelcomeMsg                      string
}

type user struct {
	PreferredName    string
	Gender           string
	Notes            string
	BanNote          string
}

var Chat_cfg_defaults = chatConfig{
	BanCommandOn:        true,
	DeleteLastWelcome:   true,
	HelpCommandOn:       true,
	InfoCommandOn:       true,
	MuteCommandOn:       true,
	RulesCommandOn:      true,

	WelcomeOn:           true,
	RulesMsg:
		"There's no rules yet. " +
		"Set rules by /setrules command.",
	WelcomeMsg:
		"Welcome!",
}

func dbInit() (maidDB, error) {
	db := maidDB{}
	db.Chats = make(map[string]*chat)
	var err error

	db_dir    := filepath.Join(config.BotDirectory, "db")
	chats_dir := filepath.Join(db_dir, "chats")

	if _, err = os.Stat(chats_dir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(chats_dir, 0700)
			if err != nil {
				return db, err
			}
		} else {
			return db, err
		}
	}

	// load chats' preferences into RAM
	chats_dirs, err := ioutil.ReadDir(chats_dir)
	if err != nil {
		return db, err
	}

	for _, chat_dir := range chats_dirs {
		chat_id       := chat_dir.Name()
		chat_dir_path := filepath.Join(chats_dir, chat_id)

		// check if chat's directory doesn't exist
		if _, err := os.Stat(chat_dir_path); err != nil {
				return db, err
		}

		chat_config := filepath.Join(chat_dir_path, "chat_config.toml")

		// check if chat's prefeneces file doesn't exist 
		if _, err := os.Stat(chat_config); err != nil {
			if os.IsNotExist(err) {
				_, err = os.Create(chat_config)
				if err != nil {
					return db, err
				}
			} else {
				return db, err
			}
		}

		data, err := ioutil.ReadFile(chat_config)
		if err != nil {
			return db, err
		}
		if string(data) != "" {
			f, err := os.Open(chat_config)
			if err != nil {
				return db, err
			}
			defer f.Close()

			// decode chat_config.toml
			db.Chats[chat_id] = new(chat)
			err = toml.NewDecoder(f).Decode(&db.Chats[chat_id].Config)
			if err != nil {
				return db, err
			}
		} else {
			f, err := os.Open(chat_config)
			if err != nil {
				return db, err
			}
			defer f.Close()

			db.Chats[chat_id].Config = Chat_cfg_defaults

			err = toml.NewEncoder(f).Encode(db.Chats[chat_id].Config)
			if err != nil {
				return db, err
			}
		}

		userdb := filepath.Join(chat_dir_path, "userdb.toml")

		// check if users database file doesn't exist 
		if _, err := os.Stat(userdb); err != nil {
			if os.IsNotExist(err) {
				_, err = os.Create(userdb)
				if err != nil {
					return db, err
				}
			} else {
				return db, err
			}
		}

		data, err = ioutil.ReadFile(userdb)
		if err != nil {
			return db, err
		}
		if string(data) != "" {
			f, err := os.Open(userdb)
			if err != nil {
				return db, err
			}
			defer f.Close()

			// decode userdb.toml
			err = toml.NewDecoder(f).Decode(&db.Chats[chat_id].Users)
			if err != nil {
				return db, err
			}
		}
	}
	return db, err
}

func dbWriteChatConfig(chat_id string, chat_cfg chatConfig, db *maidDB) (error) {
	var err error = nil

	db_dir    := filepath.Join(config.BotDirectory, "db")
	chats_dir := filepath.Join(db_dir, "chats")

	chat_dir_path := filepath.Join(chats_dir, chat_id)

	if _, err = os.Stat(chat_dir_path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(chat_dir_path, 0700)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	chat_cfg_path := filepath.Join(chat_dir_path, "chat_config.toml")

	f, err := os.Create(chat_cfg_path)
	if err != nil {
		return err
	}

	err = toml.NewEncoder(f).Encode(chat_cfg)
	if err != nil {
		return err
	}

	db.Chats[chat_id].Config = chat_cfg

	return err
}

func dbWriteChatUsers(chat_id string, users map[string]*user) (error) {
	var err error = nil

	db_dir    := filepath.Join(config.BotDirectory, "db")
	chats_dir := filepath.Join(db_dir, "chats")

	chat_dir_path := filepath.Join(chats_dir, chat_id)

	if _, err = os.Stat(chat_dir_path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(chat_dir_path, 0700)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	chat_cfg_path := filepath.Join(chat_dir_path, "userdb.toml")

	f, err := os.Create(chat_cfg_path)
	if err != nil {
		return err
	}

	err = toml.NewEncoder(f).Encode(users)
	if err != nil {
		return err
	}

	return err
}
