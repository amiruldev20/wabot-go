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

// Package router includes the router for handling events from whatsapp web api.
package router

import (
	"errors"

	"go.mau.fi/whatsmeow/types/events"
	"go.ramdanhere.dev/whatsrhyno"
)

// Route represent the router object of the application.
type Route struct {
	Name        string // name of the route
	Description string // description of the route
	Handler     HandlerFunc
	Middleware  []MiddlewareFunc
	Matcher     func(string) bool
	Routes      []*Route

	// AutoRead will mark the message from group as read as soon the message received.
	AutoRead bool
}

// HandlerFunc is the function that will be called when the route is called.
type HandlerFunc func(*Context)

// MiddlewareFunc is the function that will be called before the route is called.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Desc returns the description of the route.
func (r *Route) Desc() string {
	return r.Description
}

// Find finds the route by name.
func (r *Route) Find(name string) *Route {
	for _, route := range r.Routes {
		if route.Matcher(name) {
			return route
		}
	}
	return nil
}

// Matcher returns the matcher function of the route.
func Matcher(r *Route) func(string) bool {
	return func(cmd string) bool {
		return cmd == r.Name
	}
}

// AddRoute adds a route to the route.
// Will return an error if the route already exists.
func (r *Route) AddRoute(route *Route) error {
	if nr := r.Find(route.Name); nr != nil {
		return errors.New("route already exists")
	}

	r.Routes = append(r.Routes, route)
	return nil
}

// OnMatch adds a handler to the route.
func (r *Route) OnMatch(name string, matcher func(string) bool, hf HandlerFunc) *Route {
	if nr := r.Find(name); nr != nil {
		return nr
	}

	nf := hf

	for _, v := range r.Middleware {
		nf = v(nf)
	}

	nr := &Route{
		Name:        name,
		Handler:     nf,
		Matcher:     matcher,
		Description: r.Description,
	}

	r.AddRoute(nr)
	return nr
}

// On Register a new route.
func (r *Route) On(name string, hf HandlerFunc) *Route {
	nr := r.OnMatch(name, nil, hf)
	nr.Matcher = Matcher(nr)
	return nr
}

// Use adds a middleware to the route.
func (r *Route) Use(mf ...MiddlewareFunc) *Route {
	r.Middleware = mf
	return r
}

// EventHandler is the function that will be called when the event is emited.
func (r *Route) EventHandler(cli *whatsrhyno.Client) func(evt interface{}) {
	return func(evt interface{}) {
		switch evt.(type) {
		case *events.Message:
			r.eventMessageHandler(cli, evt.(*events.Message))
		default:
			r.parseEvt(cli, evt)
		}
	}
}
