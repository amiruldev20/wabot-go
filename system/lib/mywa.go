package lib

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type renz struct {
	RClient *whatsmeow.Client
	Msg     *events.Message
}

func NewSimp(Cli *whatsmeow.Client, m *events.Message) *renz {
	return &renz{
		RClient: Cli,
		Msg:     m,
	}
}

//-- Reply Message
func (zx *renz) Reply(teks string) {
	zx.RClient.SendMessage(context.Background(), zx.Msg.Info.Chat, "", &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(teks),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &zx.Msg.Info.ID,
				Participant:   proto.String(zx.Msg.Info.Sender.String()),
				QuotedMessage: zx.Msg.Message,
			},
		},
	})
}

func (simp *SimpleImpl) CreateStickerIMG(data []byte) *waProto.Message {
	RawPath := fmt.Sprintf("tmp/%s%s", simp.Msg.Info.ID, ".jpg")
	ConvertedPath := fmt.Sprintf("tmp/sticker/%s%s", simp.Msg.Info.ID, ".webp")
	err := os.WriteFile(RawPath, data, 0600)
	if err != nil {
		fmt.Printf("Failed to save image: %v", err)
	}
	exc := exec.Command("cwebp", "-q", "80", RawPath, "-o", ConvertedPath)
	err = exc.Run()
	if err != nil {
		fmt.Println("Failed to Convert Image to WebP")
	}
	createExif := fmt.Sprintf("webpmux -set exif %s %s -o %s", "tmp/exif/raw.exif", ConvertedPath, ConvertedPath)
	cmd := exec.Command("bash", "-c", createExif)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Failed to set webp metadata", err)
	}
	stc, err := os.ReadFile(ConvertedPath)
	if err != nil {
		fmt.Printf("Failed to read %s: %s\n", ConvertedPath, err)
	}
	uploaded, err := simp.VClient.Upload(context.Background(), stc, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
	}
	return &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(stc)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &simp.Msg.Info.ID,
				Participant:   proto.String(simp.Msg.Info.Sender.String()),
				QuotedMessage: simp.Msg.Message,
			},
		},
	}
}

func (zx *renz) Hydrated(jid types.JID, teks string, foter string, buttons []*waProto.HydratedTemplateButton) {
	zx.RClient.SendMessage(context.Background(), jid, "", &waProto.Message{
		TemplateMessage: &waProto.TemplateMessage{
			HydratedTemplate: &waProto.HydratedFourRowTemplate{
				HydratedContentText: proto.String(teks),
				HydratedFooterText:  proto.String(foter),
				HydratedButtons:     buttons,
			},
		},
	})
}

func (zx *renz) SendContact(jid types.JID, number string, nama string) {
	zx.RClient.SendMessage(context.Background(), jid, "", &waProto.Message{
		ContactMessage: &waProto.ContactMessage{
			DisplayName: proto.String(nama),
			Vcard:       proto.String(fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN:%s;;;\nFN:%s\nitem1.TEL;waid=%s:+%s\nitem1.X-ABLabel:Mobile\nEND:VCARD", nama, nama, number, number)),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &zx.Msg.Info.ID,
				Participant:   proto.String(zx.Msg.Info.Sender.String()),
				QuotedMessage: zx.Msg.Message,
			},
		},
	})
}

//-- SEND STICKER
/*
func (zx *renz) Sticker(Jid types.JID)
*/

//-- GET LINK GC
func (zx *renz) LinkGc(Jid types.JID, reset bool) string {
	link, err := zx.RClient.GetGroupInviteLink(Jid, reset)

	if err != nil {
		panic(err)
	}
	return link
}

func (zx *renz) FetchGroupAdmin(Jid types.JID) ([]string, error) {
	var Admin []string
	resp, err := zx.RClient.GetGroupInfo(Jid)
	if err != nil {
		return Admin, err
	} else {
		for _, group := range resp.Participants {
			if group.IsAdmin || group.IsSuperAdmin {
				Admin = append(Admin, group.JID.String())
			}
		}
	}
	return Admin, nil
}

func (zx *renz) GetGroupAdmin(jid types.JID, sender string) bool {
	if !zx.Msg.Info.IsGroup {
		return false
	}
	admin, err := zx.FetchGroupAdmin(jid)
	if err != nil {
		return false
	}
	for _, v := range admin {
		if v == sender {
			return true
		}
	}
	return false
}

func (zx *renz) GetCMD() string {
	extended := zx.Msg.Message.GetExtendedTextMessage().GetText()
	text := zx.Msg.Message.GetConversation()
	imageMatch := zx.Msg.Message.GetImageMessage().GetCaption()
	videoMatch := zx.Msg.Message.GetVideoMessage().GetCaption()
	tempBtnId := zx.Msg.Message.GetTemplateButtonReplyMessage().GetSelectedId()
	btnId := zx.Msg.Message.GetButtonsResponseMessage().GetSelectedButtonId()
	listId := zx.Msg.Message.GetListResponseMessage().GetSingleSelectReply().GetSelectedRowId()
	var command string
	if text != "" {
		command = text
	} else if imageMatch != "" {
		command = imageMatch
	} else if videoMatch != "" {
		command = videoMatch
	} else if extended != "" {
		command = extended
	} else if tempBtnId != "" {
		command = tempBtnId
	} else if btnId != "" {
		command = btnId
	} else if listId != "" {
		command = listId
	}
	return command
}
