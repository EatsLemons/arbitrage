package main

import (
	"arbitrage/api"
	"arbitrage/exchange"
	"arbitrage/exchange/bitlish"
	"arbitrage/exchange/exmo"
	"os"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Port int `long:"port" env:"ARBITRAGE_SERVICE_PORT" default:"8080" description:"port"`
}

func main() {
	p := flags.NewParser(&opts, flags.Default)
	if _, e := p.ParseArgs(os.Args[1:]); e != nil {
		os.Exit(1)
	}

	exchange := exchange.MakeExchange()
	exchange.AddMarket(exmo.MakeExmoAPI())
	exchange.AddMarket(bitlish.MakeBitlishAPI())

	server := &api.Service{
		Exchange: exchange,
	}

	server.Run(opts.Port)
}
