package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "helicopter/generated/proto"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	serverAddr := flag.String("server_addr", "localhost:1228", "The server address in the format of host:port")
	flag.Parse()

	log.Printf("Using server address %s\n", *serverAddr)
	// create gRPC client
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	client := pb.NewHelicopterClient(conn)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("No apiKey provided")
	}

	// create OpenAI client
	openaiClient := openai.NewClient(apiKey)

	// send initial GetNodesRequest
	root := "chatGPT"
	last := "0"
	resp, err := client.GetNodes(context.Background(), &pb.GetNodesRequest{Root: root, Last: last})
	if err != nil {
		log.Fatalf("Failed to get nodes: %v", err)
	}
	last = resp.Nodes[len(resp.Nodes)-1].Lseq

	for {
		// send periodic GetNodesRequest
		resp, err := client.GetNodes(context.Background(), &pb.GetNodesRequest{Root: root, Last: last})
		if err != nil {
			log.Printf("Failed to get nodes: %v", err)
		} else if len(resp.Nodes) > 0 {
			last = resp.Nodes[len(resp.Nodes)-1].Lseq
			for _, node := range resp.Nodes {
				if string(node.Content[:8]) == "AskChat:" {
					// trim "AskChat:" prefix from content
					content := string(node.Content[8:])
					// send chat message to OpenAI API
					resp, err := openaiClient.CreateChatCompletion(
						context.Background(),
						openai.ChatCompletionRequest{
							Model: openai.GPT3Dot5Turbo,
							Messages: []openai.ChatCompletionMessage{
								{
									Role:    openai.ChatMessageRoleUser,
									Content: content,
								},
							},
						},
					)
					if err != nil {
						log.Printf("Failed to create chat completion: %v", err)
					} else {
						// add response node to the tree
						_, err := client.AddNode(context.Background(), &pb.AddNodeRequest{
							Parent:  node.Lseq,
							Content: []byte(resp.Choices[0].Message.Content),
						})
						if err != nil {
							log.Printf("Failed to add node: %v", err)
						}
					}
				}
			}
		}
		time.Sleep(3 * time.Second)
	}
}
