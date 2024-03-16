package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	wolfram "github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events:")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	witClient := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}

	go printCommandEvents(bot.CommandEvents())

	bot.Command("query <message>", &slacker.CommandDefinition{
		Description: "Send your question to Wolfram",
		Examples:    []string{"What's the meaning of life?", "What am I doing here?"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")
			// fmt.Println("Your query: " + query)
			msg, _ := witClient.Parse(&witai.MessageRequest{
				Query: query,
			})
			// fmt.Println(msg)
			data, _ := json.MarshalIndent(msg, "", "    ")
			rough := string(data[:])
			fmt.Println(rough)
			value := gjson.Get(rough, "entities.wolfram_search_query.0.value")
			fmt.Println(value)
			question := value.String()
			res, err := wolframClient.GetSpokentAnswerQuery(question, wolfram.Metric, 1000)
			if err != nil {
				fmt.Println(err)
			}

			response.Reply("your answer is: " + res)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
