package main

import (
	"context"
	pb "go3/tasks/task3/server/proto/consignment"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Solve(*pb.Coefficients) *pb.Solution
	GetAll() *pb.Solutions
}

// Repository struct
type Repository struct {
	solutions pb.Solutions
}

type service struct {
	repo repository
}

func solve(a int32, b int32, c int32) int32 {
	if a == 0 {
		if b != 0 {
			return 1
		}

		return 0
	}

	d := (b * b) - (4 * a * c)

	if d > 0 {
		return 2
	} else if d == 0 {
		return 1
	}

	return 0
}

// Solve method
func (r *Repository) Solve(coefficients *pb.Coefficients) *pb.Solution {
	solution := pb.Solution{Coefs: coefficients, NRoots: solve(coefficients.A, coefficients.B, coefficients.C)}

	r.solutions.Solutions = append(r.solutions.Solutions, &solution)

	return &solution
}

// GetAll method
func (r *Repository) GetAll() *pb.Solutions {
	return &r.solutions
}

// Solve method
func (s *service) Solve(ctx context.Context, coefs *pb.Coefficients) (*pb.Solution, error) {
	solution := s.repo.Solve(coefs)

	log.Printf("Solution: %v", solution)

	return solution, nil
}

// GetAll method
func (s *service) GetAll(ctx context.Context, req *pb.GetRequest) (*pb.Solutions, error) {
	return s.repo.GetAll(), nil
}

func main() {
	repo := &Repository{}
	ourService := &service{repo: repo}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen port %v, error: %v", port, err)
	}

	server := grpc.NewServer()
	pb.RegisterSolverServer(server, ourService)
	//Чтобы выходные параметры сервера сохранялись в go-runtime
	reflection.Register(server)

	log.Println("gRPC server running on port:", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve from port %v", port)
	}
}
