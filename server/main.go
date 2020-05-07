package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"

	pb "github.com/campoy/justforfunc/12-say-grpc/api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "8000", "required port")
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		logrus.Fatalf("could not connect to the port %d", err)
	}

	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		logrus.Println("could not listen", err)
	}

}

type server struct{}

func (*server) say(ctx context.Context, text *pb.Text) (*pb.Speech, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("unable to create a tempfile %s", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close %s: %v", f.Name(), err)
	}
	cmd := exec.Command("flite", text.Text, f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("flite failed: %s", data)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create a tempfile %s", err)
	}

	return &pb.Speech{Audio: data}, nil

}
