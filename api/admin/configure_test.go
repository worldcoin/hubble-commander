package admin

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ConfigureTestSuite struct {
	*require.Assertions
	suite.Suite
}

func (s *ConfigureTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ConfigureTestSuite) TestConfigure_BatchCreation() {
	enabled := false
	api := &API{
		cfg:                 &config.APIConfig{AuthenticationKey: authKeyValue},
		enableBatchCreation: func(enable bool) { enabled = enable },
	}

	err := api.Configure(contextWithAuthKey(authKeyValue), dto.ConfigureParams{
		CreateBatches: ref.Bool(true),
	})
	s.NoError(err)
	s.True(enabled)

	err = api.Configure(contextWithAuthKey(authKeyValue), dto.ConfigureParams{
		CreateBatches: ref.Bool(false),
	})
	s.NoError(err)
	s.False(enabled)
}

func (s *ConfigureTestSuite) TestConfigure_AcceptTransactions() {
	enabled := false
	api := &API{
		cfg:                 &config.APIConfig{AuthenticationKey: authKeyValue},
		enableTxsAcceptance: func(enable bool) { enabled = enable },
	}

	err := api.Configure(contextWithAuthKey(authKeyValue), dto.ConfigureParams{
		AcceptTransactions: ref.Bool(true),
	})
	s.NoError(err)
	s.True(enabled)

	err = api.Configure(contextWithAuthKey(authKeyValue), dto.ConfigureParams{
		AcceptTransactions: ref.Bool(false),
	})
	s.NoError(err)
	s.False(enabled)
}

func TestConfigureTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigureTestSuite))
}
