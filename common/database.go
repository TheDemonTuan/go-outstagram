package common

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DBConn *gorm.DB

func ConnectDB() {
	gormConfig := &gorm.Config{
		PrepareStmt:            false,
		SkipDefaultTransaction: true,
	}

	if os.Getenv("APP_ENV") == "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	dbConn, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), gormConfig)

	if err != nil {
		panic("Database connection failed")
	}

	DBConn = dbConn
	defer runMigrate()
}

func runMigrate() {
	if os.Getenv("APP_ENV") == "development" {
		//if err := DBConn.Migrator().DropTable(&entity.User{}, &entity.Token{}, &entity.PostSave{}, &entity.Friend{}, &entity.Inbox{}, &entity.InboxFile{}, &entity.Post{}, &entity.PostComment{}, &entity.PostFile{}, &entity.PostLike{}); err != nil {
		//	panic(err)
		//}
		//if err := DBConn.AutoMigrate(&entity.User{}, &entity.PostSave{}, &entity.Token{}, &entity.Friend{}, &entity.Inbox{}, &entity.InboxFile{}, &entity.Post{}, &entity.PostComment{}, &entity.PostFile{}, &entity.PostLike{}, &entity.CommentLike{}, &entity.Report{}, &entity.Otp{}); err != nil {
		//	panic(err)
		//}

		log.Println("Success to migrate")
	}
}
