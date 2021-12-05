package main

import "github.com/jinzhu/gorm"

func connectToSql(database *gorm.DB, databaseERR error) {
	if databaseERR != nil {
		panic(databaseERR)
	}
	defer func(database *gorm.DB) {
		err := database.Close()
		if err != nil {

		}
	}(database)
	database.SingularTable(true)
	database.InstantSet("gorm:table_options", "ENGINE=InnoDB")
	// 刷新数据库中的表格定义，使其保持最新（只增不减）
	// 创建（新的）表、缺少的外键、约束、列和索引，并且会更改现有列的类型
	database.AutoMigrate(&User{})
	database.AutoMigrate(&Salesman{})
	database.AutoMigrate(&Repairman{})
	database.AutoMigrate(&TypeOverview{})
	database.AutoMigrate(&PartsOverview{})
	database.AutoMigrate(&Vehicle{})
	database.AutoMigrate(&Attorney{})
	database.AutoMigrate(&Arrangement{})
	database.AutoMigrate(&AuthSession{})
	database.AutoMigrate(&Notification{})

	// 构建表格
	database.Model(&Vehicle{}).AddForeignKey("user_id", "user(number)", "RESTRICT", "RESTRICT")
	database.Model(&Attorney{}).AddForeignKey("user_id", "user(number)", "RESTRICT", "RESTRICT")
	database.Model(&Attorney{}).AddForeignKey("vehicle_number", "vehicle(number)", "RESTRICT", "RESTRICT")
	database.Model(&Attorney{}).AddForeignKey("salesman_id", "salesman(number)", "RESTRICT", "RESTRICT")
	database.Model(&Arrangement{}).AddForeignKey("order_number", "attorney(number)", "RESTRICT", "RESTRICT")
	database.Model(&Arrangement{}).AddForeignKey("repairman_number", "repairman(number)", "RESTRICT", "RESTRICT")
	database.Model(&Arrangement{}).AddForeignKey("project_number", "type_overview(project_number)", "RESTRICT", "RESTRICT")
	database.Model(&Arrangement{}).AddForeignKey("parts_number", "parts_overview(parts_number)", "RESTRICT", "RESTRICT")
	//database.Model(&Notification{}).AddForeignKey("user_id", "user(number)", "RESTRICT", "RESTRICT")
	//database.Model(&Notification{}).AddForeignKey("user_id", "repairman(number)", "RESTRICT", "RESTRICT")
	//database.Model(&Notification{}).AddForeignKey("user_id", "salesman(number)", "RESTRICT", "RESTRICT")
}
