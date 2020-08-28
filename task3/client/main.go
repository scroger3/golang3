package main

import (
	"context"
	"encoding/json"
	pb "go3/tasks/task3/client/proto/consignment"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "coefficients.json"
)

func parseJSON(file string) (*pb.Coefficients, error) {
	var coefficients *pb.Coefficients

	fileBody, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(fileBody, &coefficients)

	return coefficients, nil
}

func main() {
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can not connect to address %v, error: %v", address, err)
	}
	defer connection.Close()

	client := pb.NewSolverClient(connection)
	coefficients, err := parseJSON(defaultFilename)
	if err != nil {
		log.Fatalf("Can not read file %v, error: %v", defaultFilename, err)
	}

	solution, err := client.Solve(context.Background(), coefficients)
	if err != nil {
		log.Fatalf("Can not solve: %v", err)
	}

	log.Println(solution.NRoots)

	getAll, err := client.GetAll(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Can not get all solutions: %v", err)
	}

	for _, v := range getAll.Solutions {
		log.Printf("A: %v B: %v C: %v NRoots: %v", v.Coefs.A, v.Coefs.B, v.Coefs.C, v.NRoots)
	}
}
