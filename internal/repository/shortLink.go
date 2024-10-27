package repository

import (
	"awesome-shortLink/internal/domain"
	"awesome-shortLink/internal/repository/cache"
	"awesome-shortLink/internal/repository/dao"
	"awesome-shortLink/internal/repository/filter"
	"awesome-shortLink/tools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrDuplicatedLongURL = dao.ErrDuplicatedLongURL
	ErrNotFound          = dao.ErrNotFound
)

const (
	bloomFilterKey = "bloomFilter:shortURL"
	countKey       = "CntIdKey"
)

type ShortLinkRepository interface {
	Create(ctx context.Context, longURL string) (domain.ShortLink, error)
	FindByShort(ctx context.Context, shortURL string) (domain.ShortLink, error)
}

type ShortLinkRepositoryV1 struct {
	dao dao.ShortLinkDAO
	l   *zap.Logger
}

func (repo *ShortLinkRepositoryV1) FindByShort(ctx context.Context, shortURL string) (domain.ShortLink, error) {
	sl, err := repo.dao.FindByShort(ctx, shortURL)
	if err != nil {
		return domain.ShortLink{}, err
	}
	return entityToDomain(sl), nil
}

func (repo *ShortLinkRepositoryV1) Create(ctx context.Context, url string) (domain.ShortLink, error) {
	sl, err := repo.dao.InsertV1(ctx, url)
	if err != nil {
		return domain.ShortLink{}, err
	}
	return entityToDomain(sl), nil
}

func NewShortLinkRepositoryV1(dao dao.ShortLinkDAO, l *zap.Logger) ShortLinkRepository {
	return &ShortLinkRepositoryV1{
		dao: dao,
		l:   l,
	}
}

type ShortLinkRepositoryV2 struct {
	dao    dao.ShortLinkDAO
	l      *zap.Logger
	cache  cache.ShortLinkCache
	filter filter.BloomFilter
}

func (repo *ShortLinkRepositoryV2) Create(ctx context.Context, longURL string) (domain.ShortLink, error) {
	slStr, err := repo.cache.Get(ctx, getLongKey(longURL))
	switch {
	case errors.Is(err, redis.Nil):
		// 缓存中没有数据
		id, err := repo.cache.Incr(ctx, countKey)
		if err != nil {
			repo.l.Error("ID发号器失败", zap.Error(err))
			return domain.ShortLink{}, err
		}
		sl, err := repo.dao.InsertV2(ctx, dao.ShortLink{Long: longURL, Short: tools.Encode(id)})
		go func() {
			// 同步缓存
			var err error
			slByte, err := json.Marshal(sl)
			if err != nil {
				repo.l.Error("结构体序列化失败", zap.Error(err))
				return
			}
			go func() {
				err := repo.cache.Set(ctx, getLongKey(sl.Long), string(slByte))
				if err != nil {
					repo.l.Error("缓存同步失败", zap.Error(err))
					return
				}
			}()
			go func() {
				err := repo.cache.Set(ctx, getShortKey(sl.Short), string(slByte))
				if err != nil {
					repo.l.Error("缓存同步失败", zap.Error(err))
					return
				}
			}()
		}()
		go func() {
			// 存入布隆过滤器
			err := repo.filter.BFAdd(ctx, bloomFilterKey, GetFilterVal(sl.Short))
			if err != nil {
				repo.l.Error("布隆过滤器存放失败", zap.Error(err))
			}
		}()
		return entityToDomain(sl), err
	case err == nil:
		// 缓存中有数据
		var sl dao.ShortLink
		err = json.Unmarshal([]byte(slStr), &sl)
		if err != nil {
			repo.l.Error("反序列化失败", zap.Error(err))
			return domain.ShortLink{}, err
		}
		return entityToDomain(sl), nil
	default:
		return domain.ShortLink{}, err
	}
}

func (repo *ShortLinkRepositoryV2) FindByShort(ctx context.Context, shortURL string) (domain.ShortLink, error) {
	ok := repo.filter.BFExists(ctx, bloomFilterKey, GetFilterVal(shortURL))
	if !ok {
		return domain.ShortLink{}, ErrNotFound
	}
	slStr, err := repo.cache.Get(ctx, getShortKey(shortURL))
	switch {
	case errors.Is(err, redis.Nil):
		// 缓存没有，进一步查找数据库
		sl, err := repo.dao.FindByShort(ctx, shortURL)
		go func() {
			// 添加缓存
			err := repo.cache.Set(ctx, getShortKey(sl.Short), sl.Long)
			if err != nil {
				repo.l.Error("缓存同步失败", zap.Error(err))
			}
		}()
		return entityToDomain(sl), err
	case err != nil:
		return domain.ShortLink{}, err
	default:
		var sl dao.ShortLink
		err = json.Unmarshal([]byte(slStr), &sl)
		if err != nil {
			repo.l.Error("unmarshal失败", zap.Error(err))
			return domain.ShortLink{}, err
		}
		return entityToDomain(sl), nil
	}
}

func NewShortLinkRepositoryV2(dao dao.ShortLinkDAO, l *zap.Logger,
	cache cache.ShortLinkCache, filter filter.BloomFilter) ShortLinkRepository {
	return &ShortLinkRepositoryV2{
		dao:    dao,
		l:      l,
		cache:  cache,
		filter: filter,
	}
}

func entityToDomain(sl dao.ShortLink) domain.ShortLink {
	return domain.ShortLink{
		Short: sl.Short,
		Long:  sl.Long,
	}
}

// GetFilterVal 生成布隆过滤器中的val字段值
func GetFilterVal(s string) string {
	return fmt.Sprintf("shortLink:filter:%s", s)
}

// getShortKey 生成Redis中shortURL对应的key，value为dao.shortLink
func getShortKey(s string) string {
	return fmt.Sprintf("shortLink:short:%s", s)
}

// getLongKey 生成Redis中longURL对应的key，value为dao.shortLink
func getLongKey(s string) string {
	return fmt.Sprintf("shortLink:long:%s", s)
}
