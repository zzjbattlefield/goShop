package main

import (
	"crypto/md5"
	"encoding/hex"
	"goshop/user_srv/model"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	dsn := "root:root@tcp(192.168.58.130:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&model.User{})
	// options := &password.Options{16, 100, 32, sha512.New}
	// salt, encodedPwd := password.Encode("admin_123", options)
	// NewPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	// for i := 0; i < 10; i++ {
	// 	user := model.User{
	// 		NickName: fmt.Sprintf("测试人员%d", i),
	// 		Mobile:   fmt.Sprintf("1511259714%d", i),
	// 		Password: NewPassword,
	// 	}
	// 	db.Save(&user)
	// }
	// options := &password.Options{16, 100, 32, sha512.New}
	// salt, encodedPwd := password.Encode("generic password", options)
	// NewPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	// PasswordInfo := strings.Split(NewPassword, "$")
	// check := password.Verify("generic password", PasswordInfo[2], PasswordInfo[3], options)
	// fmt.Println(check) // true
}
