package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc/metadata"

	"github.com/rubiagatra/grpc-vacation-app/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const port = ":9000"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterEmployeeServiceServer(s, new(employeeService))
	log.Println("Starting server on port" + port)
	s.Serve(lis)
}

type employeeService struct{}

func (s *employeeService) GetByBadgeNumber(ctx context.Context,
	req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Println("metada received: %v\n", md)
	}

	for _, e := range employees {
		if req.BadgeNumber == e.BadgeNumber {
			return &pb.EmployeeResponse{Employee: &e}, nil
		}
	}

	return nil, errors.New("Employee not found")
}

func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	for _, e := range employees {
		stream.Send(&pb.EmployeeResponse{Employee: &e})
	}
	return nil
}
func (s *employeeService) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	employees = append(employees, *req.Employee)
	return &pb.EmployeeResponse{Employee: req.Employee}, nil
}
func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		employees = append(employees, *emp.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: emp.Employee})
	}
	for _, e := range employees {
		fmt.Println(e)
	}
	return nil
}
func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	imgData := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File received with length: %v\n", len(imgData))
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		imgData = append(imgData, data.Data...)
	}
}

func (s *employeeService) GetAllDefault(ctx context.Context, req *pb.GetAllRequestDefault) (*pb.EmployeeResponseRepeated, error) {
	var employeeResp []*pb.Employee

	for _, e := range employees {
		employeeResp = append(employeeResp, &pb.Employee{
			Id:                  e.Id,
			BadgeNumber:         e.BadgeNumber,
			FirstName:           e.FirstName,
			LastName:            e.LastName,
			VacationAccrualRate: e.VacationAccrualRate,
			VacationAccrued:     e.VacationAccrued,
		})
	}
	return &pb.EmployeeResponseRepeated{Istrue: req.Istrue, Employee: employeeResp}, nil

}

var employees = []pb.Employee{
	pb.Employee{
		Id:                  1,
		BadgeNumber:         2080,
		FirstName:           "Grace",
		LastName:            "Decker",
		VacationAccrualRate: 2,
		VacationAccrued:     30,
	},
	pb.Employee{
		Id:                  2,
		BadgeNumber:         7538,
		FirstName:           "Amity",
		LastName:            "Fuller",
		VacationAccrualRate: 2.3,
		VacationAccrued:     23.4,
	},
	pb.Employee{
		Id:                  3,
		BadgeNumber:         5144,
		FirstName:           "Keaton",
		LastName:            "Willis",
		VacationAccrualRate: 3,
		VacationAccrued:     31.7,
	},
}
