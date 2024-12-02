package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/dao"
	"context"
)

var(
	ErrUserDuplicationEmail = dao.ErrUserDuplicationEmail
	ErrUserNotFound = dao.ErrUserNotFound
)
type UserRepository struct {
	dao *dao.UserDAO
}
func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
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
func (r *UserRepository) FindById(int64)  {

}