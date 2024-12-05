package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicationEmail = errors.New("邮箱冲突")
	ErrUserNotFound = gorm.ErrRecordNotFound
)
type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062;
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱冲突
			return ErrUserDuplicationEmail
		}
	}
	return err
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	Email string `gorm:"unique"`
	Password string

	//创建时间，毫秒数
	Ctime int64
	//更新时间，毫秒数
	Utime int64
}
