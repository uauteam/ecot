package repo

import (
	"github.com/uauteam/ecot/dto/qry"
	"github.com/uauteam/ecot/entity"
)

func Create(e entity.Entity) (err error) {
	db := DB(e.DBName())
	if err = db.Error; err != nil {
		return
	}

	if err = db.Create(e).Error; err != nil {
		return
	}

	return
}

func Get(id uint, e entity.Entity) (err error) {
	db := DB(e.DBName())
	if err = db.Error; err != nil {
		return
	}

	if err = db.First(e, id).Error; err != nil {
		return
	}

	return
}

func Find(e entity.Entity, results interface{}) (err error) {
	db := DB(e.DBName())
	if err = db.Error; err != nil {
		return
	}

	if err = db.Where(e).Find(results).Error; err != nil {
		return
	}

	return
}

func FindPage(e entity.Entity, pageQuery qry.PageQuery, results interface{}) (total uint, err error) {
	db := DB(e.DBName())
	if err = db.Error; err != nil {
		return
	}

	if err = db.Model(e).Where(e).Count(&total).Error; err != nil {
		return
	}

	db = db.Offset(pageQuery.Page * pageQuery.Size).Limit(pageQuery.Size)
	if err = db.Where(e).Find(results).Error; err != nil {
		return
	}

	return
}

func Update(id uint, e entity.Entity) (err error) {
	db := DB(e.DBName())
	if err = db.Error; err != nil {
		return
	}

	if err = db.Model(e).Omit(e.ProtectedFields()...).Where("id = ?", id).Updates(e).Error; err != nil {
		return
	}

	return
}
