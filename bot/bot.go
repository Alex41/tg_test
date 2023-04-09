package bot

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"sort"
	"tg_test/base"
)

var bot *tg.BotAPI

func ListenBot() {
	var err error
	bot, err = tg.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updChan := bot.GetUpdatesChan(tg.UpdateConfig{Timeout: 60})

	for upd := range updChan {
		if upd.Message != nil {
			dbMsg := base.Message{
				ID:      int64(upd.Message.MessageID),
				ChatID:  upd.Message.Chat.ID,
				Text:    upd.Message.Text,
				Creator: base.MessageCreatorUser,
			}

			if upd.Message.Document != nil {
				dbMsg.FileID = &upd.Message.Document.FileID
			}

			if len(upd.Message.Photo) > 0 { // there is one photo with difference sizes
				sort.Slice(upd.Message.Photo, func(i, j int) bool { // sort by size
					return upd.Message.Photo[i].FileSize > upd.Message.Photo[j].FileSize
				})
				dbMsg.FileID = &upd.Message.Photo[0].FileID // get the largest photo
			}

			//goland:noinspection GoUnhandledErrorResult
			base.SaveMessage(&dbMsg)
		}
	}
}

func Send(c tg.Chattable) (tg.Message, error) {
	return bot.Send(c)
}

func GetFileDirectURL(f string) (string, error) {
	return bot.GetFileDirectURL(f)
}
