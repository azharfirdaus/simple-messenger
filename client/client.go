package client

import (
	"github.com/azhar.firdaus/simple-messenger/config"
	"github.com/azhar.firdaus/simple-messenger/dao"
)

var GlobalClient *Client

type Client struct {
	ChatDAO dao.ChatDAO
}

func NewClient(Config *config.Config) (*Client, error) {
	dao, err := dao.NewMongoDAO(*Config.MongoDbUri, "simple_messenger", "chat")
	if err != nil {
		return nil, err
	}

	return &Client{
		ChatDAO: dao,
	}, nil
}
