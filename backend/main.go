package main

import (
	//"os"
	"os/exec"
	"log"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	pb "github.com/ktinubu/text2speech/api"
	"golang.org/x/net/context"
	"io/ioutil"
)

func main() {
	port := flag.Int("p", 8080, "port to listen to")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("could not listen to port %d: %v", *port, err)
	}

	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("could not server: %v", err)
	}
}

type server struct {
}

func (server) Say(ctx context.Context, text *pb.Text) (*pb.Speech, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("could not return temp file %v", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close: %s: %v", f.Name(), err)
	}

	cmd := exec.Command("flite", "-t", text.Text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("flite failed: %s", data)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read tmp file: %v", err)
	}
	//log.Printf("returning %v, from text: %v file: %v", data, text.Text,f.Name())
	return &pb.Speech{Audio: data}, nil
}

