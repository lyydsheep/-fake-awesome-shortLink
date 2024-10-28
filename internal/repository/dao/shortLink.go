package dao

import (
	"awesome-shortLink/tools"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ShortLink struct {
	Id    int64  `gorm:"primaryKey" json:"id"`
	Long  string `gorm:"type:varchar(100);uniqueIndex:idx_long_short;index:idx_short_long, priority:2" json:"long"`
	Short string `gorm:"type:varchar(30);uniqueIndex:idx_long_short;index:idx_short_long, priority:1" json:"short"`
	Ctime int64  `json:"ctime"`
	Utime int64  `json:"utime"`
}

var (
	ErrDuplicatedLongURL = gorm.ErrDuplicatedKey
	ErrNotFound          = gorm.ErrRecordNotFound
)

type ShortLinkDAO interface {
	InsertV1(ctx context.Context, longURL string) (ShortLink, error)
	InsertV2(ctx context.Context, sl ShortLink) (ShortLink, error)
	FindByShort(ctx context.Context, shortURL string) (ShortLink, error)
}

type ShortLinkDAOV1 struct {
	db *gorm.DB
}

func (dao *ShortLinkDAOV1) InsertV2(ctx context.Context, sl ShortLink) (ShortLink, error) {
	now := time.Now().UnixMilli()
	sl.Ctime = now
	sl.Utime = now
	return sl, dao.db.Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Select("`short`").Where("`long` = ?", sl.Long).First(&sl).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return tx.WithContext(ctx).Create(&sl).Error
		default:
			return err
		}
	})
}

func (dao *ShortLinkDAOV1) FindByShort(ctx context.Context, shortURL string) (ShortLink, error) {
	var sl ShortLink
	err := dao.db.WithContext(ctx).Where("`short` = ?", shortURL).First(&sl).Error
	return sl, err
}

// InsertV1 shortURL和MySQL耦合在一起，扩展性低
func (dao *ShortLinkDAOV1) InsertV1(ctx context.Context, url string) (ShortLink, error) {
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
	return sl, err
}

func NewShortLinkDAOV1(db *gorm.DB) ShortLinkDAO {
	return &ShortLinkDAOV1{
		db: db,
	}
}
