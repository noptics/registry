package data

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/noptics/golog"
	"github.com/noptics/registry/registrygrpc"
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
		expression.Name(db.config["field_files"]),
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
			db.config["field_cluster"]: {
				S: aws.String(cluster),
			},
			db.config["field_channel"]: {
				S: aws.String(channel),
			},
		},
		TableName:        aws.String(db.tablePrefix + db.config["table"]),
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
			db.config["field_cluster"]: {
				S: aws.String(cluster),
			},
			db.config["field_channel"]: {
				S: aws.String(channel),
			},
		},
		TableName: aws.String(db.tablePrefix + db.config["table"]),
	})

	if err != nil {
		return "", nil, err
	}

	if len(resp.Item) == 0 {
		return "", nil, nil
	}

	f := []*registrygrpc.File{}

	err = dynamodbattribute.Unmarshal(resp.Item[db.config["field_files"]], &f)
	if err != nil {
		return "", nil, err
	}

	msg := ""
	if att, ok := resp.Item[db.config["field_message"]]; ok {
		msg = *att.S
	}

	return msg, f, nil
}

func (db *dynamo) SetChannelMessage(cluster, channel, message string) error {
	update := expression.Set(
		expression.Name(db.config["field_message"]),
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
			db.config["field_cluster"]: {
				S: aws.String(cluster),
			},
			db.config["field_channel"]: {
				S: aws.String(channel),
			},
		},
		TableName:        aws.String(db.tablePrefix + db.config["table"]),
		UpdateExpression: expr.Update(),
	}

	db.l.Debugw("set channel message", "channel", channel, "cluster", cluster, "message", message, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return err
}

func (db *dynamo) GetChannels(cluster string) ([]string, error) {
	keyExp := expression.Key(db.config["field_cluster"]).Equal(expression.Value(cluster))
	proj := expression.NamesList(expression.Name(db.config["field_channel"]))

	exp, err := expression.NewBuilder().WithKeyCondition(keyExp).WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}

	client := dbClient(db)
	resp, err := client.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: exp.Values(),
		ExpressionAttributeNames:  exp.Names(),
		KeyConditionExpression:    exp.KeyCondition(),
		ProjectionExpression:      exp.Projection(),
		TableName:                 aws.String(db.tablePrefix + db.config["table"]),
	})

	if err != nil {
		return nil, err
	}

	chans := []string{}

	db.l.Debugw("query items", "cluster", cluster, "items", resp.Items)

	for _, i := range resp.Items {
		chans = append(chans, *i["channel"].S)
	}

	return chans, nil
}

func (db *dynamo) SaveChannelData(cluster, channel, message string, files []*registrygrpc.File) error {
	update := expression.Set(
		expression.Name(db.config["field_message"]),
		expression.Value(message),
	).Set(
		expression.Name(db.config["field_files"]),
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
			db.config["field_cluster"]: {
				S: aws.String(cluster),
			},
			db.config["field_channel"]: {
				S: aws.String(channel),
			},
		},
		TableName:        aws.String(db.tablePrefix + db.config["table"]),
		UpdateExpression: expr.Update(),
	}

	db.l.Debugw("set channel data", "channel", channel, "cluster", cluster, "message", message, "files", files, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return err
}
