package cli

import (
	"flag"
	"os"
)

// AOptions is the configuration options for stuff
type AOptions struct {
	Db        string
	Dbpath    string
	InputPath string
	IndexPath string
	Nosync    bool
	IP        string
	Port      string
}

func getEnvString(key string, fallback *string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = *fallback
	}
	return value
}

// GetArgs will load options from command line or environ
func GetArgs() AOptions {
	db := flag.String("db", "bolt", "Database Backend. (default \"bolt\")")
	dbpath := flag.String("dbpath", "/tmp/testdb", "Path to the database. (default \"/tmp/testdb\")")
	inputpath := flag.String("inputpath", "", "Path to input file (default \"\")")
	indexpath := flag.String("indexpath", "/tmp/index", "Path to the database. (default \"/tmp/index\")")
	nosync := flag.Bool("nosync", true, "Should db not sync to disk (default true)")
	ip := flag.String("ip", "127.0.0.1", "IP server should bind too (default 127.0.0.1)")
	port := flag.String("port", "8080", "Port server should bind too (default 8080)")

	flag.Parse()

	return AOptions{
		Db:        getEnvString("A_DB", db),
		Dbpath:    getEnvString("A_DBPATH", dbpath),
		InputPath: getEnvString("A_INPUTPATH", inputpath),
		Nosync:    *nosync,
		IP:        getEnvString("A_IP", ip),
		Port:      getEnvString("PORT", port),
		IndexPath: getEnvString("A_INDEXPATH", indexpath),
	}
}
