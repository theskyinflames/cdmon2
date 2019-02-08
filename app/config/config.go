package config

import (
	"os"
	"strconv"
)

const (
	APIPort               = "CDMON2_API_PORT"
	TotalNumberOfCores    = "CDMON2_TOTAL_NUMBER_OF_CORES"
	TotalSizeOfMemoryMb   = "CDMON2_TOTAL_SIZE_OF_MEMORY"
	TotalSizeOfDiskMb     = "CDMON2_TOTAL_SIZE_OF_DISK"
	MinimalNumberOfCores  = "CDMON2_MINIMAL_NUMBER_OF_CORES"
	MinimalSizeOfMemoryMb = "CDMON2_MINIMAL_SIZE_OF_MEMORY"
	MininalSizeOfDiskMb   = "CDMON2_MININAML_SIZE_OF_DISK"
)

type (
	Config struct {
		APIPort               string
		TotalNumberOfCores    int
		TotalSizeOfMemoryMb   int
		TotalSizeOfDiskMb     int
		MinimalNumberOfCores  int
		MinimalSizeOfMemoryMb int
		MinimalSizeOfDiskMb   int
	}
)

func (c *Config) Load() (err error) {
	c.APIPort = getEnv(APIPort)
	c.TotalNumberOfCores, err = strconv.Atoi(getEnv(TotalNumberOfCores))
	if err == nil {
		c.TotalSizeOfMemoryMb, err = strconv.Atoi(getEnv(TotalSizeOfMemoryMb))
	}
	if err == nil {
		c.TotalSizeOfDiskMb, err = strconv.Atoi(getEnv(TotalSizeOfDiskMb))
	}
	if err == nil {
		c.MinimalNumberOfCores, err = strconv.Atoi(getEnv(MinimalNumberOfCores))
	}
	if err == nil {
		c.MinimalSizeOfMemoryMb, err = strconv.Atoi(getEnv(MinimalSizeOfMemoryMb))
	}
	if err == nil {
		c.MinimalSizeOfDiskMb, err = strconv.Atoi(getEnv(MininalSizeOfDiskMb))
	}
	return
}

func getEnv(env string) (value string) {
	value = os.Getenv(env)
	if len(value) == 0 {
		panic("environment variable " + env + " does not exist")
	}
	return
}
