package domain

type ShorLink struct {
	Id    int64
	Short string
	Long  string
	Ctime int64
	Utime int64
}
