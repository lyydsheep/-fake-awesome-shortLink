package dao

import (
	"awesome-shortLink/tools"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ShortLink struct {
	Id    int64  `gorm:"primaryKey"`
	Long  string `gorm:"type:varchar(100);uniqueIndex:idx_long;index:idx_short_long, priority:2"`
	Short string `gorm:"type:varchar(10);index:idx_short_long, priority:1"`
	Ctime int64
	Utime int64
}

var (
	ErrDuplicatedLongURL = gorm.ErrDuplicatedKey
)

type ShortLinkDAO interface {
	Insert(ctx context.Context, url string) (ShortLink, error)
}

type ShortLinkDAOV1 struct {
	db *gorm.DB
}

func (dao *ShortLinkDAOV1) Insert(ctx context.Context, url string) (ShortLink, error) {
	now := time.Now().UnixMilli()
	sl := ShortLink{
		Long:  url,
		Ctime: now,
		Utime: now,
	}
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&sl).Error
		if err != nil {
			return err
		}
		return tx.Model(&ShortLink{}).Updates(map[string]any{
			"short": tools.Encode(sl.Id),
		}).Error
	})
	if errors.Is(err, ErrDuplicatedLongURL) {

	}
}

func NewShortLinkDAOV1(db *gorm.DB) ShortLinkDAO {
	return &ShortLinkDAOV1{
		db: db,
	}
}
