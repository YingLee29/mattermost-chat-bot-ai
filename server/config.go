package main

import (
    "sync"
)

type Configuration struct {
    OpenAIKey string
}

type ChatbotPlugin struct {
    configurationLock sync.RWMutex
    configuration     *Configuration
}

func (p *ChatbotPlugin) getConfiguration() *Configuration {
    p.configurationLock.RLock()
    defer p.configurationLock.RUnlock()

    if p.configuration == nil {
        return &Configuration{}
    }

    return p.configuration
}

func (p *ChatbotPlugin) OnConfigurationChange() error {
    var configuration = new(Configuration)

    if err := p.API.LoadPluginConfiguration(configuration); err != nil {
        return err
    }

    p.configurationLock.Lock()
    p.configuration = configuration
    p.configurationLock.Unlock()

    return nil
}
