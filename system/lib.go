/*
###################################
# Name: Mywa BOT                  #
# Version: 1.0.1                  #
# Developer: Amirul Dev           #
# Library: waSocket               #
# Contact: 085157489446           #
###################################
# Thanks to: 
# Vnia
*/
package system

import (
    "context"
    "fmt"
    "strings"
    //"time"
    "net/http"
    "os"
    "os/exec"
    //"image/jpeg"
    //"image/webp"
    "github.com/amiruldev20/waSocket"
    waProto "github.com/amiruldev20/waSocket/binary/proto"
    "github.com/amiruldev20/waSocket/types"
    "github.com/amiruldev20/waSocket/types/events"
    "google.golang.org/protobuf/proto"
)

type renz struct {
    RClient * waSocket.Client
    Msg * events.Message
}

func NewSimp(Cli * waSocket.Client, m * events.Message) * renz {
    return &renz {
        RClient: Cli,
        Msg: m,
    }
}


/* pair code */
func(sock * renz) loginPhone(jid string)(string, int, error) {
    // Request Pairing Code
    code, err: = sock.RClient.PairPhone(jid, true, waSocket.PairClientChrome, "Chrome (Linux)")
    if err != nil {
        return "", 0, err
    }
    return code, 160, nil
}

/* parse jid */
func(sock * renz) parseJID(arg string)(types.JID, bool) {
    if arg[0] == '+' {
        arg = arg[1: ]
    }
    if !strings.ContainsRune(arg, '@') {
        return types.NewJID(arg, types.DefaultUserServer), true
    } else {
        recipient, err: = types.ParseJID(arg)
        if err != nil {
            fmt.Println("Invalid JID %s: %v", arg, err)
            return recipient, false
        } else if recipient.User == "" {
            fmt.Println("Invalid JID %s: no server specified", arg)
            return recipient, false
        }
        return recipient, true
    }
}

/* send react */
func(sock * renz) React(react string) {
    _, err: = sock.RClient.SendMessage(context.Background(), sock.Msg.Info.Chat, sock.RClient.BuildReaction(sock.Msg.Info.Chat, sock.Msg.Info.Sender, sock.Msg.Info.ID, react))
    if err != nil {
        return
    }
}

/* send message */
func(sock * renz) SendMsg(jid types.JID, teks string) {
    _, err: = sock.RClient.SendMessage(context.Background(), jid, & waProto.Message {
        Conversation: proto.String(teks)
    })
    if err != nil {
        return
    }
}

/* send reply */
func(sock * renz) Reply(teks string) {
    _, err: = sock.RClient.SendMessage(context.Background(), sock.Msg.Info.Chat, & waProto.Message {
        ExtendedTextMessage: & waProto.ExtendedTextMessage {
            Text: proto.String(teks),
            ContextInfo: & waProto.ContextInfo {
                Expiration: proto.Uint32(86400),
                StanzaId: & sock.Msg.Info.ID,
                Participant: proto.String(sock.Msg.Info.Sender.String()),
                QuotedMessage: sock.Msg.Message,
            },
        },
    })
    if err != nil {
        return
    }
}

/* send reply */
func(sock * renz) ReplyAd(teks string) {
    _, err: = sock.RClient.SendMessage(context.Background(), sock.Msg.Info.Chat, & waProto.Message {
        ExtendedTextMessage: & waProto.ExtendedTextMessage {
            Text: proto.String(teks),
            ContextInfo: & waProto.ContextInfo {
                ExternalAdReply: & waProto.ContextInfo_ExternalAdReplyInfo {
                    Title: proto.String("MywaBOT 2023"),
                    Body: proto.String("Lightweight whatsapp bot 2023"),
                    //MediaType: waProto.ContextInfo_ExternalAdReplyInfo_MediaType(proto.Number(1)),
                    ThumbnailUrl: proto.String("https://telegra.ph/file/eb7261ee8de82f8f48142.jpg"),
                    RenderLargerThumbnail: proto.Bool(true),
                },
                Expiration: proto.Uint32(86400),
                StanzaId: & sock.Msg.Info.ID,
                Participant: proto.String(sock.Msg.Info.Sender.String()),
                QuotedMessage: sock.Msg.Message,
            },
        },
    })
    if err != nil {
        return
    }
}

