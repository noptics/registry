package data

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
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
			"channels_table":             "noptics_registry_channels",
			"channels_field_cluster":     "clusterID",
			"channels_field_channel":     "channel",
			"channels_field_files":       "files",
			"channels_field_message":     "message",
			"clusters_table":             "noptics_registry_clusters",
			"clusters_field_id":          "clusterID",
			"clusters_field_name":        "clusterName",
			"clusters_field_description": "clusterDescription",
			"clusters_field_servers":     "servers",
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
		expression.Name(db.config["channels_field_files"]),
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
			db.config["channels_field_cluster"]: {
				S: aws.String(cluster),
			},
			db.config["channels_field_channel"]: {
				S: aws.String(channel),
			},
		},
		TableName:        aws.String(db.tablePrefix + db.config["channels_table"]),
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
			db.config["channels_field_cluster"]: {
				S: aws.String(cluster),
			},
			db.config["channels_field_channel"]: {
				S: aws.String(channel),
			},
		},
		TableName: aws.String(db.tablePrefix + db.config["channels_table"]),
	})

	if err != nil {
		return "", nil, err
	}

	if len(resp.Item) == 0 {
		return "", nil, nil
	}

	f := []*registrygrpc.File{}

	err = dynamodbattribute.Unmarshal(resp.Item[db.config["channels_field_files"]], &f)
	if err != nil {
		return "", nil, err
	}

	msg := ""
	if att, ok := resp.Item[db.config["channels_field_message"]]; ok {
		msg = *att.S
	}

	return msg, f, nil
}

func (db *dynamo) SetChannelMessage(cluster, channel, message string) error {
	update := expression.Set(
		expression.Name(db.config["channels_field_message"]),
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
			db.config["channels_field_cluster"]: {
				S: aws.String(cluster),
			},
			db.config["channels_field_channel"]: {
				S: aws.String(channel),
			},
		},
		TableName:        aws.String(db.tablePrefix + db.config["channels_table"]),
		UpdateExpression: expr.Update(),
	}

	db.l.Debugw("set channel message", "channel", channel, "cluster", cluster, "message", message, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return err
}

func (db *dynamo) GetChannels(cluster string) ([]string, error) {
	keyExp := expression.Key(db.config["channels_field_cluster"]).Equal(expression.Value(cluster))
	proj := expression.NamesList(expression.Name(db.config["channels_field_channel"]))

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
		TableName:                 aws.String(db.tablePrefix + db.config["channels_table"]),
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
		expression.Name(db.config["channels_field_message"]),
		expression.Value(message),
	).Set(
		expression.Name(db.config["channels_field_files"]),
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
			db.config["channels_field_cluster"]: {
				S: aws.String(cluster),
			},
			db.config["channels_field_channel"]: {
				S: aws.String(channel),
			},
		},
		TableName:        aws.String(db.tablePrefix + db.config["channels_table"]),
		UpdateExpression: expr.Update(),
	}

	db.l.Debugw("set channel data", "channel", channel, "cluster", cluster, "message", message, "files", files, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return err
}

func (db *dynamo) SaveCluster(cluster *registrygrpc.Cluster) (string, error) {
	if cluster.Id == "" {
		cluster.Id = uuid.New().String()
	}

	data, err := dynamodbattribute.Marshal(cluster)
	if err != nil {
		return "", err
	}

	update := expression.Set(
		expression.Name(db.config["clusters_field_id"]),
		expression.Value(cluster.Id),
	).Set(
		expression.Name(db.config["clusters_field_data"]),
		expression.Value(data),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return "", err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			db.config["clusters_field_id"]: {
				S: aws.String(cluster.Id),
			},
		},
		TableName:        aws.String(db.tablePrefix + db.config["clusters_table"]),
		UpdateExpression: expr.Update(),
	}

	db.l.Debugw("save cluster data", "clustser", cluster, "input", input)

	_, err = dbClient(db).UpdateItem(input)

	return cluster.Id, err
}

func (db *dynamo) GetCluster(id string) (*registrygrpc.Cluster, error) {
	item, err := dbClient(db).GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			db.config["clusters_field_id"]: {
				S: aws.String(id),
			},
		},
		TableName: aws.String(db.tablePrefix + db.config["clusters_table"]),
	})

	if err != nil {
		return nil, err
	}

	cluster := &registrygrpc.Cluster{}

	if data, ok := item.Item[db.config["clusters_field_data"]]; !ok {
		return nil, fmt.Errorf("could not find cluster %s", id)
	} else {
		err = dynamodbattribute.Unmarshal(data, cluster)
	}

	return cluster, err
}

func (db *dynamo) GetClusters() ([]*registrygrpc.Cluster, error) {
	proj := expression.NamesList(expression.Name(db.config["clusters_field_data"]))

	exp, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}

	resp, err := dbClient(db).Scan(&dynamodb.ScanInput{
		TableName:            aws.String(db.tablePrefix + db.config["clusters_table"]),
		ProjectionExpression: exp.Projection(),
	})

	if err != nil {
		return nil, err
	}

	clusters := []*registrygrpc.Cluster{}

	for _, i := range resp.Items {
		c := &registrygrpc.Cluster{}
		err = dynamodbattribute.Unmarshal(i[db.config["clusters_field_data"]], c)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling cluster")
		}
		clusters = append(clusters, c)
	}

	return clusters, nil
}
