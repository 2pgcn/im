package data

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type User struct {
	Id         int64  `gorm:"column:id;primaryKey;type:bigint;index;comment:uid"`
	AppId      string `gorm:"column:app_id;type:varchar(32);comment:appid"`
	RoomId     string `gorm:"column:room_id;type:varchar(32);comment:群id"`
	CreateTime int64  `gorm:"column:create_time;type:bigint;not null"`
	UpdateTime int64  `gorm:"column:update_time;type:bigint;not null"`
	DeleteTime int64  `gorm:"column:delete_time;type:bigint;default:0;comment:删除时间"`
}

func (u *User) Encode() (string, error) {
	data, err := json.Marshal(u)
	return string(data), err
}

func (u *User) UnEncode(data string) error {
	return json.Unmarshal([]byte(data), &u)
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.UpdateTime == 0 {
		u.UpdateTime = time.Now().Unix()
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.CreateTime == 0 {
		u.CreateTime = time.Now().Unix()
	}
	if u.UpdateTime == 0 {
		u.UpdateTime = time.Now().Unix()
	}
	return nil
}

type UserRepo interface {
}

func (d *Data) FindByID(ctx context.Context, id int64) (*User, error) {
	var user *User
	err := d.mysqlClient.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &User{
		Id:         user.Id,
		AppId:      user.AppId,
		RoomId:     user.RoomId,
		CreateTime: user.CreateTime,
		DeleteTime: user.DeleteTime,
	}, nil
}

func (d *Data) FindByIDCache(ctx context.Context, id int64) (user *User, err error) {
	rst := d.redisClient.Get(ctx, string(id))
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	err = user.UnEncode(rst.String())
	return

}

func (d *Data) CreateTokenLinkId(ctx context.Context, id int64, token string) error {
	return d.redisClient.SetNX(ctx, getRedisUserTokenKey(token), id,
		time.Hour*3+time.Minute*time.Duration(rand.Int63())).Err()
}

func (d *Data) GetIdByToken(ctx context.Context, token string) (int64, error) {
	rst := d.redisClient.Get(ctx, getRedisUserTokenKey(token))
	id, err := rst.Int64()
	if err == redis.Nil {
		return -1, nil
	}
	return id, err
}

func getRedisUserTokenKey(token string) string {
	return "pg:auth:token:" + token
}
