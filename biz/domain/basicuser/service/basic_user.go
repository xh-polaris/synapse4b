package service

import (
	"context"

	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/entity"
)

type BasicUser interface {
	UserIDExist(ctx context.Context, userId string) (*entity.BasicUser, bool, error)
	PhoneExist(ctx context.Context, phone string) (is bool, err error)
	CodeExist(ctx context.Context, schoolId, code string) (is bool, err error)
	EmailExist(ctx context.Context, email string) (is bool, err error)
	Register(ctx context.Context, authType, authId, extraAuthId, password string) (*entity.BasicUser, error)
	LoginByPhone(ctx context.Context, requirePassword bool, phone, verify string) (*entity.BasicUser, error)
	LoginByEmail(ctx context.Context, requirePassword bool, email, verify string) (*entity.BasicUser, error)
	LoginByCode(ctx context.Context, schoolId, code, verify string) (*entity.BasicUser, error)
	ResetPassword(ctx context.Context, basicUserId string, password string) error
	CreateBasicUser(ctx context.Context, unitId, code, phone, email, password string, encryptType int64) (*entity.BasicUser, error)
}
