package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
}

// // ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
// func (p *HelloWorldPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Hello, world!")
// }

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
		AutoCompleteDesc: "Reply to a message (Use it by clicking reply first then slash command)",
		AutoCompleteHint: "Let's talk here",
	})
}

// ExecuteCommand to send message
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	rawText := strings.TrimSpace((strings.Replace(args.Command, "/"+trigger, "", 1)))
	userID := args.UserId
	rootID := args.RootId

	post, err := p.API.GetPost(rootID)
	if err != nil {
		return nil, err
	}

	otherUserID := post.UserId

	channel, err1 := p.API.GetDirectChannel(userID, otherUserID)
	if err1 != nil {
		return nil, err1
	}

	channelID := channel.Id

	postModel := &model.Post{
		UserId:    userID,
		ChannelId: channelID,
		// RootId:    rootID,
		Message: rawText,
	}

	_, err2 := p.API.CreatePost(postModel)
	if err2 != nil {
		return nil, err2
	}

	return &model.CommandResponse{}, nil
}

func main() {
	plugin.ClientMain(&Plugin{})
}
