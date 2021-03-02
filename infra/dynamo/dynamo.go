package dynamo

import (
	"fmt"
	"os"
	"sanepar-level/domain/entity"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type reportItem struct {
	Datetime string       `json:"datetime"`
	Dams     []entity.Dam `json:"dams"`
}

// SaveReport -
func SaveReport(report entity.Report) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	item := reportItem{
		Datetime: report.UpdatedAt.Format(time.RFC3339),
		Dams:     report.Dams,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling new movie item:")
		fmt.Println(err.Error())
		return err
	}

	// Create item in table Movies
	tableName := os.Getenv("DYNAMO_TABLE")
	fmt.Println("tablename ", tableName)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Successfully added  %+v", item)

	return nil
}
