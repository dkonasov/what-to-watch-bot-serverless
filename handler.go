package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(res http.ResponseWriter, req *http.Request) {
	var update tgbotapi.Update
	data, err := io.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occured while parsing request body: %s", err.Error())
	} else {
		err = json.Unmarshal(data, &update)

		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while parsing update: %s", err.Error())
		} else {
			bot := Bot{os.Getenv("TOKEN")}
			bot.HandleRequest(&update)
		}
	}
	res.WriteHeader(204)
}
