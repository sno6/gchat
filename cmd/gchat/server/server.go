package main

import (
	"flag"
	"log"
	"os"

	"github.com/sno6/gchat/chat"
)

func main() {
	addr := flag.String("h", "127.0.0.1:8080", "Host and port address the server will run from. (defaults to localhost:8080).")
	flag.Parse()

	srv := chat.NewServer(*addr, log.New(os.Stdout, "", log.LstdFlags))
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
