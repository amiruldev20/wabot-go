// MYWA BOT GOLANG
// LIBRARY WHATSMEOW
// MADE BY AMIRUL DEV
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.amirul.dev/system/msg"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"github.com/probandula/figlet4go"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func main() {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	dxz, err := base64.StdEncoding.DecodeString("TVlXQSBCT1Q=")
	if err != nil {
		panic("malformed input")
		log.Println(dxz)
	}
	container, err := sqlstore.New("sqlite3", "file:system/session/mywabot.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "ERROR", true)
	// Client
	client := whatsmeow.NewClient(deviceStore, clientLog)
	eventHandler := registerHandler(client)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				dxz, err := base64.StdEncoding.DecodeString("TWFkZSBieSBBbWlydWwgRGV2LiBmb2xsb3cgSUcgQGFtaXJ1bC5kZXY=")
				if err != nil {
					panic("malformed input")
				}
				log.Println(string(dxz))

				log.Println("Silahkan scan qr...")
			} else {
				log.Println("Login success...")
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		log.Println("Succes Login")
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func init() {
	ascii := figlet4go.NewAsciiRender()
	dxz, err := base64.StdEncoding.DecodeString("TVlXQSBCT1Q=")
	if err != nil {
		panic("malformed input")
	}
	renderStr, _ := ascii.Render(string(dxz))
	// Set Browser
	store.DeviceProps.PlatformType = waProto.DeviceProps_SAFARI.Enum()
	store.DeviceProps.Os = proto.String(string(dxz))
	// Print Banner
	fmt.Print(renderStr)
}

func registerHandler(client *whatsmeow.Client) func(evt interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			go message.Msg(client, v)
			break
		}
	}
}
