// Copyright 2015 Horacio Duran.
// Licenced under the MIT license, see LICENCE file for details.

package main

import (
	"flag"
	"log"

	"github.com/perrito666/got/bug"
	"github.com/perrito666/got/git"
	"github.com/perrito666/got/registry"
)

// Interface Checking
var _ registry.Command = &bug.Command{}

var args []string

func init() {
	flag.Parse()
	args = flag.Args()
}

func main() {
	if len(args) == 0 {
		// error is ignored here because calling git without
		// arguments returns 1.
		// TODO(perrito666) add some form of usage for got too.
		c := &git.Call{}
		_ = c.Run()
		return
	}
	commands := registry.Commands()
	command, ok := commands[args[0]]
	if !ok {
		c := git.New(args[0], args[1:])
		err := c.Run()
		if err != nil {
			log.Fatalln(err)
		}
		return
	}
	c, err := command()
	if err != nil {
		panic(err)
	}

	if err := c.Run(args[1:]); err != nil {
		log.Fatal(err)
	}

}
