package main

import (
    "fmt"
    "github.com/mattermost/mattermost-server/v6/model"
    "github.com/mattermost/mattermost-server/v6/plugin"
)

type ChatbotPlugin struct {
    plugin.MattermostPlugin
    botUserID string
}

func (p *ChatbotPlugin) OnActivate() error {
    botID, err := p.Helpers.EnsureBot(&model.Bot{
        Username:    "chatbotgpt",
        DisplayName: "Chatbot GPT",
        Description: "A chatbot powered by OpenAI GPT.",
    })
    if err != nil {
        return err
    }
    p.botUserID = botID
    return nil
}

func (p *ChatbotPlugin) MessageHasBeenPosted(_ *plugin.Context, post *model.Post) {
    if post.UserId == p.botUserID {
        return
    }

    if !p.isMentioned(post.Message) {
        return
    }

    response, err := p.callOpenAI(post.Message)
    if err != nil {
        p.API.LogError("Failed to call OpenAI", "error", err.Error())
        return
    }

    newPost := &model.Post{
        UserId:    p.botUserID,
        ChannelId: post.ChannelId,
        Message:   response,
        RootId:    post.Id,
    }
    p.API.CreatePost(newPost)
}

func (p *ChatbotPlugin) isMentioned(message string) bool {
    return model.IsStringInStringsIgnoreCase([]string{"@chatbotgpt"}, message)
}
