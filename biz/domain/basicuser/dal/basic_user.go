package dal

import (
	"context"
	"errors"

	"github.com/xh-polaris/synapse4b/biz/pkg/errorx"
	"github.com/xh-polaris/synapse4b/biz/types/errno"

	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/dal/model"
	"github.com/xh-polaris/synapse4b/biz/domain/basicuser/dal/query"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/id"
	"github.com/xh-polaris/synapse4b/biz/infra/contract/orm"
	"gorm.io/gorm"
)

func NewBasicUserDAO(db *orm.DB) *BasicUserDAO {
	return &BasicUserDAO{query: query.Use(db)}
}

type BasicUserDAO struct {
	query *query.Query
}

func (d *BasicUserDAO) FindByID(ctx context.Context, uid string) (*model.BasicUser, error) {
	buid, err := id.FromHex(uid)
	if err != nil {
		return nil, err
	}
	user, err := d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.ID.Eq(buid)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *BasicUserDAO) FindByPhone(ctx context.Context, phone string) (*model.BasicUser, error) {
	user, err := d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.Phone.Eq(phone)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, err
}

func (d *BasicUserDAO) FindByPhoneAndUnit(ctx context.Context, phone, unitId string) (*model.BasicUser, error) {
	uid, err := id.FromHex(unitId)
	if err != nil {
		return nil, errorx.New(errno.InvalidParameter, errorx.KV("parameter", "单位id"))
	}
	user, err := d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.Phone.Eq(phone), d.query.BasicUser.UnitID.Eq(uid)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, err
}

func (d *BasicUserDAO) FindByEmail(ctx context.Context, email string) (*model.BasicUser, error) {
	user, err := d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.Email.Eq(email)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, err
}

func (d *BasicUserDAO) FindByEmailAndUnit(ctx context.Context, email, unitId string) (*model.BasicUser, error) {
	uid, err := id.FromHex(unitId)
	if err != nil {
		return nil, errorx.New(errno.InvalidParameter, errorx.KV("parameter", "单位id"))
	}

	user, err := d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.Email.Eq(email), d.query.BasicUser.UnitID.Eq(uid)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, err
}

func (d *BasicUserDAO) FindManyByUnitID(ctx context.Context, unitId string) ([]*model.BasicUser, error) {
	uid, err := id.FromHex(unitId)
	if err != nil {
		return nil, errorx.New(errno.InvalidParameter, errorx.KV("parameter", "单位id"))
	}
	return d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.UnitID.Eq(uid)).Find()
}

func (d *BasicUserDAO) FindByCode(ctx context.Context, unitId, code string) (*model.BasicUser, error) {
	uid, err := id.FromHex(unitId)
	if err != nil {
		return nil, errorx.New(errno.InvalidParameter, errorx.KV("parameter", "单位id"))
	}
	user, err := d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.UnitID.Eq(uid), d.query.BasicUser.Code.Eq(code)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, err
}

func (d *BasicUserDAO) Create(ctx context.Context, nu *model.BasicUser) (*model.BasicUser, error) {
	err := d.query.WithContext(ctx).BasicUser.Create(nu)
	return nu, err
}

func (d *BasicUserDAO) ResetPassword(ctx context.Context, basicUserId, password string) error {
	buid, err := id.FromHex(basicUserId)
	if err != nil {
		return err
	}
	_, err = d.query.WithContext(ctx).BasicUser.Where(d.query.BasicUser.ID.Eq(buid)).
		Updates(map[string]interface{}{
			d.query.BasicUser.Password.ColumnName().String(): password,
			d.query.BasicUser.Encrypt.ColumnName().String():  0,
		})
	return err
}

// FindCompletely 完全匹配
// 原User和新User的所有非空字段一致视为完全匹配
func (d *BasicUserDAO) FindCompletely(ctx context.Context, unitID, code, phone, email string) (*model.BasicUser, error) {
	uid, err := id.FromHex(unitID)
	if err != nil {
		return nil, errorx.New(errno.InvalidParameter, errorx.KV("parameter", "单位id"))
	}

	// code, phone, email至少有一个非空，否则无法匹配
	if code == "" && phone == "" && email == "" {
		return nil, errorx.New(errno.MissingParameter, errorx.KV("parameter", "code、phone、email至少需要一个"))
	}

	var user model.BasicUser
	err = d.query.WithContext(ctx).BasicUser.UnderlyingDB().
		Model(&model.BasicUser{}).
		Where("unit_id = ?", uid).
		Where("(code IS NULL OR code = '' OR ? = '' OR code = ?)", code, code).
		Where("(phone IS NULL OR phone = '' OR ? = '' OR phone = ?)", phone, phone).
		Where("(email IS NULL OR email = '' OR ? = '' OR email = ?)", email, email).
		Where("((? <> '' AND code IS NOT NULL AND code <> '' AND code = ?) OR (? <> '' AND phone IS NOT NULL AND phone <> '' AND phone = ?) OR (? <> '' AND email IS NOT NULL AND email <> '' AND email = ?))", code, code, phone, phone, email, email).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindPartly 部分匹配
// 匹配到unitId和code,phone,email三组中任意一组
func (d *BasicUserDAO) FindPartly(ctx context.Context, unitID, code, phone, email string) (*model.BasicUser, error) {
	uid, err := id.FromHex(unitID)
	if err != nil {
		return nil, errorx.New(errno.InvalidParameter, errorx.KV("parameter", "单位id"))
	}

	if code != "" {
		user, err := d.query.WithContext(ctx).BasicUser.Where(
			d.query.BasicUser.UnitID.Eq(uid),
			d.query.BasicUser.Code.Eq(code),
		).First()
		if err == nil {
			return user, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	if phone != "" {
		user, err := d.query.WithContext(ctx).BasicUser.Where(
			d.query.BasicUser.UnitID.Eq(uid),
			d.query.BasicUser.Phone.Eq(phone),
		).First()
		if err == nil {
			return user, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	if email != "" {
		user, err := d.query.WithContext(ctx).BasicUser.Where(
			d.query.BasicUser.UnitID.Eq(uid),
			d.query.BasicUser.Email.Eq(email),
		).First()
		if err == nil {
			return user, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return nil, nil
}
