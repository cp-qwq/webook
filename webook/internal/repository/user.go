package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
)

var(
	ErrUserDuplicationEmail = dao.ErrUserDuplicationEmail
	ErrUserNotFound = dao.ErrUserNotFound
)
type UserRepository struct {
	dao *dao.UserDAO
	cache *cache.UserCache
}
func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao: dao,
		cache: c,
	}
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id: u.Id,
		Email: u.Email,
		Password: u.Password,
	}, nil
}
func (r *UserRepository) Create(ctx context.Context, u domain.User) error  {
	return r.dao.Insert(ctx, dao.User{
		Email: u.Email,
		Password: u.Password,
	})
	//在这里操作缓存。
}
func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error)  {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 必然有数据
		return u, nil
	}

	// 去数据库取数据
	// 不过也可能是因为 redis 崩了，这里一定要对数据库做好保护（限流）
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id: ue.Id,
		Email: ue.Email,
		Password: ue.Password,
	}
	err = r.cache.Set(ctx, id, u)
	if err != nil {
		// 打个日志，做监控

	}
	return u, nil
}