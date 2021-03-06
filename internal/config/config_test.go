package config

import (
	"reflect"
	"testing"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/sirupsen/logrus"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		wantConfig Config
		wantErr    error
	}{
		{
			name:       "no args",
			args:       []string{},
			env:        map[string]string{},
			wantConfig: Config{},
			wantErr:    errNoBinaryName,
		},
		{
			name: "success",
			args: []string{
				"test-cmd",
				"--" + flagNetatmoClientID,
				"id",
				"--" + flagNetatmoClientSecret,
				"secret",
				"--" + flagNetatmoUsername,
				"username",
				"--" + flagNetatmoPassword,
				"password",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr:            defaultConfig.Addr,
				LogLevel:        logLevel(logrus.InfoLevel),
				RefreshInterval: defaultRefreshInterval,
				StaleDuration:   defaultStaleDuration,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
					Username:     "username",
					Password:     "password",
				},
			},
			wantErr: nil,
		},
		{
			name: "all env",
			args: []string{
				"test-cmd",
			},
			env: map[string]string{
				envVarListenAddress:       ":8080",
				envVarLogLevel:            "debug",
				envVarRefreshInterval:     "5m",
				envVarStaleDuration:       "10m",
				envVarNetatmoClientID:     "id",
				envVarNetatmoClientSecret: "secret",
				envVarNetatmoUsername:     "username",
				envVarNetatmoPassword:     "password",
			},
			wantConfig: Config{
				Addr:            ":8080",
				LogLevel:        logLevel(logrus.DebugLevel),
				RefreshInterval: 5 * time.Minute,
				StaleDuration:   10 * time.Minute,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
					Username:     "username",
					Password:     "password",
				},
			},
			wantErr: nil,
		},
		{
			name: "no addr",
			args: []string{
				"test-cmd",
				"--" + flagListenAddress,
				"",
				"--" + flagNetatmoClientID,
				"id",
				"--" + flagNetatmoClientSecret,
				"secret",
				"--" + flagNetatmoUsername,
				"username",
				"--" + flagNetatmoPassword,
				"password",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr: defaultConfig.Addr,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
					Username:     "username",
					Password:     "password",
				},
			},
			wantErr: errNoListenAddress,
		},
		{
			name: "no client id",
			args: []string{
				"test-cmd",
				"--" + flagNetatmoClientSecret,
				"secret",
				"--" + flagNetatmoUsername,
				"username",
				"--" + flagNetatmoPassword,
				"password",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr: defaultConfig.Addr,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
					Username:     "username",
					Password:     "password",
				},
			},
			wantErr: errNoNetatmoClientID,
		},
		{
			name: "no client secret",
			args: []string{
				"test-cmd",
				"--" + flagNetatmoClientID,
				"id",
				"--" + flagNetatmoUsername,
				"username",
				"--" + flagNetatmoPassword,
				"password",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr: defaultConfig.Addr,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
					Username:     "username",
					Password:     "password",
				},
			},
			wantErr: errNoNetatmoClientSecret,
		},
		{
			name: "no username",
			args: []string{
				"test-cmd",
				"--" + flagNetatmoClientID,
				"id",
				"--" + flagNetatmoClientSecret,
				"secret",
				"--" + flagNetatmoPassword,
				"password",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr: defaultConfig.Addr,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
					Username:     "username",
					Password:     "password",
				},
			},
			wantErr: errNoNetatmoUsername,
		},
		{
			name: "no password",
			args: []string{
				"test-cmd",
				"--" + flagNetatmoClientID,
				"id",
				"--" + flagNetatmoClientSecret,
				"secret",
				"--" + flagNetatmoUsername,
				"username",
			},
			env: map[string]string{},
			wantConfig: Config{
				Addr: defaultConfig.Addr,
				Netatmo: netatmo.Config{
					ClientID:     "id",
					ClientSecret: "secret",
					Username:     "username",
					Password:     "password",
				},
			},
			wantErr: errNoNetatmoPassword,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			getenv := func(key string) string {
				return tt.env[key]
			}

			config, err := Parse(tt.args, getenv)

			if err != tt.wantErr {
				t.Errorf("got error %q, want %q", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(config, tt.wantConfig) {
				t.Errorf("got config %v, want %v", config, tt.wantConfig)
			}
		})
	}
}
