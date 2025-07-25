package datasources

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// const (
//   host     = "interchange.proxy.rlwy.net"  // or the Docker service name if running in another container
//   port     = 25977         // default PostgreSQL port
//   user     = "postgres"     // as defined in docker-compose.yml
//   password = "oMclFKTHjJEFEJaDUhoIedagkygDAEBZ" // as defined in docker-compose.yml
//   dbname   = "railway" // as defined in docker-compose.yml
// )

func NewPostgresql() *sql.DB {
	host := os.Getenv("POSTGRESQL_HOST")
	port := os.Getenv("POSTGRESQL_PORT")
	user := os.Getenv("POSTGRESQL_USER")
	password := os.Getenv("POSTGRESQL_PASSWORD")
	dbname := os.Getenv("POSTGRESQL_DBNAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open PostgreSQL: %v", err)
	}

	// ตรวจสอบการเชื่อมต่อ
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL!")
	return db
}
