package data

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/sethjback/noptics/golog"
	"github.com/sethjback/noptics/registry/registrygrpc"
)

// Store for the registry data
type Store interface {
	SaveFiles(cluster, channel string, files []*registrygrpc.File) error
	SetChannelMessage(cluster, channel, message string) error
	GetChannelData(cluster, channel string) (string, []*registrygrpc.File, error)
}

const (
	NOPTICS_REGISTRY_TABLE         = "noptics_registry"
	NOPTICS_REGISTRY_FIELD_CLUSTER = "cluster"
	NOPTICS_REGISTRY_FIELD_CHANNEL = "channel"
	NOPTICS_REGISTRY_FIELD_FILES   = "files"
	NOPTICS_REGISTRY_FIELD_MESSAGE = "message"
)

// dynamo implements the store interface for dynamodb
type dynamo struct {
	// awsSession holds access credentials to the db
	awsSession *session.Session
	// endpoint is used for testing. Empty is the default and will connect to live dynamodb servers
	endpoint string
	// tablePrefix is used to segregate data by env
	tablePrefix string
	// config for the tables/etc.
	config map[string]string

	l golog.Logger
}

func NewDynamodb(session *session.Session, endpoint, tablePrefix string, l golog.Logger) Store {
	return &dynamo{
		awsSession:  session,
		endpoint:    endpoint,
		tablePrefix: tablePrefix,
		config: map[string]string{
			"table":         "noptics_registry",
			"field_cluster": "cluster",
			"field_channel": "channel",
			"field_files":   "files",
			"field_message": "message",
		},
		l: l,
	}
}

func dbClient(db *dynamo) *dynamodb.DynamoDB {
	var ep *string
	if len(db.endpoint) == 0 {
		ep = nil
	} else {
		ep = aws.String(db.endpoint)
	}
	return dynamodb.New(db.awsSession, &aws.Config{Endpoint: ep})
}

func (db *dynamo) SaveFiles(cluster, channel string, files []*registrygrpc.File) error {
	update := expression.Set(
		expression.Name(NOPTICS_REGISTRY_FIELD_FILES),
		expression.Value(files),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			NOPTICS_REGISTRY_FIELD_CLUSTER: {
				S: aws.String(cluster),
			},
			NOPTICS_REGISTRY_FIELD_CHANNEL: {
				S: aws.String(channel),
			},
		},
		TableName:        aws.String(db.tablePrefix + NOPTICS_REGISTRY_TABLE),
		UpdateExpression: expr.Update(),
	}

	db.l.Debugw("save files data", "files", files, "channel", channel, "cluster", cluster, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return err
}

func (db *dynamo) GetChannelData(cluster, channel string) (string, []*registrygrpc.File, error) {
	client := dbClient(db)

	db.l.Debugw("get files", "cluster", cluster, "channel", channel)

	resp, err := client.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			NOPTICS_REGISTRY_FIELD_CLUSTER: {
				S: aws.String(cluster),
			},
			NOPTICS_REGISTRY_FIELD_CHANNEL: {
				S: aws.String(channel),
			},
		},
		TableName: aws.String(db.tablePrefix + NOPTICS_REGISTRY_TABLE),
	})

	if err != nil {
		return "", nil, err
	}

	if len(resp.Item) == 0 {
		return "", nil, nil
	}

	f := []*registrygrpc.File{}

	err = dynamodbattribute.Unmarshal(resp.Item[NOPTICS_REGISTRY_FIELD_FILES], &f)
	if err != nil {
		return "", nil, err
	}

	msg := ""
	if att, ok := resp.Item[NOPTICS_REGISTRY_FIELD_MESSAGE]; ok {
		msg = *att.S
	}

	return msg, f, nil
}

func (db *dynamo) SetChannelMessage(cluster, channel, message string) error {
	update := expression.Set(
		expression.Name(NOPTICS_REGISTRY_FIELD_MESSAGE),
		expression.Value(message),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			NOPTICS_REGISTRY_FIELD_CLUSTER: {
				S: aws.String(cluster),
			},
			NOPTICS_REGISTRY_FIELD_CHANNEL: {
				S: aws.String(channel),
			},
		},
		TableName:        aws.String(db.tablePrefix + NOPTICS_REGISTRY_TABLE),
		UpdateExpression: expr.Update(),
	}

	db.l.Debugw("set channel message", "channel", channel, "cluster", cluster, "message", message, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return err
}
