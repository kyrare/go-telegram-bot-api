package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const apiEndpoint = "https://api.telegram.org/bot%s/%s"

type Bot struct {
	Token         string
	User          User
	events        []Event
	updatesOffset int
}

func New(token string) (Bot, error) {
	bot := Bot{Token: token}

	user, err := bot.GetMe()

	if err != nil {
		return Bot{}, err
	}

	bot.User = user

	return bot, nil
}

func (bot *Bot) request(endpoint string, params RequestParams) (ApiResponse, error) {
	requestUrl := fmt.Sprintf(apiEndpoint, bot.Token, endpoint)

	values := buildRequestParams(params)

	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(values.Encode()))

	if err != nil {
		panic(err)
		return ApiResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
		return ApiResponse{}, nil
	}

	defer resp.Body.Close()

	//bytes, err := io.ReadAll(resp.Body)
	//fmt.Println("response body", string(bytes))

	var apiResponse ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)

	return apiResponse, err
}

func buildRequestParams(params RequestParams) url.Values {
	if params == nil {
		return url.Values{}
	}

	values := url.Values{}

	for name, value := range params {
		values.Set(name, value)
	}

	return values
}

func (bot *Bot) GetMe() (User, error) {
	resp, err := bot.request("getMe", nil)

	if err != nil {
		panic(err)
		//return User{}, nil
	}

	var user User
	err = json.Unmarshal(resp.Result, &user)

	return user, err
}

func (bot *Bot) GetUpdates(offset int) ([]Update, error) {
	params := RequestParams{}

	if offset > 0 {
		params["offset"] = strconv.Itoa(offset)
	}

	resp, err := bot.request("getUpdates", params)

	if err != nil {
		return []Update{}, err
	}

	var updates []Update
	err = json.Unmarshal(resp.Result, &updates)

	return updates, err
}

func (bot *Bot) SendMessage(chatId int, message string) (bool, error) {
	params := make(map[string]string)

	params["chat_id"] = strconv.Itoa(chatId)
	params["text"] = message

	resp, err := bot.request("sendMessage", params)

	return resp.Ok, err
}

func (bot *Bot) Run() {
	for {
		updates, _ := bot.GetUpdates(bot.updatesOffset)

		for _, update := range updates {
			bot.updatesOffset = update.UpdateID + 1

			for _, event := range bot.events {
				if event.ExecuteChecker(update) {
					go event.ExecuteAction(update)
				}
			}
		}

		time.Sleep(time.Second * 2)
	}
}

func (bot *Bot) Command(command string, action func(message Message)) *Bot {
	bot.On(
		func(update Update) bool {
			return update.Message.isBotCommand() && strings.HasPrefix(update.Message.Text, "/"+command)
		},
		func(update Update) {
			action(update.Message)
		},
	)

	return bot
}

func (bot *Bot) On(checker func(update Update) bool, action func(update Update)) *Bot {
	bot.events = append(bot.events, Event{
		checker: checker,
		action:  action,
	})

	return bot
}
