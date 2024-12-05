package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicationEmail = repository.ErrUserDuplicationEmail
var ErrInvalidUserOrPassword = errors.New("账号或密码不对")
type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err;
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		//DEBUG
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	//加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err;
	}

	u.Password = string(hash)
	//存起来
	return svc.repo.Create(ctx, u)
}
func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error)  {
	u, err := svc.repo.FindById(ctx, id)

	// 没有这个数据
	if err == cache.ErrKeyNotExist {
		return domain.User{}, err
	}
	return u, nil
}
