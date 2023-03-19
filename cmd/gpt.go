package cmd

import (
	"Te4g/ai/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var gptCmd = &cobra.Command{
	Use:   "gpt",
	Short: "Send prompt to GPT model and get response",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No prompt provided")
			return
		}

		messages := []Message{{Role: "user", Content: strings.Join(args, " ")}}

		requestBody := GptRequestBody{
			Model:    "gpt-3.5-turbo",
			Messages: messages,
		}

		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_ACCESS_TOKEN")))
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var data GptResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err)
		}

		fmt.Println(utils.RemoveFirstCarriageReturn(data.Choices[0].Message.Content))
	},
}

func init() {
	rootCmd.AddCommand(gptCmd)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptRequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type GptResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
