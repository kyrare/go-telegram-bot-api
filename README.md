Велосипед для работы с Telegram Bot Api для изучения go (pet project)

Вдохновлено https://github.com/TelegramBot/Api

Пример использования

```go
bot, err := botApi.New(secret)

if err != nil {
    panic(err)
}

bot.Command("ping", func(message botApi.Message) {
    bot.SendMessage(message.Chat.Id, "pong")
}).On(
    func(update botApi.Update) bool {
        return strings.ToLower(update.Message.Text) == "привет"
    },
    func(update botApi.Update) {
        bot.SendMessage(update.Message.Chat.Id, "И тебе привет человек")
    },
)
```