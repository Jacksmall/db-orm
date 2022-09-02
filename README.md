# dborm
基于gorm的db脚手架，独立出来，可以适用到不同的项目或者服务中
# 目录结构会如下:
```
db-orm
 |--models
      |--user.go
 |--database
      |--dborm.go
 user_test.go
 main.go
 go.mod
 ...
```
# Installing
$> go get github.com/Jacksmall/db-orm
# Example
```
package main

import (
	"fmt"
	"log"
	"time"

	"db-orm/database"
  "db-orm/models"
)

const webPort = "8070"

type Config struct {
	db *gorm.DB
}

func main() {
  app := new(Config)

	db, err := connectToDB()
	if err != nil {
		log.Panic(err)
	}

	app.db = db
  
  // 创建user表
  db.Migrator().CreateTable(&models.User{})
  
	// 脚手架db赋值
	database.SetDB(db)
  
  app.models = models.NewModels(db)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting server on port:%s", webPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connectToDB() (*gorm.DB, error) {
	return gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/gomicro?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
}
```

# models/user.go
```
package models

import (
  "log"
  "time"
  
  "gorm.io/gorm"
  "db-orm/database"
)

type Models struct {
	User User
}

func NewModels(gdb *gorm.DB) *Models {
	db = gdb

	return &Models{
		User: User{},
	}
}

type User struct {
	UID        UID    `gorm:"primaryKey;column:uid;not null;autoIncrement"`
	UnionID    string `gorm:"index:union,unique;column:unionId;not null"`
	OpenID     string `gorm:"column:openId;not null"`
	SupplierID uint   `gorm:"index:union,unique;column:supplierId;not null"`
	Channel    string `gorm:"column:channel;not null"`
	Nickname   string `gorm:"column:nickname;not null"`
	HeadPic    string `gorm:"column:headPic;not null"`
	Sex        uint8  `gorm:"column:sex;not null;default:0"`
	Email      string
	Mobile     string
	Birthday   time.Time
	WxAppID    string    `gorm:"column:wxAppId;not null"`
	WxUnionID  string    `gorm:"column:wxUnionId;not null"`
	IsFrozen   uint8     `gorm:"column:isFrozen;not null"`
	CreatedAt  time.Time `gorm:"column:createdAt;not null"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;not null;default:0;autoUpdateTime"`
}

func (u User) TableName() string {
	return "t_users"
}

type UserCreateReq struct {
	SupplierID uint      `json:"supplier_id"`
	Mobile     string    `json:"mobile"`
	UnionID    string    `json:"union_id"`
	Nickname   string    `json:"nickname,omitempty"`
	Channel    string    `json:"channel"`
	HeadPic    string    `json:"head_pic,omitempty"`
	Email      string    `json:"email,omitempty"`
	Sex        uint8     `json:"sex,omitempty"`
	Birthday   time.Time `json:"birthday,omitempty"`
}

// 创建单个用户
func (u User) Create(f UserCreateReq) (UID, error) {
	user := User{
		UnionID:    f.UnionID,
		SupplierID: f.SupplierID,
		Channel:    f.Channel,
		Nickname:   f.Nickname,
		HeadPic:    f.HeadPic,
		Sex:        f.Sex,
		Email:      f.Email,
		Mobile:     f.Mobile,
		Birthday:   f.Birthday,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := database.M(u, db).Insert(&user); err != nil {
		log.Fatalf("error creating user: %v", err)
	}
	return user.UID, nil
}
```

# user_test.go 
```
package main

import (
	"log"
	"reflect"
	"testing"
	"time"

	"db-orm/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var model *models.Models

func init() {
	gdb, err := gorm.Open(mysql.Open("root:chenkuanwo5@tcp(127.0.0.1:3306)/gomicro?charset=utf8&loc=Local&parseTime=True"))
	if err != nil {
		log.Fatal(err)
	}
	db = gdb
	model = models.NewModels(gdb)
}

func TestCreateUser(t *testing.T) {
	app := Config{
		db:     db,
		models: model,
	}
	req := models.UserCreateReq{
		SupplierID: 1,
		Mobile:     "15666667777",
		Email:      "666666@qq.com",
		UnionID:    "sdsad12sssqqq123fdsfdsfdsf23",
		Nickname:   "lck",
		Channel:    "app",
		HeadPic:    "http://example.com/22.png",
		Sex:        1,
		Birthday:   time.Now(),
	}

	real, err := app.models.User.Create(req)
	if real == 0 || err != nil {
		t.Errorf("CreateUser failed:%v", err)
	}
}

```

# Output
$> go run main.go
```
Starting server on port:8070
```

# Test Output
$> go test -run CreateUser
```
PASS
ok      github.com/Jacksmall/db-orm    0.373s
```


