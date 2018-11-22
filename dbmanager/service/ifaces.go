package service

import (
	"github.com/7joe7/csvstreamtest/common/model"
	"github.com/7joe7/csvstreamtest/dbmanager/repository"
)

type clientRepo interface {
	Store(client *model.Client) error
	Retrieve(id int32) (*model.Client, error)
	Update(newClient *model.Client) error
	Delete(id uint) error
	Find(clientIDs []uint) ([]*model.Client, error)
	NewWithTransaction() (*repository.ORMClient, error)
	Rollback() error
	Commit() error
}
