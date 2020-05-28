package main

import (
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin that adds a slash command to reply privately
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *Configuration
}

const (
	trigger = "private"
)

// OnActivate register the plugin command
func (p *Plugin) OnActivate() error {
	return p.API.RegisterCommand(&model.Command{
		Trigger:          trigger,
		Description:      "Reply privately to message",
		DisplayName:      "Reply Privately",
		AutoComplete:     true,
		AutoCompleteDesc: "Hide a spoiler message (Use it by clicking reply first then slash command)",
		AutoCompleteHint: "Let's talk here",
	})
}

// ExecuteCommand ro send message
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	rawText := strings.TrimSpace((strings.Replace(args.Command, "/"+trigger, "", 1)))
	userID := args.UserId
	rootID := args.RootId

	post, err := p.API.GetPost(rootID)
	if err != nil {
		return nil, err
	}

	otherUserID := post.UserId

	channel, err1 := GetGroupChannel([]string{userID, otherUserID})
	if err1 != nil {
		return nil, err1
	}

	channelID := channel.UserId

	postModel := &model.Post{
		UserId:    userID,
		ChannelId: channelID,
		RootId:    rootID,
		Message:   rawText,
	}

	_, err2 := p.API.CreatePost(postModel)
	if err2 != nil {
		return nil, err2
	}

	return &model.CommandResponse{}, nil
}
