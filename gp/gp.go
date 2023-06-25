package gp

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type DefaultUser struct {
	FirstName  string `json:"FirstName"`
	MiddleName string `json:"MiddleName"`
	LastName   string `json:"LastName"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

// ProjectSettings type allow reading config file
type ProjectSettings struct {
	ProjectName     string `yaml:"ProjectName"`
	Environment     string `yaml:"Environment"` // production, development, test
	HTTPPort        string `yaml:"HTTPPort"`
	DefaultLanguage string `yaml:"DefaultLanguage"`
	AppLink         string `yaml:"AppLink"`

	// MongoDB Settings
	MongoDBURI                string        `yaml:"MongoDBURI"`
	DBName                    string        `yaml:"DBName"`
	DBConnectionRefreshTicker time.Duration `yaml:"DBConnectionRefreshTicker"`

	// JWT Settings
	JWTSecret        string        `yaml:"JWTSecret"`
	JWTTokenDuration time.Duration `yaml:"JWTTokenDuration"`

	// Default User Settings
	DefaultUser DefaultUser `yaml:"DefaultUser"`

	DebugMode bool `yaml:"DebugMode"`
}

// PConfig holds info of project settings
var PConfig ProjectSettings

// Init inits project settings
func Init(fileName string) error {

	// ─── PROJECT SETTINGS ───────────────────────────────────────────────────────────
	// We read config file data
	reader, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("can't open config file. err: %s", err.Error())
	}

	err = yaml.NewDecoder(reader).Decode(&PConfig)
	if err != nil {
		return fmt.Errorf("can't load project settings. err: %s", err.Error())
	}

	if strings.TrimSpace(PConfig.HTTPPort) == "" {
		return errors.New("please provide a HTTP Port")
	}

	if strings.TrimSpace(PConfig.MongoDBURI) == "" {
		return errors.New("please provide a MongoDB connection URI")
	}

	if strings.TrimSpace(PConfig.DBName) == "" {
		return errors.New("please provide a database name")
	}

	if PConfig.DBConnectionRefreshTicker == 0 {
		PConfig.DBConnectionRefreshTicker = 5 * time.Minute // default value
	}

	if strings.TrimSpace(PConfig.JWTSecret) == "" {
		return errors.New("please provide an JWTSecret")
	}

	if PConfig.JWTTokenDuration == 0 {
		PConfig.JWTTokenDuration = 1 * time.Hour // default value
	}

	return nil
}

// RandomNum generates random numbers
func RandomNum(nbOfNumbers int) string {

	sResult := ""

	// We init the random seed
	rSource, _ := uuid.Must(uuid.NewRandom()).Time().UnixTime()
	r := rand.New(rand.NewSource(rSource))
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// build the string
	for i := 1; i <= nbOfNumbers; i++ {
		sResult += fmt.Sprintf("%d", r.Intn(10))
	}
	return sResult
}

// RandomNumInt generates random numbers
func RandomNumInt(nbOfNumbers int) int {
	num, _ := strconv.Atoi(RandomNum(nbOfNumbers))
	return num
}
