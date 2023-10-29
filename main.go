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
package main

import (
    "mywa-bot/config"
    "context"
    "encoding/base64"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "mywa-bot/system"
    "go"
    _ "github.com/mattn/go-sqlite3"
    "github.com/mdp/qrterminal"
    "github.com/probandula/figlet4go"
    "github.com/amiruldev20/waSocket"
    waProto "github.com/amiruldev20/waSocket/binary/proto"

    "github.com/amiruldev20/waSocket/store"
    "github.com/amiruldev20/waSocket/store/sqlstore"
    "github.com/amiruldev20/waSocket/types/events"
    waLog "github.com/amiruldev20/waSocket/util/log"
    "google.golang.org/protobuf/proto"
)

func main() {
    dbLog: = waLog.Stdout("Database", "ERROR", true)
    dxz,
    err: = base64.StdEncoding.DecodeString("TXl3YSBCT1QgQnkgd2FTb2NrZXQ=")
    if err != nil {
        panic("malformed input")
        log.Println(dxz)
    }
    container,
    err: = sqlstore.New("sqlite3", "file:mywabot.db?_foreign_keys=on", dbLog)
    if err != nil {
        panic(err)
    }

    deviceStore,
    err: = container.GetFirstDevice()
    if err != nil {
        panic(err)
    }
    clientLog: = waLog.Stdout("Client", "ERROR", true)
        
	/* client */
    client: = waSocket.NewClient(deviceStore, clientLog)
    eventHandler: = registerHandler(client)
    client.AddEventHandler(eventHandler)

    if client.Store.ID == nil {

        if config.TypeLogin == "code" {
            fmt.Println("You login with pairing code")

            err = client.Connect()
            if err != nil {
                panic(err)
            }

            // don't edit
            code, err: = client.PairPhone(config.BotNumber, true, waSocket.PairClientChrome, "Chrome (Linux)")

            if err != nil {
                fmt.Println(err)
                return
            }

            log.Println("Your Code: " + code)

        } else {
            qrChan, _: = client.GetQRChannel(context.Background())

            err = client.Connect()
            if err != nil {
                panic(err)
            }

            for evt: = range qrChan {
                if evt.Event == "code" {
                    qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
                    dxz, err: = base64.StdEncoding.DecodeString("TXl3YSBCT1QgQnkgd2FTb2NrZXQ=")
                    if err != nil {
                        panic("malformed input")
                    }
                    log.Println(string(dxz))

                    log.Println("Please scan this QR...")
                } else {
                    log.Println("Login successfully!!")
                }
            }
        }
    } else {
        // Already logged in, just connect
        err = client.Connect()
        log.Println("Login Sucess!!")
        if err != nil {
            panic(err)
        }
    }

    // Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
    c: = make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM) < -c

    client.Disconnect()
}

func init() {
    ascii: = figlet4go.NewAsciiRender()
    dxz,
    err: = base64.StdEncoding.DecodeString("TXl3YSBCT1QgQnkgd2FTb2NrZXQ=")
    if err != nil {
        panic("malformed input")
    }
    renderStr,
    _: = ascii.Render(string(dxz))
    store.DeviceProps.PlatformType = waProto.DeviceProps_FIREFOX.Enum()
    store.DeviceProps.Os = proto.String(string(dxz))
    fmt.Print(renderStr)
}

func registerHandler(client * waSocket.Client) func(evt interface {}) {
    return func(evt interface {}) {
        switch v: = evt.(type) {
            case *events.Message:
                go system.Msg(client, v)
                break
        }
    }
}