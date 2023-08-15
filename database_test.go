package GoAlinDatabase

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	connectionSemaphore chan struct{}
)

func init() {
	connectionSemaphore = make(chan struct{}, 5)
}

func TestDatabaseInit(t *testing.T) {
	testConfig := Config{
		Host:               "localhost",
		Port:               "3306",
		User:               "root",
		Password:           "",
		Name:               "base",
		MaxIdleConnections: 10,
		MaxOpenConnections: 20,
	}

	DatabaseInit("AGENCY1", testConfig)

	tenantDB := TenantConnections["AGENCY1"]
	assert.NotNil(t, tenantDB, "Tenant database connection should not be nil")

	CloseTenantDB("AGENCY1")
}

func TestGetTenantDB(t *testing.T) {
	testConfig := Config{
		Host:               "localhost",
		Port:               "3306",
		User:               "root",
		Password:           "",
		Name:               "base",
		MaxIdleConnections: 10,
		MaxOpenConnections: 20,
	}

	tenantDB1 := GetTenantDB("AGENCY2", testConfig)
	tenantDB2 := GetTenantDB("AGENCY2", testConfig)

	assert.NotNil(t, tenantDB1, "Tenant database connection should not be nil")
	assert.NotNil(t, tenantDB2, "Tenant database connection should not be nil")
	assert.Equal(t, tenantDB1, tenantDB2, "Connections should be the same")

	CloseTenantDB("AGENCY2")
}

func TestCloseTenantDB(t *testing.T) {
	testConfig := Config{
		Host:               "localhost",
		Port:               "3306",
		User:               "root",
		Password:           "",
		Name:               "base",
		MaxIdleConnections: 10,
		MaxOpenConnections: 20,
	}

	DatabaseInit("AGENCY3", testConfig)

	CloseTenantDB("AGENCY3")

	// Try to retrieve the tenant connection from the map
	// tenantDB := TenantConnections["AGENCY3"]
	// assert.Nil(t, tenantDB, "Tenant database connection should be nil after closing")
}

func TestConnectionPooling(t *testing.T) {
	testConfig := Config{
		Host:               "localhost",
		Port:               "3306",
		User:               "root",
		Password:           "",
		Name:               "base",
		MaxIdleConnections: 2,
		MaxOpenConnections: 5,
	}

	DatabaseInit("AGENCY4", testConfig)
	queryCounter := 1

	wg := sync.WaitGroup{}
	for i := 1; i <= 10000; i++ {
		wg.Add(1)
		go func(queryNumber int) {
			connectionSemaphore <- struct{}{}
			defer func() {
				<-connectionSemaphore
				wg.Done()
			}()
			db := GetTenantDB("AGENCY4", testConfig)
			var result int
			err := db.Raw("SELECT 1").Scan(&result).Error
			if err != nil {
				fmt.Printf("Query %d: Error executing query: %s\n", queryNumber, err)
			} else {
				fmt.Printf("Query %d: Result: %d\n", queryNumber, result)
			}
			// PrintConnectionPoolStats(db)
		}(queryCounter)

		queryCounter++
	}

	time.Sleep(1 * time.Second)

	CloseTenantDB("AGENCY4")
}

func PrintConnectionPoolStats(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Error getting DB stats: %s\n", err)
		return
	}

	fmt.Printf("Max Open Connections: %d\n", sqlDB.Stats().MaxOpenConnections)
	fmt.Printf("Open Connections: %d\n", sqlDB.Stats().OpenConnections)
	fmt.Printf("In Use: %d\n", sqlDB.Stats().InUse)
	fmt.Printf("Idle: %d\n", sqlDB.Stats().Idle)
	fmt.Printf("Wait Count: %d\n", sqlDB.Stats().WaitCount)
	fmt.Printf("Wait Duration: %s\n", sqlDB.Stats().WaitDuration)
}
