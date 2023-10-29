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
	"mywa-bot/config"
	"mywa-bot/system"
    "context"
    "fmt"
    //"os"
    otto "github.com/robertkrimen/otto"
    "time"
    "os/exec"
    "strings"
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
    // "google.golang.org/protobuf/proto"

    "github.com/amiruldev20/waSocket"
    //waProto "github.com/amiruldev20/waSocket/binary/proto"
    "github.com/amiruldev20/waSocket/types/events"
)

func delay(duration time.Duration) {
    time.Sleep(duration)
}

func Msg(client * waSocket.Client, msg * events.Message) {
    // simple
    sock: = system.NewSimp(client, msg)
        // dll
    from: = msg.Info.Chat
    sender: = msg.Info.Sender.String()
    pushName: = msg.Info.PushName
    isOwner: = strings.Contains(sender, owner)
        //isAdmin := sock.GetGroupAdmin(from, sender)
    isBotAdm: = sock.GetGroupAdmin(from, config.BotNumber + "@s.whatsapp.net")
    isGroup: = msg.Info.IsGroup
    args: = strings.Split(sock.GetCMD(), " ")
    command: = strings.ToLower(args[0])
    query: = strings.Join(args[1: ], ` `)
    extended: = msg.Message.GetExtendedTextMessage()
    quotedMsg: = extended.GetContextInfo().GetQuotedMessage()
    quotedImage: = quotedMsg.GetImageMessage()
        //quotedVideo := quotedMsg.GetVideoMessage()
        //quotedSticker := quotedMsg.GetStickerMessage()
        // Self

    if self && !isOwner {
        return
    }

    //-- CONSOLE LOG
    fmt.Println("\n===============================\nNAME: " + pushName + "\nJID: " + sender + "\nTYPE: " + msg.Info.Type + "\nMessage: " + command + "")
    switch command {

        /* panggil bot */
        case "bot":
            sock.Reply(`Bot aktif *` + pushName + `*`)
            sock.React("🤖")

            /* eval */
        case "*":
            sock.React("⏱️")
            if !isOwner {
                sock.Reply(system.Own())
                return
            }

            vm: = otto.New()
            _, er: = vm.Run(query)
            if er != nil {
                sock.Reply(fmt.Sprintf("%s", er))
                return
            }
            x, err: = vm.Get(query)
            if err != nil {
                sock.Reply(fmt.Sprintf("%s", er))
                return
            }

            result: = x.String()
            sock.Reply(result)

            /* shell */
        case "$":
            if !isOwner {
                sock.Reply(system.Own())
                return
            }
            if !isOwner {
                sock.Reply(system.Own())
                return
            }
            sock.React("⏱️")
            out, err: = exec.Command("bash", "-c", query).Output()
            if err != nil {
                sock.Reply(fmt.Sprintf("%v", err))
                return
            }
            sock.Reply(string(out))
            sock.React("✅")

    
            // end
    }
}