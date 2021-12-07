#!/bin/bash
go build -o bookings cmd/web/*.go
./bookings -production=false -cache=false -dbname=bookings -dbuser=postgres -dbpass=123hj123