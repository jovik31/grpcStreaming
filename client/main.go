package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	pbclient "goProjects/testCSV/proto"
	"goProjects/testCSV/utils"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"goProjects/testCSV/client/calls"
)

const (
	port = ":8081"
)

var md calls.MessageMetadata

func main() {

	flag.Parse()
	var recF []byte
	var mdResponse *pbclient.Metadata
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to server %v", err)
	}
	defer conn.Close()

	client := pbclient.NewTestServiceClient(conn)

	jsonFile, err := os.Open("me.json")
	if err != nil {
		log.Fatalf("Cannot open json file %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &md)
	if err != nil {
		log.Fatalf("Json unmarshal failed, %v", err)
	}

	recF, mdResponse = calls.CallDonwload(client, "tranco_3VPWL.csv")
	log.Print(mdResponse)

	csvfile, err := os.Create("received.csv")
	if err != nil {
		log.Fatal(err)

	}
	defer csvfile.Close()
	csvWriter := csv.NewWriter(csvfile)

	stringSlice := utils.ToStrSlice(recF)

	if err := csvWriter.Write(stringSlice); err != nil {
		log.Fatalf("Could not write to csv %v", err)
	}
	csvWriter.Flush()

}
