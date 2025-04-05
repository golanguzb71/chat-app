package database

import (
	"chat-app/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//err = db.AutoMigrate(
	//	&model.User{},
	//	&model.Message{},
	//	&model.Group{},
	//	&model.GroupMember{},
	//	&model.BlockedUser{},
	//)
	//if err != nil {
	//	return nil, err
	//}

	return db, nil
}
