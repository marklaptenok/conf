package conf

import (
	"errors"
	"flag"
	"net"
	"os"
	"strconv"
	"time"

	"codelearning.online/logger"
)

type ClociConfiguration struct {
	bind_address       net.IP
	bind_port          uint16
	route              string
	tls_cert_path      string
	capacity           uint16
	rate_limit         time.Duration
	waiting_time_limit time.Duration
	timeout            time.Duration
}

var (
	cnf = ClociConfiguration{bind_address: net.IP{127, 0, 0, 1}, bind_port: 443, route: "compile", tls_cert_path: "", capacity: 10, rate_limit: 250 * time.Millisecond, waiting_time_limit: 1000 * time.Millisecond, timeout: 5 * time.Second}
)

//	Reads configuration from a file or a provided ClociConfiguration struct
//	and returns it as a final ClociConfiguration struct
func Read(options []string) error {
	//	TO-DO: read configuration from a file or a provided string and return it as a conf struct

	flags := flag.NewFlagSet("CLOCI", flag.ContinueOnError)
	//	Writes usage information to the stdout.
	flags.SetOutput(os.Stdout)

	flags.Func("bind-address", "Specify IP address to bind (e.g. \"127.0.0.1\" or \"2001:db8::68\", or \"::ffff:192.0.2.1\" (default \"127.0.0.1\"))", func(ip_address_string string) error {
		if cnf.bind_address = net.ParseIP(ip_address_string); cnf.bind_address == nil {
			return errors.New("given IP address is invalid")
		}
		return nil
	})

	flags.Func("bind-port", "Specify port to bind (e.g. 443)", func(port_string string) error {
		var (
			port uint64
			err  error
		)
		if port, err = strconv.ParseUint(port_string, 10, 16); err != nil {
			cnf.bind_port = 0
			return errors.New("given port is invalid")
		}
		cnf.bind_port = uint16(port)
		return nil
	})

	flags.Parse(options)

	if cnf.bind_address == nil {
		if location, err := logger.Get_function_name(); err == nil {
			return &logger.ClpError{1, "Given IP address is invalid", location}
		} else {
			return &logger.ClpError{1, "Given IP address is invalid", ""}
		}
	}

	if cnf.bind_port == 0 {
		if location, err := logger.Get_function_name(); err == nil {
			return &logger.ClpError{2, "Given port is invalid", location}
		} else {
			return &logger.ClpError{2, "Given port is invalid", ""}
		}
	}

	logger.Debug("%v:%d", cnf.bind_address, cnf.bind_port)

	return nil
}
