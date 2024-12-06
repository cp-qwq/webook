package service

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany
type CodeService struct {
	repo *repository.CodeRepository
	smsSvc sms.Service
}

//发验证码
func (svc *CodeService) Send(ctx context.Context,
	// 区别业务场景
	biz string,
	phone string) error {
	// 生成一个验证码
	code := svc.generateCode()

	//塞进 redis
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		// 有问题
		return err
	}
	//发送
	const codeTplId = "xxxxx"
	return svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
}

func (svc *CodeService) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	return ok, err
}

func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
