package repository

import (
	"github.com/7joe7/csvstreamtest/common/model"
	"github.com/7joe7/csvstreamtest/common/types"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// ORMClient is a repository to manage the clients.
type ORMClient struct {
	db *gorm.DB
}

// NewORMClient creates a new client repository using the given gorm DB.
func NewORMClient(db *gorm.DB) *ORMClient {
	return &ORMClient{db: db}
}

// Store saves a new client in the database.
func (repo *ORMClient) Store(client *model.Client) error {
	err := repo.db.Create(client).Error
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		// if the error is of type duplicate entry
		if mysqlError.Number == 1062 {
			return types.ErrDuplicateEntry
		}
	}
	return err
}

// Retrieve returns the client with the given ID from the database.
func (repo *ORMClient) Retrieve(id int32) (*model.Client, error) {
	var client model.Client
	err := repo.db.First(&client, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, types.ErrNotFound
	}
	return &client, err
}

// Update updates the client in the database.
func (repo *ORMClient) Update(newClient *model.Client) error {
	var client model.Client
	err := repo.db.First(&client, newClient.Id).Error
	if err == gorm.ErrRecordNotFound {
		return types.ErrNotFound
	}
	if err != nil {
		return err
	}
	err = repo.db.Save(&newClient).Error
	if msErr, ok := err.(*mysql.MySQLError); ok {
		// if the error is of type duplicate entry
		if msErr.Number == 1062 {
			return types.ErrDuplicateEntry
		}
	}
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes the Client from the database.
func (repo *ORMClient) Delete(id uint) error {
	var Client model.Client
	err := repo.db.First(&Client, id).Error
	if err == gorm.ErrRecordNotFound {
		return types.ErrNotFound
	}
	return repo.db.Delete(&Client).Error
}

// Find retrieves the Clients fitting the provided filters from the database.
func (repo *ORMClient) Find(clientIDs []uint) ([]*model.Client, error) {
	var clients []*model.Client
	query := repo.db
	if len(clientIDs) > 0 {
		query = query.Where("id IN (?)", clientIDs)
	}
	return clients, query.Find(&clients).Error
}

// NewWithTransaction creates new ORMClient with new transaction begun. The transaction can be propagated to other repositories.
// The ORMClient object then serves as a transaction object (Rollback, Commit functions)
// It is a shortcut basically, cleanest would be to have transaction object, then create new repos with that object.
func (repo *ORMClient) NewWithTransaction() (*ORMClient, error) {
	ormClient := &ORMClient{db: repo.db.Begin()}
	err := ormClient.db.Error
	if err != nil {
		return nil, err
	}
	return ormClient, nil
}

// Rollback performs gorm rollback function
func (repo *ORMClient) Rollback() error {
	return repo.db.Rollback().Error
}

// Commit performs gorm commit function
func (repo *ORMClient) Commit() error {
	return repo.db.Commit().Error
}
