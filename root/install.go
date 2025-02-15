package root

import (
	"fmt"

	"github.com/ButterHost69/PKr-cli/config"
)

// TODO: [X] Setup Username
// TODO: [X] Generate Public and Private Keys
// TODO: [ ] Register gRPC Server as a service
func Install(username, serverIp string) error {
	err := config.CreateUserIfNotExists(username, serverIp)
	if err != nil {
		return fmt.Errorf("error, could not create user.\nError%v", err)
	}
	err = config.CreateServerConfigFiles()
	if err != nil {
		return fmt.Errorf("error, could not server config file for user.\nError%v", err)
	}

	return nil
}
