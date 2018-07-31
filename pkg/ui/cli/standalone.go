// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"context"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"cirello.io/cci/pkg/coordinator"
	"cirello.io/cci/pkg/infra/repositories"
	"cirello.io/cci/pkg/models"
	"cirello.io/cci/pkg/ui/dashboard"
	"cirello.io/cci/pkg/ui/webhooks"
	"cirello.io/cci/pkg/worker"
	"cirello.io/errors"
	"github.com/jmoiron/sqlx"
)

type standaloneMode struct {
	db *sqlx.DB

	mu  sync.Mutex
	err error
}

func standalone(db *sqlx.DB) error {
	m := &standaloneMode{db: db}
	buildsDir := m.buildsDir()
	configuration := m.loadConfiguration()
	ctx, coord := m.startCoordinator(db, configuration)
	m.startWebhooksServer(ctx, coord)
	m.startDashboard(ctx, db)
	m.startWorkers(ctx, buildsDir, configuration, coord)
	if err := coord.Wait(); err != nil {
		m.setError(err)
	}
	return m.err
}

func (m *standaloneMode) setError(err error) {
	if err == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.err = err
}

func (m *standaloneMode) buildsDir() string {
	currentUser, err := user.Current()
	if err != nil {
		m.setError(errors.E(err, "cannot load current user information"))
		return ""
	}
	return filepath.Join(currentUser.HomeDir, ".cci", "builds-%v")
}

func (m *standaloneMode) loadConfiguration() models.Configuration {
	fd, err := os.Open("cci-config.yaml")
	if err != nil {
		m.setError(errors.E("cannot open configuration file"))
		return nil
	}
	configuration, err := models.LoadConfiguration(fd)
	if err != nil {
		m.setError(err)
		return nil
	}
	return configuration
}

func (m *standaloneMode) startCoordinator(db *sqlx.DB, configuration models.Configuration) (context.Context, *coordinator.Coordinator) {
	ctx, coord := coordinator.New(db, configuration)
	if err := coord.Error(); err != nil {
		m.setError(errors.E(err, "coordinator error on start"))
		return nil, nil
	}
	return ctx, coord
}

func (m *standaloneMode) startWorkers(ctx context.Context, buildsDir string,
	configuration models.Configuration, coord *coordinator.Coordinator) {
	err := worker.Start(ctx, buildsDir, configuration, coord)
	m.setError(errors.E(err, "coordinator error on start"))
}

func (m *standaloneMode) startWebhooksServer(ctx context.Context, coord *coordinator.Coordinator) {
	webhookListener, err := net.Listen("tcp", ":6500")
	if err != nil {
		m.setError(errors.E(err, "cannot start web server"))
		return
	}
	webhookServer := webhooks.New(coord)
	go func() {
		err := webhookServer.ServeContext(ctx, webhookListener)
		m.setError(err)
	}()
}

func (m *standaloneMode) startDashboard(ctx context.Context, db *sqlx.DB) {
	dashboardListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		m.setError(errors.E(err, "cannot start dashboard server"))
		return
	}
	dashboardServer := dashboard.New(repositories.Builds(db))
	go func() {
		err := dashboardServer.ServeContext(ctx, dashboardListener)
		m.setError(errors.E(err, "cannot server dashboard"))
	}()
}
