package server

import (
	"github.com/gin-gonic/gin"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	fl "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"tg_test/base"
	"tg_test/bot"
	"tg_test/dto/responses"
)

func init() {
	file, _ := os.ReadFile("./docs/swagger.json")
	swag.Register("swagger", &swag.Spec{SwaggerTemplate: string(file)})
}

func Serve() {
	router := gin.Default()

	router.GET("/messages", getMessages)
	router.POST("/message", sendMessage)

	router.GET("/files", getFile)

	router.GET("/swagger/*any", gs.WrapHandler(fl.Handler))

	panic(router.Run(":80"))
}

// @Summary GetMessages
// @Produce json
// @Tags Message
// @Success 200 {array} responses.Message
// @Router /messages [get]
func getMessages(c *gin.Context) {
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
			fileUrl := "/files?file=" + *msg.FileID
			m.File = &fileUrl
		}

		resp = append(resp, m)
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary SendMessages
// @Produce json
// @Success 204
// @Tags Message
// @Param chat_id formData int    true  "ChatID"
// @Param text    formData string false "Text"
// @Param file    formData file   false "File"
// @Router /message [post]
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
		mp, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		//goland:noinspection GoUnhandledErrorResult
		defer mp.Close()

		doc := tg.NewDocument(int64(chatID), tg.FileReader{Name: file.Filename, Reader: mp})
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
	fileUrl, err := bot.GetFileDirectURL(c.Query("file"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp, err := http.Get(fileUrl)
	if err != nil {
		log.Println(err)
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
