package main

import (
	"io"
	"log"
	"net"
	"os"

	pbserver "goProjects/testCSV/proto"
	"google.golang.org/grpc"
)

const (
	port = ":8081"
)

type testServer struct {
	pbserver.TestServiceServer
}

func (srv *testServer) Download(req *pbserver.FileRequest, responseStream pbserver.TestService_DownloadServer) error {

	bufferSize := 64 * 1024 // 64KiB
	file, err := os.Open(req.GetFileName())
	if err != nil {
		log.Fatalf("Not able to open file %v", err)
		return err
	}

	defer file.Close()
	buff := make([]byte, bufferSize)
	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Could not read file, %v", err)
			}
			metda := &pbserver.Metadata{

				NumberVariable:  4,
				PredictVariable: "NLab",
				VariableNames:   "Temp, Lux, O2",
				VariableTypes:   "C, lumen, PPP",
				IsCleaned:       true,
				IsEncoded:       false,
				IsStandardized:  false,
				Description:     "Test Dataset",
			}
			resp := &pbserver.FileTransfer{
				Metadata: metda,
			}
			err = responseStream.Send(resp)
			if err != nil {
				log.Fatalf("error while sending metadata %v", err)
			}

			break
		}
		resp := &pbserver.FileTransfer{
			FileChunk: buff[:bytesRead],
		}
		err = responseStream.Send(resp)
		if err != nil {
			log.Fatalf("Error while sending chunk: %v", err)
			return err

		}
	}
	return nil

}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to liste server %v", err)

	}
	grpcServer := grpc.NewServer()
	pbserver.RegisterTestServiceServer(grpcServer, &testServer{})
	log.Printf("Registered Server at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start? %v", err)
	}

}
