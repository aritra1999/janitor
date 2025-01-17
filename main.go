package main

import (
	"janitor/db"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func loadEnvFromFile() error {
	log.Print("Loading environment variables from file")
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	return nil
}

func init() {
	if err := loadEnvFromFile(); err != nil {
		if len(os.Getenv("SECRET_XATA_PG_ENDPOINT")) == 0 {
			log.Fatal().Err(err).Msg("Error loading environment variables")
			os.Exit(1)
		}
	}
}

func main() {
	dbConn, err := db.ConnectDB(os.Getenv("SECRET_XATA_PG_ENDPOINT"))
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}

	expiredUptimeChecks, err := db.GetExpiredUptimeChecks(dbConn)
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting expired uptime checks")
	}
	log.Printf("Expired uptime checks: %v", len(expiredUptimeChecks))

	count, err := db.DeleteUptimeChecksBatch(dbConn, expiredUptimeChecks)
	if err != nil {
		log.Fatal().Err(err).Msg("Error deleting uptime checks")
	}
	log.Printf("Deleted %v expired uptime checks", count)

	defer dbConn.Close()
}
