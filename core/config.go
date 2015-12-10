package core

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Config struct {
	Workers int 'json:"workers"'
	Listen int 'json:"listen"'

	Log struct {
		Tank string 'json:"tank"'
		Level string 'json:"level"'
		File string 'json:"file"'
	}'json:"log"'
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

