package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "helicopter/generated/proto"
)

func main() {
	serverAddr := flag.String("server_addr", "localhost:1228", "The server address in the format of host:port")
	mode := flag.String("mode", "listen", "The mode of the CLI messenger, can be 'listen' or 'write'")
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewHelicopterClient(conn)

	if *mode == "listen" {
		listen(client)
	} else if *mode == "write" {
		write(client)
	} else {
		fmt.Println("Invalid mode, must be 'listen' or 'write'")
	}
}

func listen(client pb.HelicopterClient) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter thread ID: ")
	threadID, _ := reader.ReadString('\n')
	threadID = strings.TrimSpace(threadID)

	lastLseq := "0"
	for {
		req := &pb.GetNodesRequest{Root: threadID, Last: lastLseq}
		resp, err := client.GetNodes(context.Background(), req)
		if err != nil {
			fmt.Printf("Error fetching nodes: %v\n", err)
			continue
		}

		for _, node := range resp.Nodes {
			fmt.Printf("↩ %s: [%s] : %s\n", node.Parent, node.Lseq, node.Content)
			lastLseq = node.Lseq
		}

		time.Sleep(3 * time.Second)
	}
}

func write(client pb.HelicopterClient) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter chat ID: ")
	chatID, _ := reader.ReadString('\n')
	chatID = strings.TrimSpace(chatID)

	for {
		fmt.Print("Enter message: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		if strings.HasPrefix(msg, ":ref") {
			parts := strings.SplitN(msg[5:], " ", 2)
			if len(parts) != 2 {
				fmt.Println("Invalid reference format, must be ':ref refId contents'")
				continue
			}
			refID := strings.TrimSpace(parts[0])
			contents := []byte(parts[1])
			req := &pb.AddNodeRequest{Parent: refID, Content: contents}
			resp, err := client.AddNode(context.Background(), req)
			if err != nil {
				fmt.Printf("Error adding node: %v\n", err)
				continue
			}
			fmt.Printf("Added node with lseq %s\n", resp.Node.Lseq)
		} else {
			contents := []byte(msg)
			req := &pb.AddNodeRequest{Parent: chatID, Content: contents}
			resp, err := client.AddNode(context.Background(), req)
			if err != nil {
				fmt.Printf("Error adding node: %v\n", err)
				continue
			}
			fmt.Printf("Added node with lseq %s\n", resp.Node.Lseq)
		}
	}
}
