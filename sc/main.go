package main

import (
"os/exec"


"go.ramdanhere.dev/whatsrhyno"
"go.ramdanhere.dev/whatsrhyno/router"
)

func main() {
// initialize the client.
client, err := whatsrhyno.NewClient("ninjabot.db", "ERROR")
if err != nil {
panic(err)
}

// before adding some route, we need to specify prefix to identify command from the message.
client.SetCommandPrefix(".")

// Initialize the message router.
r := &router.Route{}

r.Description = "description command"

//r.Use(OnlyGroup)      
// only group

//-- COMMAND --//
r.On("tes", tes)          
r.On("exec", ex)
r.On("args", ex)

// Send the router to the client.
evtHandler := r.EventHandler(client)
client.RegisterEventHandler(evtHandler)

// Start the client.
err = client.Connect()
if err != nil {
panic(err)
}
}

func tes(ctx *router.Context) {
ctx.Reply("halo " + ctx.Msg.Info.PushName)
}

//-- func exec
func ex(ctx *router.Context){
	if len(ctx.Args) == 0 {
		ctx.Reply("Please provide arguments.")
	}
	b, _ := exec.Command(ctx.Args[0]).Output()
	out := string(b)
	ctx.Reply(out)
}

// OnlyGroup is the middleware to set that the command can be used on group only.
func OnlyGroup(h router.HandlerFunc) router.HandlerFunc {
return func(ctx *router.Context) {
// check if the message was in private or not,
// if it was, then return because we will only
// accept the command in group.
if ctx.IsPrivate {
return
}

h(ctx)
}
}
