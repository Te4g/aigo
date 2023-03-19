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
)

var davinciCmd = &cobra.Command{
	Use:   "davinci",
	Short: "Send prompt to DaVinci model and get response",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No prompt provided")
			return
		}

		requestBody := DavinciRequestBody{
			Model:       "text-davinci-003",
			Prompt:      strings.Join(args, " "),
			MaxTokens:   512,
			Temperature: 0.1,
		}

		requestBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBodyBytes))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_ACCESS_TOKEN")))
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var data DavinciResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err)
		}

		fmt.Println(utils.RemoveFirstCarriageReturn(data.Choices[0].Text))
	},
}

func init() {
	rootCmd.AddCommand(davinciCmd)
}

type DavinciRequestBody struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
}

type DavinciResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
