// Copyright 2022 Ade M Ramdani <ramdanhere04@gmail.com>
// This file is part of whatsrhyno
//
// whatsrhyno is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// whatsrhyno is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with whatsrhyno.  If not, see <http://www.gnu.org/licenses/>.

package router

import (
	"strings"

	"go.mau.fi/whatsmeow/types/events"
	"go.ramdanhere.dev/whatsrhyno"
)

const (
	groupChange = "groupChange"
	joinedGroup = "joinedGroup"
)

var (
	globalMessageHandler HandlerFunc
	notFoundHandler      HandlerFunc
	eventsHandler        = make(map[string]func(*Context))

	// hook
	beforeMessageHandler HandlerFunc
	afterMessageHandler  HandlerFunc
)

// SetBeforeMessageHandler sets the function that will be called before the message is sent.
func (r *Route) SetBeforeMessageHandler(handler HandlerFunc) {
	beforeMessageHandler = handler
}

// SetAfterMessageHandler sets the function that will be called after the message is sent.
func (r *Route) SetAfterMessageHandler(handler HandlerFunc) {
	afterMessageHandler = handler
}

// parseEvt parses the event and calls the appropriate handler.
func (r *Route) parseEvt(cli *whatsrhyno.Client, evt interface{}) {
	ctx := &Context{
		Cli:    cli.GetCli(),
		Routes: r.Routes,
	}

	switch evt.(type) {
	case *events.GroupInfo:
		ctx.GroupInfo = evt.(*events.GroupInfo)
		if fn := eventsHandler[groupChange]; fn != nil {
			fn(ctx)
		}
	case *events.JoinedGroup:
		ctx.JoinedGroup = evt.(*events.JoinedGroup)
		if fn := eventsHandler[joinedGroup]; fn != nil {
			fn(ctx)
		}
	}
}

// SetGroupChangedHandler sets the function that will be called when the group info is changed.
func (r *Route) SetGroupChangedHandler(handler func(*Context)) {
	eventsHandler[groupChange] = handler
}

// SetJoinedGroupHandler sets the function that will be called when the
// current user joined a group.
func (r *Route) SetJoinedGroupHandler(handler func(*Context)) {
	eventsHandler[joinedGroup] = handler
}

// SetGlobalMessageHandler sets the function that will be called when the
// message is not a command.
func (r *Route) SetGlobalMessageHandler(handler HandlerFunc) {
	globalMessageHandler = handler
}

// SetNotFoundHandler sets the function that will be called when the command is not found.
func (r *Route) SetNotFoundHandler(handler HandlerFunc) {
	notFoundHandler = handler
}

// eventMessageHandler is the function that will be called when the event witj type message is emited.
func (r *Route) eventMessageHandler(cli *whatsrhyno.Client, msg *events.Message) {
	// str is the message that was sent.
	// sometimes the message string is not inside the Conversation object, so we need to find it from another object.
	var str string

	if msg.Message.Conversation != nil {
		str = msg.Message.GetConversation()
	} else if msg.Message.ExtendedTextMessage != nil {
		str = msg.Message.ExtendedTextMessage.GetText()
	} else if msg.Message.ImageMessage != nil {
		str = msg.Message.ImageMessage.GetCaption()
	}

	ctx := &Context{
		Cli:       cli.GetCli(),
		Msg:       msg,
		Str:       str,
		IsPrivate: !msg.Info.IsGroup,
		Routes:    r.Routes,
	}

	// call before hook.
	if beforeMessageHandler != nil {
		beforeMessageHandler(ctx)
	}

	// call after hook.
	if afterMessageHandler != nil {
		defer afterMessageHandler(ctx)
	}

	// prevent handling message from self.
	if msg.Info.IsFromMe {
		return
	}

	// if the AutoRead is enabled, we need to mark as read the message from group.
	if r.AutoRead {
		ctx.MarkAsRead()
	}

	var commands string
	if strings.HasPrefix(str, cli.GetCommandPrefix()) {
		str = strings.TrimPrefix(str, cli.GetCommandPrefix())
		strSlice := strings.Split(str, " ")
		commands = strSlice[0]
		ctx.Args = strSlice[1:]

		// find the handler and execute it.
		nr, i := r.getHandler(commands)
		if i > 0 {
			nr.Handler(ctx)
			return
		}

		if notFoundHandler != nil {
			notFoundHandler(ctx)
		}

		return
	}

	// if the message is not a command and the global message handler is set, call it.
	if globalMessageHandler != nil {
		globalMessageHandler(ctx)
	}

	return
}

// isCommand checks if the message is a command.
func (r *Route) isCommand(str string, cli *whatsrhyno.Client) bool {
	prefix := cli.GetCommandPrefix()
	if prefix == "" {
		return false
	}

	return strings.HasPrefix(str, prefix)
}

// getHandler returns the handler function for the given command.
// TODO: implement subroutes (sub commands).
func (r *Route) getHandler(command string) (*Route, int) {
	nr := r
	i := 0
	if rt := nr.Find(command); rt != nil {
		nr = rt
		i = 1
	}

	return nr, i
}
