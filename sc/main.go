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
client.SetCommandPrefix("")

// Initialize the message router.
r := &router.Route{}

// Add the ping command to the router.
r.Description = "Send ping message" // set the route description.
//r.Use(OnlyGroup)                    // use the OnlyGroup middleware.
r.On("tes", PingHandler)          
r.On("menu", menu)
r.On("ex", ex)
r.On("args", ex)

// serve the ping command.

// Send the router to the client.
evtHandler := r.EventHandler(client)
client.RegisterEventHandler(evtHandler)

// Start the client.
err = client.Connect()
if err != nil {
panic(err)
}
}

// tes

// PingHandler is a handler for the ping command.
func PingHandler(ctx *router.Context) {
ctx.Reply("halo " + ctx.Msg.Info.PushName)
}

//-- func menu
func menu(ctx *router.Context){
renz := ctx.Msg.Info
ctx.Reply(`Halo, *`+renz.PushName+`*
berikut list menu yang tersedia

*DETAIL BOT*
Name: *GO-WA BOT*
Language: *GO*
Library: *Whatsmeow*
Database: *NoSQL*

*-= MAIN MENU =-*
• .menu
• .about
• .source`)
}

//-- func sticker
func sticker(ctx *router.Context){

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
