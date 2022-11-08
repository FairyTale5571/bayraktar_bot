package mongo

import (
	"context"
	"github.com/fairytale5571/bayraktar_bot/pkg/errs"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"

	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// DuplicateKeyError is error code which is thrown when you try
	// to insert duplicate field in an unique field index collection
	DuplicateKeyError = 11000
	// ConnectionTimeoutInSecond is number of seconds after which the database operation will get timed out.
	ConnectionTimeoutInSecond = 600
)

type mongoDB interface {
	// MongoClient should only be used to initialize our own external library such as UCS library.
	MongoClient() (*mgo.Client, error)
	MongoDB() *mgo.Database
}

type Mongo struct {
	client  *mgo.Client
	session *mgo.Session
	logger  *logger.Wrapper
	cfg     models.Config
}

func New(cfg models.Config) (*Mongo, error) {
	// create mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mgo.Connect(ctx, options.Client().ApplyURI(cfg.MongoUri))
	if err != nil {
		return nil, err
	}
	return &Mongo{
		client: client,
		cfg:    cfg,
		logger: logger.New("Mongo"),
	}, nil
}

// MongoDB returns db instance of mongodb setting mainDb as default
// database
func (m *Mongo) MongoDB() *mgo.Database {
	return m.client.Database(m.cfg.MongoDatabase)
}

// Ok check if database is running or not
func (m *Mongo) MongoClient() (*mgo.Client, error) {
	if m.client == nil {
		return nil, errs.ErrMongoClientNotInitialized
	}
	return m.client, nil
}

// Ok check if database is running or not
func (m *Mongo) Ok() (bool, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), ConnectionTimeoutInSecond*time.Second)
	defer cancelFunc()
	if err := m.client.Ping(ctx, readpref.Primary()); err != nil {
		return false, err
	}
	return true, nil
}

func (m *Mongo) Close() error {
	return m.client.Disconnect(context.Background())
}

func (m *Mongo) Write(database, collection string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.client.Database(database).Collection(collection).InsertOne(ctx, data)
	if err != nil {
		m.logger.Errorf("Write: %s", err)
		return err
	}
	return nil
}

func (m *Mongo) Version() string {
	return "1.6"
}
