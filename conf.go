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
	bind_address              net.IP
	bind_port                 uint16
	route                     string
	tls_cert_path             string
	capacity                  uint16
	rate_limit                time.Duration
	waiting_time_limit        time.Duration
	timeout                   time.Duration
	response_write_timeout    time.Duration
	tls_handshake_timeout     time.Duration
	request_read_timeout      time.Duration
	request_header_size_limit uint64
	request_body_size_limit   uint64
}

func (this *ClociConfiguration) Bind_address() net.IP {
	return this.bind_address
}

func (this *ClociConfiguration) Bind_port() uint16 {
	return this.bind_port
}

func (this *ClociConfiguration) Route() string {
	return this.route
}

func (this *ClociConfiguration) Tls_cert_path() string {
	return this.tls_cert_path
}

func (this *ClociConfiguration) Timeout() time.Duration {
	return this.timeout
}

func (this *ClociConfiguration) Response_write_timeout() time.Duration {
	return this.response_write_timeout
}

func (this *ClociConfiguration) TLS_handshake_timeout() time.Duration {
	return this.tls_handshake_timeout
}

func (this *ClociConfiguration) Request_read_timeout() time.Duration {
	return this.request_read_timeout
}

func (this *ClociConfiguration) Request_header_size_limit() uint64 {
	return this.request_header_size_limit
}

func (this *ClociConfiguration) Request_body_size_limit() uint64 {
	return this.request_body_size_limit
}

var (
	default_cnf = ClociConfiguration{
		bind_address: net.IP{127, 0, 0, 1},
		bind_port:    443,
		route:        "compile",
		//	TO-DO: get this path automatically using working directory location.
		tls_cert_path:             "/root/cloci/certificates",
		capacity:                  10,
		rate_limit:                250 * time.Millisecond,
		waiting_time_limit:        1000 * time.Millisecond,
		timeout:                   5 * time.Second,
		response_write_timeout:    10000 * time.Millisecond,
		tls_handshake_timeout:     200 * time.Millisecond,
		request_read_timeout:      200 * time.Millisecond,
		request_header_size_limit: 1 << 10, //	1 Kb
		request_body_size_limit:   1 << 14} //	16 Kb
)

//	Reads configuration from a file or a provided ClociConfiguration struct
//	and returns it as a final ClociConfiguration struct
func Read(options []string) (*ClociConfiguration, error) {
	//	TO-DO: read configuration from a file or a provided string and return it as a conf struct

	cnf := default_cnf

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
			return nil, &logger.ClpError{1, "Given IP address is invalid", location}
		} else {
			return nil, &logger.ClpError{1, "Given IP address is invalid", ""}
		}
	}

	if cnf.bind_port == 0 {
		if location, err := logger.Get_function_name(); err == nil {
			return nil, &logger.ClpError{2, "Given port is invalid", location}
		} else {
			return nil, &logger.ClpError{2, "Given port is invalid", ""}
		}
	}

	logger.Debug("%v:%d", cnf.bind_address, cnf.bind_port)

	return &cnf, nil
}
