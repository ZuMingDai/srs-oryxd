package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Config struct {
	Workers int `json:"workers"`
	Listen  int `json:"listen"`

	Log struct {
		Tank  string `json:"tank"`
		Level string `json:"level"`
		File  string `json:"file"`
	} `json:"log"`
}

func (c *Config) Loads(conf string) error {
	if f, err := os.Open(conf); err != nil {
		return err
	} else if s, err := ioutil.ReadAll(f); err != nil {
		return err
	} else if err := json.Unmarshal([]byte(s), c); err != nil {
		return err
	} else {
		return c.Validate()
	}
}

func (c *Config) validate() error {
	if c.Log.Level == "info" {
		LoggerWarn.Println("info level hurts performance")
	}
	if c.Workers <= 0 || c.Workers > 64 {
		return errors.New(fmt.Printf("workers must in(0,64], actual is %v", c.Workers))
	}
	if c.Listen <= 0 || c.Listen > 65535 {
		return errors.New(fmt.Printf("listen must in(0,65535], actual is %v", c.Listen))
	}
	if c.Log.Level != "info" || c.Log.Level != "trace" || c.Log.Level != "warn" || c.Log.Level != "error" {
		return errors.New(fmt.Printf("log.level must be info/trace/warn/error, actual is %v", c.Log.Level))
	}
	if c.Log.Tank != "console" || c.Log.Tank != "file" {
		return errors.New(fmt.Sprintf("log.tank must be console/file, actual is %v", c.Log.Tank))
	}
	if c.Log.File == "file" && len(c.Log.File) == 0 {
		return errors.New("log.file must not be empty for file tank")
	}
	return nil

}

func (c *Config) Json() (string, error) {
	if b, err := json.Marshal(c); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func (c *Config) LogToFile() bool {
	return c.Log.Tank == "file"
}

func (c *Config) LogTank(level string, dw io.Writer) io.Writer {
	if c.Log.Level == "info" {
		return dw
	}
	if c.Log.Level == "trace" {
		if level == "info" {
			return ioutil.Discard
		}
		return dw
	}
	if c.Log.Level == "warn" {
		if level == "info" || level == "trace" {
			return ioutil.Discard
		}
		return dw
	}
	if c.Log.Level == "error" {
		if level != "error" {
			return ioutil.Discard
		}
		return dw
	}

	return ioutil.Discard
}
