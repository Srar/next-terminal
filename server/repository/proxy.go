package repository

import (
	"next-terminal/server/model"

	"gorm.io/gorm"
)

type ProxyRepository struct {
	DB *gorm.DB
}

func NewProxyRepository(db *gorm.DB) *ProxyRepository {
	proxyRepository = &ProxyRepository{DB: db}
	return proxyRepository
}

func (p ProxyRepository) FindById(id string) (o model.Proxy, err error) {
	err = p.DB.Where("id = ?", id).First(&o).Error
	return
}

func (p ProxyRepository) Find(pageIndex, pageSize int, name, proxyType, order, field string) (o []model.Proxy, total int64, err error) {
	t := model.Proxy{}
	db := p.DB.Table(t.TableName())
	dbCounter := p.DB.Table(t.TableName())

	if len(name) > 0 {
		db = db.Where("name like ?", "%"+name+"%")
		dbCounter = dbCounter.Where("name like ?", "%"+name+"%")
	}

	if len(proxyType) > 0 {
		db = db.Where("type = ?", proxyType)
		dbCounter = dbCounter.Where("type = ?", proxyType)
	}

	err = dbCounter.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if order == "ascend" {
		order = "asc"
	} else {
		order = "desc"
	}

	if field == "name" {
		field = "name"
	} else {
		field = "created"
	}

	err = db.Order(field + " " + order).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&o).Error
	if o == nil {
		o = make([]model.Proxy, 0)
	}

	return
}

func (p ProxyRepository) UpdateById(o *model.Proxy, id string) error {
	o.ID = id
	return p.DB.Updates(o).Error
}

func (p ProxyRepository) DeleteByID(id string) error {
	return p.DB.Delete(&model.Proxy{ID: id}).Error
}

func (p ProxyRepository) Create(o *model.Proxy) error {
	return p.DB.Create(o).Error
}
