package dao

import (
	"awesome-shortLink/tools"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ShortLink struct {
	Id    int64  `gorm:"primaryKey"`
	Long  string `gorm:"type:varchar(100);uniqueIndex:idx_long_short;index:idx_short_long, priority:2"`
	Short string `gorm:"type:varchar(10);uniqueIndex:idx_long_short;index:idx_short_long, priority:1"`
	Ctime int64
	Utime int64
}

var (
	ErrDuplicatedLongURL = gorm.ErrDuplicatedKey
	ErrNotFound          = gorm.ErrRecordNotFound
)

type ShortLinkDAO interface {
	Insert(ctx context.Context, longURL string) (ShortLink, error)
	FindByShort(ctx context.Context, shortURL string) (ShortLink, error)
}

type ShortLinkDAOV1 struct {
	db *gorm.DB
}

func (dao *ShortLinkDAOV1) FindByShort(ctx context.Context, shortURL string) (ShortLink, error) {
	var sl ShortLink
	err := dao.db.WithContext(ctx).Where("`short` = ?", shortURL).First(&sl).Error
	return sl, err
}

func (dao *ShortLinkDAOV1) Insert(ctx context.Context, url string) (ShortLink, error) {
	now := time.Now().UnixMilli()
	sl := ShortLink{
		Long:  url,
		Ctime: now,
		Utime: now,
	}
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 看看有木有
		err := tx.WithContext(ctx).Select("short").Where("`long` = ?", sl.Long).First(&sl).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			// 没有就插入
			err = tx.WithContext(ctx).Create(&sl).Error
			if err != nil {
				return err
			}
			sl.Short = tools.Encode(sl.Id)
			return tx.WithContext(ctx).Model(&ShortLink{}).Where("`long` = ?", sl.Long).Updates(map[string]any{
				"short": sl.Short,
			}).Error
		default:
			return err
		}
	})
	fmt.Println("")
	return sl, err
}

func NewShortLinkDAOV1(db *gorm.DB) ShortLinkDAO {
	return &ShortLinkDAOV1{
		db: db,
	}
}
