package database

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetById(id uuid.UUID, preloads ...string) error
	GetWhere(where map[string]interface{}, preloads ...string) error
	GetDeleted(where map[string]interface{}) error

	FindAllWhere(where map[string]interface{}, preloads ...string) error
	FindAllUnscoped(preloads ...string) error

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
	return queryDb.Error
}

func (r *repositoryImpl) GetWhere(where map[string]interface{}, preloads ...string) error {
	queryDb := r.checkPreloads(preloads...)
	queryDb = r.applyConditions(queryDb, where).First(r.entity)
	return queryDb.Error
}

func (r *repositoryImpl) GetDeleted(where map[string]interface{}) error {
	queryDb := r.db.Unscoped().Where(where).First(r.entity)
	return queryDb.Error
}

func (r *repositoryImpl) FindAllWhere(where map[string]interface{}, preloads ...string) error {
	queryDb := r.checkPreloads(preloads...)
	queryDb = r.applyConditions(queryDb, where).Find(r.entity)
	return queryDb.Error
}

func (r *repositoryImpl) FindAllUnscoped(preloads ...string) error {
	queryDB := r.checkPreloads(preloads...)
	queryDB = queryDB.Unscoped().Find(r.entity) // ignora o filtro deleted_at
	return queryDB.Error
}

func (r *repositoryImpl) Create() error {
	queryDb := r.db.Create(r.entity)
	if queryDb.Error != nil {
		return queryDb.Error
	}
	if queryDb.RowsAffected == 0 {
		return fmt.Errorf("rows affected equals zero")
	}
	return nil
}

func (r *repositoryImpl) Save() error {
	queryDb := r.db.Save(r.entity)
	if queryDb.Error != nil {
		return queryDb.Error
	}
	if queryDb.RowsAffected == 0 {
		return fmt.Errorf("rows affected equals zero")
	}
	return nil
}

func (r *repositoryImpl) Delete(id uuid.UUID) error {
	queryDb := r.db.Delete(r.entity, id)
	if queryDb.Error != nil {
		return queryDb.Error
	}
	if queryDb.RowsAffected == 0 {
		return fmt.Errorf("rows affected equals zero")
	}
	return nil
}

func (r *repositoryImpl) Update(fields interface{}, id uuid.UUID) error {
	queryDb := r.db.Model(r.entity).Where("id = ?", id).Updates(fields).Find(r.entity)
	if queryDb.Error != nil {
		return queryDb.Error
	}
	if queryDb.RowsAffected == 0 {
		return fmt.Errorf("rows affected equals zero")
	}
	return nil
}

func (r *repositoryImpl) checkPreloads(preloads ...string) *gorm.DB {
	queryDb := r.db
	for _, preload := range preloads {
		queryDb = queryDb.Preload(preload)
	}
	return queryDb
}

func (r *repositoryImpl) applyConditions(db *gorm.DB, conditions map[string]interface{}) *gorm.DB {
	for key, value := range conditions {
		if strings.ContainsAny(key, "<>!=") {
			db = db.Where(fmt.Sprintf("%s ?", key), value)
		} else {
			db = db.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	return db
}
