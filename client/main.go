package main

import (
	"flag"
	"fmt"
	pb "github.com/ktinubu/text2speech/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func getEnvWithBackup(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func textToSpeech(backend *string, output *string, conn *grpc.ClientConn) {

	defer conn.Close()

	client := pb.NewTextToSpeechClient(conn)
	text := &pb.Text{
		Text: flag.Arg(0),
	}
	res, err := client.Say(context.Background(), text)
	if err != nil {
		log.Fatalf("Could not say %s: %v", text.Text, err)
	}
	if err = ioutil.WriteFile(*output, res.Audio, 0666); err != nil {
		log.Fatalf("could not write to %s: %v", *output, err)
	}
}

// https://bbengfort.github.io/snippets/2017/03/21/sanely-grpc-dial-a-remote.html
func main() {
	backend := flag.String("b", "localhost:8080", "address of the say backend")
	output := flag.String("o", "output.wav", "wav file where output will be created")
	withinCluster := flag.Bool("c", false, "true if client making call from node in cluster, false if not (default is false")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \"text to speak\"\n", os.Args[0])
		os.Exit(1)
	}
	var conn *grpc.ClientConn
	var err error
	timeout := 3 * time.Second

	if *withinCluster {
		addr := "set-0.say.default.svc.cluster.local"
		conn, err = grpc.Dial(
			addr, grpc.WithInsecure(), grpc.WithBlock(),
			grpc.WithTimeout(timeout))
		if err != nil {
			log.Fatalf("could not connect to %s: %v", *backend, err)
		}
		log.Print("...")
	} else {
		addr := *backend
		conn, err = grpc.Dial(
			addr, grpc.WithInsecure(), grpc.WithBlock(),
			grpc.WithTimeout(timeout))
		if err != nil {
			log.Fatalf("could not connect to %s: %v", *backend, err)
		}
		log.Print("(((")
	}
	textToSpeech(backend, output, conn)

}
