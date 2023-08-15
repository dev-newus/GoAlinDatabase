package GoAlinDatabase

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Host               string
	Port               string
	User               string
	Password           string
	Name               string
	MaxIdleConnections int
	MaxOpenConnections int
}

type TenantConfig map[string]Config

var TenantConnections map[string]*gorm.DB = make(map[string]*gorm.DB)

// tenantID = agency_app_service_schema
func DatabaseInit(tenantID string, tenantConfig Config) {
	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		tenantConfig.User,
		tenantConfig.Password,
		tenantConfig.Host,
		tenantConfig.Port,
		tenantConfig.Name,
	)

	database, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database for tenant %s: %s\n", tenantID, err)
		return
	}

	sqlDB, err := database.DB()
	if err != nil {
		fmt.Printf("Failed to set connection pool for tenant %s: %s\n", tenantID, err)
		return
	}
	sqlDB.SetMaxIdleConns(tenantConfig.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(tenantConfig.MaxOpenConnections)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	TenantConnections[tenantID] = database

	fmt.Printf("Connection to database for tenant %s success\n", tenantID)
}

func GetTenantDB(tenantID string, config Config) *gorm.DB {
	if tenantDB, ok := TenantConnections[tenantID]; ok {
		return tenantDB
	}

	DatabaseInit(tenantID, config)
	return TenantConnections[tenantID]
}

func CloseTenantDB(tenantID string) {
	if db, ok := TenantConnections[tenantID]; ok {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}
