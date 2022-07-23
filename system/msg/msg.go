package message

import (
	"fmt"
	"os/exec"
	"strings"
        "context"
	
	"go.amirul.dev/system/help"
	"go.amirul.dev/system/lib"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

// Config
var (
	prefix = "."
	self   = false
	owner  = "687852104"
)

func Msg(client *whatsmeow.Client, msg *events.Message) {
	// simple
	zx := lib.NewSimp(client, msg)
	// dll
	from := msg.Info.Chat
	sender := msg.Info.Sender.String()
	pushName := msg.Info.PushName
	isOwner := strings.Contains(sender, owner)
	isAdmin := zx.GetGroupAdmin(from, sender)
	isBotAdm := zx.GetGroupAdmin(from, "6285742431407@s.whatsapp.net")
	isGroup := msg.Info.IsGroup
	args := strings.Split(zx.GetCMD(), " ")
	command := strings.ToLower(args[0])
	query := strings.Join(args[1:], ` `)
	extended := msg.Message.GetExtendedTextMessage()
	quotedMsg := extended.GetContextInfo().GetQuotedMessage()
	quotedImage := quotedMsg.GetImageMessage()
	//quotedVideo := quotedMsg.GetVideoMessage()
	//quotedSticker := quotedMsg.GetStickerMessage()
	// Self

	if self && !isOwner {
		return
	}

	//-- CONSOLE LOG
	fmt.Println("\n===============================\nNAME: " + pushName + "\nJID: " + sender + "\nTYPE: " + msg.Info.Type + "\nMessage: " + command + ".")
	switch command {
	case "bot":
		zx.Reply(`Bot aktif *` + pushName + `*`)
	case ">":
		if !isOwner {
			zx.Reply(helper.Own())
			return
		}

	case "$":
		if !isOwner {
			zx.Reply(helper.Own())
			return
		}
		out, err := exec.Command("bash", "-c", query).Output()
		if err != nil {
			zx.Reply(fmt.Sprintf("%v", err))
			return
		}
		zx.Reply(string(out))

	case prefix + "menu":
		buttons := []*waProto.HydratedTemplateButton{
			{
				HydratedButton: &waProto.HydratedTemplateButton_QuickReplyButton{
					QuickReplyButton: &waProto.HydratedQuickReplyButton{
						DisplayText: proto.String("OWNER"),
						Id:          proto.String(prefix + "owner"),
					},
				},
			},
		}
		zx.Hydrated(from, helper.Menu(pushName, prefix), "Library : Whatsmeow", buttons)
case prefix + "sticker":
		if quotedImage != nil {
			data, _ := client.Download(quotedImage)
			stc := zx.CreateStickerIMG(data)
			client.SendMessage(context.Background(), from, "", stc)
		} else if msg.Message.GetImageMessage() != nil {
			data, _ := client.Download(msg.Message.GetImageMessage())
			stc := zx.CreateStickerIMG(data)
			client.SendMessage(context.Background(), from, "", stc)
		}
		
	//-- get link gc
	case prefix + "linkgc":
		if !isGroup {
			zx.Reply(helper.Gc())
			return
		}
		if !isBotAdm {
			zx.Reply(helper.BotAdm())
			return
		}
		if !isAdmin {
			zx.Reply(helper.Adm())
			return
		}

		link := zx.LinkGc(from, false)
		zx.Reply(`Berikut link group: ` + link + ``)

	case prefix + "owner":
		zx.SendContact(from, owner, "Amirul Dev")

	}
}
