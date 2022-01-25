package config

import "fmt"

func (c *Config) Validate() error {
	if c.API.AuthenticationKey == nil {
		return fmt.Errorf("authentication key is required")
	}
	return nil
}
