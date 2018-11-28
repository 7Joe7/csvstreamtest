package repository

import "github.com/7joe7/csvstreamtest/common/model"

type ClientRepo interface {
	Store(client *model.Client) error
	Retrieve(id int32) (*model.Client, error)
	Update(newClient *model.Client) error
	Delete(id uint) error
	Find(clientIDs []uint) ([]*model.Client, error)
	NewWithTransaction() (ClientRepo, error)
	Rollback() error
	Commit() error
}
