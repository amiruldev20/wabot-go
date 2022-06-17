package whatsrhyno

import (
	"context"
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	walog "go.mau.fi/whatsmeow/util/log"
)

func pf() string {
	x := "77686174737268796e6f"
	b, _ := hex.DecodeString(x)
	return string(b)
}

// Client is a wrapper for whatsmeow client with some additional features.
type Client struct {
	cli           *whatsmeow.Client // The underlying client.
	evtHandlerID  uint32
	commandPrefix string
}

// NewClient initializes a new client.
// dbFile is the file path that will be used to store session and other data.
// whatsrhyno will use `sqlite3` as the sql dialect.
// loglevel is the log level of the underlying client.
// e.g "WARN", "INFO", "DEBUG", "ERROR"
func NewClient(dbFile, loglevel string) (*Client, error) {
	dbLog := walog.Stdout("Database", loglevel, true)
	waVersion := store.GetWAVersion()
	store.SetOSInfo(pf(), waVersion)
	container, err := sqlstore.New("sqlite3", "file:"+dbFile+"?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}

	devicestore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	clog := walog.Stdout("Client", loglevel, true)
	client := whatsmeow.NewClient(devicestore, clog)

	return &Client{
		cli: client,
	}, nil
}

// SetCommandPrefix sets the command prefix.
func (cli *Client) SetCommandPrefix(prefix string) {
	cli.commandPrefix = prefix
}

// GetCommandPrefix returns the command prefix.
func (cli *Client) GetCommandPrefix() string {
	return cli.commandPrefix
}

// RegisterEventHandler registers a handler for events.
func (cli *Client) RegisterEventHandler(evtHandler func(evt any)) {
	cli.evtHandlerID = cli.cli.AddEventHandler(evtHandler)
}

// Connect connects the client to whatsapp web websocket.
func (cli *Client) Connect() error {
	if cli.cli.Store.ID == nil {
		qrchan, _ := cli.cli.GetQRChannel(context.Background())
		err := cli.cli.Connect()
		if err != nil {
			return err
		}

		for evt := range qrchan {
			switch evt.Event {
			case "code":
				// show qr on the terminal
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)

			default:
				log.Printf("Login event: %s\n", evt.Event)
			}
		}
	} else {
		err := cli.cli.Connect()
		if err != nil {
			return err
		}
	}

	// send the presence to the server so other user will see online status.
	cli.cli.SendPresence(types.PresenceAvailable)

	log.Printf("Client %s connected to whatsapp.", cli.cli.Store.ID.String())

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	cli.Disconnect()

	return nil
}

// Disconnect disconnects the client from whatsapp web websocket.
func (cli *Client) Disconnect() {
	cli.cli.Disconnect()
}

// GetCli returns the underlying client.
func (cli *Client) GetCli() *whatsmeow.Client {
	return cli.cli
}