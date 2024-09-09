package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

type OpenAIRequest struct {
    Model    string  `json:"model"`
    Messages []Message `json:"messages"`
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type OpenAIResponse struct {
    Choices []struct {
        Message Message `json:"message"`
    } `json:"choices"`
}

func (p *ChatbotPlugin) callOpenAI(message string) (string, error) {
    apiKey := p.getConfiguration().OpenAIKey

    messages := []Message{
        {Role: "system", Content: "You are a helpful assistant."},
        {Role: "user", Content: message},
    }

    requestBody, err := json.Marshal(OpenAIRequest{
        Model:    "gpt-3.5-turbo",
        Messages: messages,
    })
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(requestBody))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var openAIResp OpenAIResponse
    if err := json.Unmarshal(body, &openAIResp); err != nil {
        return "", err
    }

    if len(openAIResp.Choices) > 0 {
        return openAIResp.Choices[0].Message.Content, nil
    }

    return "Sorry, I couldn't generate a response.", nil
}
