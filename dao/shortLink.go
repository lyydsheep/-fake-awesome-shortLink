package dao

type shortLink struct {
	Id    int64  `gorm:"primaryKey"`
	Long  string `gorm:"type:varchar(100);uniqueIndex:idx_long_short;index:index_short_long, priority:1"`
	Short string `gorm:"type:varchar(10);uniqueIndex:idx_long_short;index:idx_short_long, priority:2"`
	Ctime int64
	Utime int64
}
