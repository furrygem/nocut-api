package apiserver

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BindAddr string         `yaml:"bind_addr" json:"bind_addr"`
	BindPort uint16         `yaml:"bind_port" json:"bind_port"`
	LogLevel string         `yaml:"log_level" json:"log_level"`
	MongoDB  *MongoDBConfig `yaml:"mongodb"`
}

type MongoDBConfig struct {
	Host         string        `yaml:"host" json:"host"`
	Port         uint16        `yaml:"port" json:"port"`
	Database     string        `yaml:"database" json:"database"`
	AuthDB       string        `yaml:"auth_db" json:"auth_db"`
	Collection   string        `yaml:"collection" json:"collection"`
	UsernameFile string        `yaml:"username_file" json:"username_file"`
	PasswordFile string        `yaml:"password_file" json:"password_file"`
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	LinkTTL      time.Duration `yaml:"link_ttl" json:"link_ttl"`
}

func NewMongoDBConfig() *MongoDBConfig {
	c := &MongoDBConfig{
		Host:       "127.0.0.1",
		Port:       27017,
		Database:   "nocut",
		Collection: "links",
		AuthDB:     "admin",
		LinkTTL:    time.Duration(3 * time.Minute),
	}
	return c
}

// Prepares config. Reads credentials from files if needed
func (mdbc *MongoDBConfig) prepare() error {
	// Read username from username file if it's set
	if f := mdbc.UsernameFile; f != "" {
		username, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		mdbc.Username = string(username)
	}

	// Read password from password file if it's set
	if f := mdbc.PasswordFile; f != "" {
		password, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		mdbc.Password = string(password)
	}
	return nil
}

/*
NewConfig creates New config and populates it with default values
BindAddr: "0.0.0.0"
BindPort: 8080
LogLevel: info
*/
func NewConfig() *Config {
	return &Config{
		BindAddr: "0.0.0.0",
		BindPort: 8080,
		LogLevel: "info",
		MongoDB:  NewMongoDBConfig(),
	}
}

func (c *Config) FromYAML(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		return err
	}
	return c.MongoDB.prepare()
}
