package main

import (
	"flag"
	"log"

	"github.com/sno6/gchat/chat"
	"github.com/sno6/gchat/ui"
)

func main() {
	addr := flag.String("h", "127.0.0.1:8080", "Host and port address the server will run from. (defaults to localhost:8080).")
	name := flag.String("u", "", "Client user name.")
	flag.Parse()

	cli, err := chat.NewClient(*addr, *name)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	cui, err := ui.NewCUI(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer cui.Close()

	go func() {
		if err := cli.ReadPump(cui); err != nil {
			log.Fatalf("ReadPump error: %v", err)
		}
	}()

	if err := cui.Run(); err != nil {
		log.Fatal(err)
	}
}
