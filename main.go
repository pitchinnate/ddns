package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

type IpResponse struct {
	Ip string `json:"ip"`
}

func main() {
	log.Print("running ddns")
	f, err := os.OpenFile("ddns.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := route53.NewFromConfig(cfg)

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	updateIp(svc)
	for {
		select {
		case <-ticker.C:
			updateIp(svc)
		}
	}
}

func updateIp(svc *route53.Client) {
	log.Print("Getting ip from service")
	response, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		log.Printf("Error getting ip", err)
		return
	}
	ipResponse := IpResponse{}
	err = json.NewDecoder(response.Body).Decode(&ipResponse)
	if err != nil {
		log.Printf("Error decoding ip response", err)
		return
	}
	log.Printf("Got ip response: %s", ipResponse.Ip)

	hostedZoneID := os.Getenv("HOSTED_ZONE_ID")
	recordName := os.Getenv("RECORD_NAME")
	recordType := types.RRTypeA
	recordValue := ipResponse.Ip
	ttl := int64(300)

	log.Printf("Trying to created record set, %s %s %s %d", recordName, recordType, recordValue, ttl)

	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(recordName),
						Type: recordType,
						TTL:  aws.Int64(ttl),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(recordValue),
							},
						},
					},
				},
			},
		},
	}

	// Call Route 53 to change the record set
	_, err = svc.ChangeResourceRecordSets(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to create record set, %v", err)
	} else {
		log.Printf("Created record set, %s %s %s", recordName, recordType, recordValue)
	}
}
