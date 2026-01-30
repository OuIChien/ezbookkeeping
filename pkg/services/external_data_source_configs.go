package services

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// ExternalDataSourceConfigService represents external data source config service
type ExternalDataSourceConfigService struct {
	ServiceUsingDB
}

// Initialize a external data source config service singleton instance
var (
	ExternalDataSourceConfigs = &ExternalDataSourceConfigService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
	}
)

// GetConfig returns the config for a specific type
func (s *ExternalDataSourceConfigService) GetConfig(c core.Context, configType models.ExternalDataSourceType) (*models.ExternalDataSourceConfig, error) {
	config := &models.ExternalDataSourceConfig{}
	has, err := s.UserDataDB(0).NewSession(c).Where("type=?", configType).Get(config)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return config, nil
}

// SaveConfig saves or updates the config
func (s *ExternalDataSourceConfigService) SaveConfig(c core.Context, config *models.ExternalDataSourceConfig) error {
	now := time.Now().Unix()
	
	existingConfig, err := s.GetConfig(c, config.Type)
	if err != nil {
		return err
	}

	if existingConfig == nil {
		config.CreatedUnixTime = now
		config.UpdatedUnixTime = now
		_, err = s.UserDataDB(0).NewSession(c).Insert(config)
	} else {
		config.ConfigId = existingConfig.ConfigId
		config.UpdatedUnixTime = now
		_, err = s.UserDataDB(0).NewSession(c).ID(existingConfig.ConfigId).Update(config)
	}

	return err
}
