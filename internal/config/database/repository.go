package database

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(id uuid.UUID, preloads ...string) error
	GetWhere(where map[string]interface{}, preloads ...string) error
	GetDeleted(where map[string]interface{}) error

	FindAllWhere(where map[string]interface{}, preloads ...string) error

	Create() error
	Save() error
	Delete(id uuid.UUID) error
	Update(fields interface{}, id uuid.UUID) error
}

type repositoryImpl struct {
	db     *gorm.DB
	entity interface{}
	ctx    *gin.Context
}

func NewRepository(entity interface{}, ctx *gin.Context) Repository {
	return &repositoryImpl{
		db:     instance,
		entity: entity,
		ctx:    ctx,
	}
}

func (r *repositoryImpl) GetById(id uuid.UUID, preloads ...string) error {
	queryDb := r.checkPreloads(preloads...)
	queryDb = queryDb.First(r.entity, id)
	if queryDb.Error != nil {
		return queryDb.Error
	}

	return nil
}

func (r *repositoryImpl) GetWhere(where map[string]interface{}, preloads ...string) error {
	queryDb := r.checkPreloads(preloads...)

	if len(where) > 1 {
		first := true
		for key, value := range where {
			condition := map[string]interface{}{key: value}
			if first {
				queryDb = queryDb.Where(condition)
				first = false
			} else {
				queryDb = queryDb.Or(condition)
			}
		}
	} else {
		queryDb = queryDb.Where(where).First(r.entity)
	}

	queryDb = queryDb.Find(r.entity).First(r.entity)
	if queryDb.Error != nil {
		return queryDb.Error
	}

	return nil
}

func (r *repositoryImpl) GetDeleted(where map[string]interface{}) error {
	queryDb := r.db
	queryDb = queryDb.Unscoped().Where(where).First(r.entity)
	if queryDb.Error != nil {
		return queryDb.Error
	}

	return nil
}

func (r *repositoryImpl) FindAllWhere(where map[string]interface{}, preloads ...string) error {
	queryDb := r.checkPreloads(preloads...)

	if len(where) > 1 {
		first := true
		for key, value := range where {
			condition := map[string]interface{}{key: value}
			if first {
				queryDb = queryDb.Where(condition)
				first = false
			} else {
				queryDb = queryDb.Or(condition)
			}
		}
	} else {
		queryDb = queryDb.Where(where)
	}

	queryDb = queryDb.Find(r.entity)

	if queryDb.Error != nil {
		return queryDb.Error
	}

	return nil
}

func (r *repositoryImpl) Create() error {
	queryDb := r.db
	queryDb = queryDb.Create(r.entity)

	if queryDb.Error != nil {
		return queryDb.Error
	}

	if queryDb.RowsAffected == 0 {
		return fmt.Errorf("rows affected equals zero")
	}

	return nil
}

func (r *repositoryImpl) Save() error {
	queryDb := r.db
	queryDb = queryDb.Save(r.entity)
	if queryDb.Error != nil {
		return queryDb.Error
	}

	if queryDb.RowsAffected == 0 {
		return fmt.Errorf("rows affected equals zero")
	}

	return nil
}

func (r *repositoryImpl) Delete(id uuid.UUID) error {
	queryDb := r.db
	queryDb = queryDb.Delete(r.entity, id)
	if queryDb.Error != nil {
		return r.db.Error
	}
	if queryDb.RowsAffected == 0 {
		return fmt.Errorf("rows affected equals zero")
	}

	return nil
}

func (r *repositoryImpl) Update(fields interface{}, id uuid.UUID) error {
	queryDb := r.db
	queryDb = queryDb.Model(r.entity).Where("id = ?", id).Updates(fields).Find(r.entity)

	if queryDb.Error != nil {
		return queryDb.Error
	}

	if queryDb.RowsAffected == 0 {
		return fmt.Errorf("rows affected equals zero")
	}

	return nil
}

func (r *repositoryImpl) checkPreloads(args ...string) *gorm.DB {
	if len(args) == 0 {
		return r.db
	}

	queryDb := r.db
	fmt.Printf("Query: %t", queryDb.Statement.Unscoped)
	for _, arg := range args {
		queryDb = queryDb.Preload(arg)
	}
	return queryDb
}
