package errorHandler

import (
	"creatif-sdk-seed/storage"
	"fmt"
)

func CompleteCleanup() {
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", PROJECT_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", MAP_VARIABLES))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", VARIABLE_MAP))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", LIST_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", LIST_VARIABLES_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", USERS_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", GROUPS_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", CONNECTIONS_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", VARIABLE_GROUPS_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", PUBLISHED_LISTS_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", PUBLISHED_MAPS_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", VERSION_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", FILE_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", PUBLISHED_FILES_TABLE))
	storage.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", PUBLISHED_GROUPS_TABLE))
}
