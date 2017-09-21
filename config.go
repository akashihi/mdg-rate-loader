package main

import (
	"flag"
)

type Configuration struct {
	Host     string
	Database string
	User     string
	Password string
	Period   int
}

func config() Configuration {
	addressPtr := flag.String("host", "127.0.0.1", "Database host")
	portPtr := flag.String("database", "mdg", "Database name")
	userPtr := flag.String("user", "mdg", "Database user")
	passwordPtr := flag.String("password", "mdg", "Database password")
	periodPtr := flag.Int("period", 60, "Period of currency rates polling in minutes")

	flag.Parse()

	return Configuration{
		Host:     *addressPtr,
		Database: *portPtr,
		User:     *userPtr,
		Password: *passwordPtr,
		Period:   *periodPtr,
	}
}
