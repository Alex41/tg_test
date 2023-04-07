package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"tg_test/base"
	"tg_test/bot"
	"tg_test/dto/responses"
)

func Serve() {
	router := gin.Default()

	router.GET("/messages", getMessage)
	router.POST("/message", sendMessage)

	router.GET("/files", getFile)

	panic(router.Run(":80"))
}

func getMessage(c *gin.Context) {
	messages, err := base.GetMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := make([]responses.Message, 0)

	for _, msg := range messages {
		m := responses.Message{
			ID:        msg.ID,
			ChatID:    msg.ChatID,
			Text:      msg.Text,
			Creator:   string(msg.Creator),
			CreatedAt: msg.CreatedAt,
		}

		if msg.FileID != nil {
			fileUrl, err := bot.GetFileUrl(*msg.FileID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
			fileUrl = "/files?file=" + url.PathEscape(fileUrl)

			m.File = &fileUrl
		}

		resp = append(resp, m)
	}

	c.JSON(http.StatusOK, resp)
}

func sendMessage(c *gin.Context) {
	var (
		chatID, err1 = strconv.Atoi(c.PostForm("chat_id"))
		text         = c.PostForm("text")
		file, err2   = c.FormFile("file")
	)

	if err1 != nil {
		c.JSON(http.StatusBadRequest, err1)
		return
	}

	var tgMsg tg.Chattable

	if err2 == http.ErrMissingFile {
		// message without file
		tgMsg = tg.NewMessage(int64(chatID), text)

	} else if err2 != nil {
		// another error
		c.JSON(http.StatusInternalServerError, err2)
		return

	} else {
		// file exist
		fl, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		//goland:noinspection GoUnhandledErrorResult
		defer fl.Close()

		doc := tg.NewDocument(int64(chatID), tg.FileReader{Name: file.Filename, Reader: fl})
		doc.Caption = text

		tgMsg = doc

	}

	msg, err := bot.Send(tgMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	dbMsg := base.Message{
		ID:      int64(msg.MessageID),
		ChatID:  msg.Chat.ID,
		Text:    msg.Text,
		Creator: base.MessageCreatorOperator,
	}
	if msg.Document != nil {
		dbMsg.FileID = &msg.Document.FileID
	}

	err = base.SaveMessage(&dbMsg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func getFile(c *gin.Context) {
	link := fmt.Sprintf(tg.FileEndpoint, os.Getenv("BOT_TOKEN"), c.Query("file"))
	resp, err := http.Get(link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "cannot get image") // can be token in the error
		return
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	//goland:noinspection GoUnhandledErrorResult
	io.Copy(c.Writer, resp.Body)

	for k, v := range resp.Header {
		c.Header(k, v[0])
	}
}