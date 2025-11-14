package testutils

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go_mcp_server/internal/model"

	_ "github.com/lib/pq"
)

var (
	testDBContainer *TestDBContainer
	once            sync.Once
)

var allModels = []interface{}{
	&model.Tool{},
}

type TestDBContainer struct {
	Container testcontainers.Container
	DB        *gorm.DB
	URI       string
}

func SetupTestDB(t *testing.T) *TestDBContainer {

	once.Do(func() {
		ctx := context.Background()

		req := testcontainers.ContainerRequest{
			Image:        "postgres:15",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
				return fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=testdb sslmode=disable", host, port.Port())
			}),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           nil,
		})
		if err != nil {
			t.Fatalf("failed to start postgres container: %v", err)
		}

		host, _ := container.Host(ctx)
		port, _ := container.MappedPort(ctx, "5432")

		dsn := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=testdb sslmode=disable", host, port.Port())
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to connect to DB: %v", err)
		}

		if err := db.AutoMigrate(allModels...); err != nil {
			t.Fatalf("failed to migrate DB schema: %v", err)
		}

		testDBContainer = &TestDBContainer{
			Container: container,
			DB:        db,
			URI:       dsn,
		}
	})

	// 모든 모델 테이블 초기화
	for _, model := range allModels {
		stmt := &gorm.Statement{DB: testDBContainer.DB}
		if err := stmt.Parse(model); err != nil {
			t.Fatalf("failed to parse model: %v", err)
		}
		tableName := stmt.Schema.Table
		if err := testDBContainer.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)).Error; err != nil {
			t.Fatalf("failed to truncate table %s: %v", tableName, err)
		}
	}

	return testDBContainer
}

func TestMain(m *testing.M) {
	code := m.Run()
	if testDBContainer != nil && testDBContainer.Container != nil {
		testDBContainer.Container.Terminate(context.Background())
	}
	os.Exit(code)
}
