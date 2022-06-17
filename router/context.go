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
	"errors"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	waproto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// Context is the default context for the router.
type Context struct {
	Cli         *whatsmeow.Client
	IsPrivate   bool
	Msg         *events.Message
	GroupInfo   *events.GroupInfo
	JoinedGroup *events.JoinedGroup
	Args        []string
	Str         string
	Routes      []*Route

	// vars for the router
	vars map[string]any
	Mut  sync.Mutex
}

// Set sets a variable to the context.
func (c *Context) Set(key string, value any) {
	c.Mut.Lock()
	c.vars[key] = value
	c.Mut.Unlock()
}

// Get gets a variable from the context.
func (c *Context) Get(key string) any {
	c.Mut.Lock()
	defer c.Mut.Unlock()
	return c.vars[key]
}

// MarkAsRead marks the current message as read.
func (c *Context) MarkAsRead() error {
	return c.Cli.MarkRead([]string{c.Msg.Info.ID}, time.Now(), c.Msg.Info.Chat, c.Msg.Info.Sender)
}

// LeaveGroup leaves the group.
func (c *Context) LeaveGroup(jid types.JID) error {
	return c.Cli.LeaveGroup(jid)
}

// JoinGroup joins the group by the given link.
func (c *Context) JoinGroup(link string) error {
	_, err := c.Cli.JoinGroupWithLink(link)
	return err
}

// GetGroupInfo gets the group info.
func (c *Context) GetGroupInfo(jid types.JID) (*types.GroupInfo, error) {
	return c.Cli.GetGroupInfo(jid)
}

// DeleteMessage deletes the message.
func (c *Context) DeleteMessage(chat types.JID, message types.MessageID) error {
	_, err := c.Cli.RevokeMessage(chat, message)
	return err
}

// DeleteBulkMessages deletes the messages.
func (c *Context) DeleteBulkMessages(chat types.JID, messages []types.MessageID) error {
	for _, v := range messages {
		err := c.DeleteMessage(chat, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// sendAmsg sends a proto message to specific JID.
func (c *Context) sendAmsg(to types.JID, message *waproto.Message) error {
	// prevent send message to self
	if c.Cli.Store.ID.String() == to.String() {
		return errors.New("cannot send message to self")
	}
	// prevent send message to status@broadcast
	if to.String() == "status@broadcast" {
		return errors.New("cannot send message to status@broadcast")
	}

	_, err := c.Cli.SendMessage(to, "", message)
	return err
}

// SendMessageWithMention sends a message to a specific JID with a mention.
func (c *Context) SendMessageWithMention(to types.JID, message string, mention []string) error {
	return c.sendAmsg(to, &waproto.Message{
		ExtendedTextMessage: &waproto.ExtendedTextMessage{
			Text: &message,
			ContextInfo: &waproto.ContextInfo{
				MentionedJid: mention,
			},
		},
	})
}

// SendMessage sends a basic string message to a specific JID.
func (c *Context) SendMessage(to types.JID, message string) error {
	return c.sendAmsg(to, &waproto.Message{
		Conversation: &message,
	})
}

// Reply reply to the message.
func (c *Context) Reply(message string) error {
	return c.sendAmsg(c.Msg.Info.Chat, &waproto.Message{
		ExtendedTextMessage: &waproto.ExtendedTextMessage{
			Text: &message,
			ContextInfo: &waproto.ContextInfo{
				StanzaId:      &c.Msg.Info.ID,
				Participant:   proto.String(c.Msg.Info.Sender.String()),
				QuotedMessage: c.Msg.Message,
			},
		},
	})
}
