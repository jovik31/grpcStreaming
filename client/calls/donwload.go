package calls

import (
	"context"
	pbclient "goProjects/testCSV/proto"
	"io"
	"log"
)

type MessageMetadata struct {
	NumberVariable  uint32 `json:"numberVariable"`
	PredictVariable string `json:"predictVariable"`
	VariableNames   string `json:"variableNames"`
	VariableTypes   string `json:"variableTypes"`
	IsCleaned       bool   `json:"isCleaned"`
	IsEncoded       bool   `json:"isEncoded"`
	IsStandardized  bool   `json:"isStandardized"`
	Description     string `json:"description"`
}

func CallDonwload(client pbclient.TestServiceClient, fileName string) ([]byte, *pbclient.Metadata) {

	var recF []byte
	var mdResponse *pbclient.Metadata
	log.Printf("Downloading file")

	request := &pbclient.FileRequest{
		FileName: fileName,
	}

	fileStream, err := client.Download(context.Background(), request)
	if err != nil {
		log.Fatalf("Could not download file")
	}

	for {
		chunkResponse, err := fileStream.Recv()
		if err == io.EOF {
			log.Println("Received all chunks")
			break
		}
		if err != nil {
			log.Fatalf("Error in receiving chunks %v", err)
		}

		recF = append(recF, chunkResponse.FileChunk...)
		mdResponse = chunkResponse.Metadata
		log.Print(chunkResponse)
	}

	return recF, mdResponse
}
