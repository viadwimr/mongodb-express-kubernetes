package models

import (
	"errors"
	"html"
	"strings"
	//"time"

	"github.com/jinzhu/gorm"
)

type Class struct {
	ID_2        uint64    `gorm:"primary_key;auto_increment" json:"id_2"`
	Theme     string    `gorm:"size:255;not null;unique" json:"theme"`
	Chapter   string    `gorm:"size:255;not null;" json:"chapter"`
	Learner    User      `json:"learner"`
	LearnerID  uint32    `gorm:"not null" json:"learner_id"`
	//CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	//UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Class) Prepare() {
	c.ID_2 = 0
	c.Theme = html.EscapeString(strings.TrimSpace(c.Theme))
	c.Chapter = html.EscapeString(strings.TrimSpace(c.Chapter))
	c.Learner = User{}
	//c.CreatedAt = time.Now()
	//c.UpdatedAt = time.Now()
}

func (c *Class) Validate() error {

	if c.Theme == "" {
		return errors.New("Required Theme")
	}
	if c.Chapter == "" {
		return errors.New("Required Chapter")
	}
	if c.LearnerID < 1 {
	 	return errors.New("Required Learner")
	}
	return nil
}

func (c *Class) SaveClass(db *gorm.DB) (*Class, error) {
	var err error
	err = db.Debug().Model(&Class{}).Create(&c).Error
	if err != nil {
		return &Class{}, err
	}
	if c.ID_2 != 0 {
		err = db.Debug().Model(&User{}).Where("id_2 = ?", c.LearnerID).Take(&c.Learner).Error
		if err != nil {
			return &Class{}, err
		}
	}
	return c, nil
}

func (c *Class) FindAllClasses(db *gorm.DB) (*[]Class, error) {
	var err error
	classes := []Class{}
	err = db.Debug().Model(&Class{}).Limit(100).Find(&classes).Error
	if err != nil {
		return &[]Class{}, err
	}
	if len(classes) > 0 {
		for i, _ := range classes {
			err := db.Debug().Model(&User{}).Where("id_2 = ?", classes[i].LearnerID).Take(&classes[i].Learner).Error
			if err != nil {
				return &[]Class{}, err
			}
		}
	}
	return &classes, nil
}

func (c *Class) FindClassByID(db *gorm.DB, cid uint64) (*Class, error) {
	var err error
	err = db.Debug().Model(&Class{}).Where("id_2 = ?", cid).Take(&c).Error
	if err != nil {
		return &Class{}, err
	}
	if c.ID_2 != 0 {
		err = db.Debug().Model(&User{}).Where("id_2 = ?", c.LearnerID).Take(&c.Learner).Error
		if err != nil {
			return &Class{}, err
		}
	}
	return c, nil
}

func (c *Class) UpdateAClass(db *gorm.DB) (*Class, error) {

	var err error

	err = db.Debug().Model(&Class{}).Where("id_2 = ?", c.ID_2).Updates(Class{Theme: c.Theme, Chapter: c.Chapter}).Error
	if err != nil {
		return &Class{}, err
	}
	if c.ID_2 != 0 {
		err = db.Debug().Model(&User{}).Where("id_2 = ?", c.LearnerID).Take(&c.Learner).Error
		if err != nil {
			return &Class{}, err
		}
	}
	return c, nil
}

func (c *Class) DeleteAClass(db *gorm.DB, cid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Class{}).Where("id_2 = ? and learner_id = ?", cid, uid).Take(&Class{}).Delete(&Class{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Class not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}