/* send contact */
func(sock * renz) SendContact(jid types.JID, number string, nama string) {
    _, err: = sock.RClient.SendMessage(context.Background(), jid, & waProto.Message {
        ContactMessage: & waProto.ContactMessage {
            DisplayName: proto.String(nama),
            Vcard: proto.String(fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN:%s;;;\nFN:%s\nitem1.TEL;waid=%s:+%s\nitem1.X-ABLabel:Mobile\nEND:VCARD", nama, nama, number, number)),
            ContextInfo: & waProto.ContextInfo {
                StanzaId: & sock.Msg.Info.ID,
                Participant: proto.String(sock.Msg.Info.Sender.String()),
                QuotedMessage: sock.Msg.Message,
            },
        },
    })
    if err != nil {
        return
    }
}

/* send image 
func (sock *renz) SendImage(jid types.JID, number string, nama string) {
_, err := sock.RClient.SendMessage(context.Background(), jid, &waProto.Message{
ImageMessage: &waProto.ImageMessage{
JpegThumbnail: dataWaThumbnail,
Caption:       proto.String(dataWaCaption),
Url:           proto.String(uploadedImage.URL),
DirectPath:    proto.String(uploadedImage.DirectPath),
MediaKey:      uploadedImage.MediaKey,
Mimetype:      proto.String(http.DetectContentType(dataWaImage)),
FileEncSha256: uploadedImage.FileEncSHA256,
FileSha256:    uploadedImage.FileSHA256,
FileLength:    proto.Uint64(uint64(len(dataWaImage))),
ViewOnce:      proto.Bool(request.ViewOnce),
ContextInfo: &waProto.ContextInfo{
StanzaId:&sock.Msg.Info.ID,
Participant: proto.String(sock.Msg.Info.Sender.String()),
QuotedMessage: sock.Msg.Message,
},
},
})
if err != nil {
return
}
}
/* create sticker */
func(sock * renz) CreateSticker(data[] byte) * waProto.Message {
    RawPath: = fmt.Sprintf("tmp/%s%s", sock.Msg.Info.ID, ".jpg")
    ConvertedPath: = fmt.Sprintf("tmp/sticker/%s%s", sock.Msg.Info.ID, ".webp")
    err: = os.WriteFile(RawPath, data, 0600)
    if err != nil {
        fmt.Printf("Failed to save image: %v", err)
    }
    exc: = exec.Command("cwebp", "-q", "80", RawPath, "-o", ConvertedPath)
    err = exc.Run()
    if err != nil {
        fmt.Println("Failed to Convert Image to WebP")
    }
    createExif: = fmt.Sprintf("webpmux -set exif %s %s -o %s", "tmp/exif/raw.exif", ConvertedPath, ConvertedPath)
    cmd: = exec.Command("bash", "-c", createExif)
    err = cmd.Run()
    if err != nil {
        fmt.Println("Failed to set webp metadata", err)
    }
    stc,
    err: = os.ReadFile(ConvertedPath)
    if err != nil {
        fmt.Printf("Failed to read %s: %s\n", ConvertedPath, err)
    }
    uploaded,
    err: = sock.RClient.Upload(context.Background(), stc, waSocket.MediaImage)
    if err != nil {
        fmt.Printf("Failed to upload file: %v\n", err)
    }
    return &waProto.Message {
        StickerMessage: & waProto.StickerMessage {
            Url: proto.String(uploaded.URL),
            DirectPath: proto.String(uploaded.DirectPath),
            MediaKey: uploaded.MediaKey,
            Mimetype: proto.String(http.DetectContentType(stc)),
            FileEncSha256: uploaded.FileEncSHA256,
            FileSha256: uploaded.FileSHA256,
            FileLength: proto.Uint64(uint64(len(data))),
            ContextInfo: & waProto.ContextInfo {
                StanzaId: & sock.Msg.Info.ID,
                Participant: proto.String(sock.Msg.Info.Sender.String()),
                QuotedMessage: sock.Msg.Message,
            },
        },
    }

}

func(sock * renz) createChannel(params[] string) {
    _, err: = sock.RClient.CreateNewsletter(waSocket.CreateNewsletterParams {
        Name: strings.Join(params, " "),
    })
    if err != nil {
        return
    }
}

func(sock * renz) FetchGroupAdmin(Jid types.JID)([] string, error) {
    var Admin[] string
    resp, err: = sock.RClient.GetGroupInfo(Jid)
    if err != nil {
        return Admin, err
    } else {
        for _, group: = range resp.Participants {
            if group.IsAdmin || group.IsSuperAdmin {
                Admin = append(Admin, group.JID.String())
            }
        }
    }
    return Admin, nil
}

func(sock * renz) GetGroupAdmin(jid types.JID, sender string) bool {
        if !sock.Msg.Info.IsGroup {
            return false
        }
        admin, err: = sock.FetchGroupAdmin(jid)
        if err != nil {
            return false
        }
        for _, v: = range admin {
            if v == sender {
                return true
            }
        }
        return false
    }
    //-- GET LINK GC
func(sock * renz) LinkGc(Jid types.JID, reset bool) string {
    link, err: = sock.RClient.GetGroupInviteLink(Jid, reset)

    if err != nil {
        panic(err)
    }
    return link
}

/* update participants 
func (sock *renz) UpdateParticipants(jid types.JID, participantChanges map[types.JID]ParticipantChange){
_, err := sock.RClient.UpdateGroupParticipants(sock.Msg.Info.Chat, map[types.JID]waSocket.ParticipantChange{
jid: waSocket.ParticipantChangeRemove
})
if err != nil {
return
}
}
*/
/*
func (sock *renz) userInfo(Jid []types.JID) map {
resp, err := sock.RClient.GetUserInfo(Jid)

if err != nil {
panic(err)
}
return resp
}
}
}
*/


func(sock * renz) GetCMD() string {
    extended: = sock.Msg.Message.GetExtendedTextMessage().GetText()
    text: = sock.Msg.Message.GetConversation()
    imageMatch: = sock.Msg.Message.GetImageMessage().GetCaption()
    videoMatch: = sock.Msg.Message.GetVideoMessage().GetCaption()
    tempBtnId: = sock.Msg.Message.GetTemplateButtonReplyMessage().GetSelectedId()
    btnId: = sock.Msg.Message.GetButtonsResponseMessage().GetSelectedButtonId()
    listId: = sock.Msg.Message.GetListResponseMessage().GetSingleSelectReply().GetSelectedRowId()
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