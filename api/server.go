package api

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

var TycoonHost = "http://127.0.0.1"
var TycoonPort = ":2377"

type dispatcher struct {
	handler http.Handler
}

func (d *dispatcher) SetHandler(handler http.Handler) {
	d.handler = handler
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if d.handler == nil {
		http.Error(w, "No dispatcher defined", http.StatusInternalServerError)
	}
	d.handler.ServeHTTP(w, r)
}

func newListener(proto, addr string) (net.Listener, error) {
	l, err := net.Listen(proto, addr)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") && strings.Contains(addr, TycoonPort) {
			return nil, fmt.Errorf("%s: is Docker already running on this machine? Try using a different port", err)
		}
		return nil, err
	}
	return l, nil
}

type Server struct {
	hosts      []string
	dispatcher *dispatcher
}

func NewServer(hosts []string) *Server {
	return &Server{
		hosts:      hosts,
		dispatcher: &dispatcher{},
	}
}

func (s *Server) SetHandler(handler http.Handler) {
	s.dispatcher.SetHandler(handler)
}

func (s *Server) ListenAndServe() error {
	chErrors := make(chan error, len(s.hosts))

	for _, host := range s.hosts {
		protoAddrParts := strings.SplitN(host, "://", 2)
		if len(protoAddrParts) == 1 {
			protoAddrParts = append([]string{"tcp"}, protoAddrParts...)
		}

		go func() {
			// log.WithFields(log.Fields{"proto": protoAddrParts[0], "addr": protoAddrParts[1]}).Info("Listening for HTTP")

			var (
				l      net.Listener
				err    error
				server = &http.Server{
					Addr:    protoAddrParts[1],
					Handler: s.dispatcher,
				}
			)

			switch protoAddrParts[0] {
			case "unix":
				l, err = newUnixListener(protoAddrParts[1])
			case "tcp":
				l, err = newListener("tcp", protoAddrParts[1])
			default:
				err = fmt.Errorf("unsupported protocol: %q", protoAddrParts[0])
			}
			if err != nil {
				chErrors <- err
			} else {
				chErrors <- server.Serve(l)
			}

		}()
	}

	for i := 0; i < len(s.hosts); i++ {
		err := <-chErrors
		if err != nil {
			return err
		}
	}
	return nil
}
