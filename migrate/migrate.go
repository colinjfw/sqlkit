package migrate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/coldog/sqlkit/db"
)

const table = `CREATE TABLE migrations (id INT PRIMARY KEY, version INT)`

type Migration struct {
	ID      int
	Version int
}

type Migrator struct {
	DB            db.DB
	MigrationsDir string
}

// Migrate will migrate to the specified version, 0 represents the initial state
// of the database, therefore migrating to 0 will wipe the db. Migrating down
// will migrate down to a state where the specified version was the last to be
// applied.
func (m *Migrator) Migrate(desired int) error {
	history, err := m.List()
	if err != nil {
		return err
	}

	versions, err := m.Versions()
	if err != nil {
		return err
	}

	dir, path, err := solve(history[0].Version, desired, versions)
	if err != nil {
		return err
	}

	for _, version := range path {

	}
}

func (m *Migrator) List() ([]Migration, error) {
	ctx := context.Background()
	var out []Migration
	err := m.DB.Query(ctx, db.
		Select("id", "version").
		From("migrations").
		OrderBy("id", "DESC")).
		Decode(&out)
	return out, err
}

func (m *Migrator) Versions() (versions []int, err error) {
	err = filepath.Walk(
		m.MigrationsDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.Contains(path, ".up.") {
				return nil
			}
			v := strings.Split(filepath.Base(path), ".")[0]
			vi, err := strconv.ParseInt(v, 10, 0)
			if err != nil {
				return fmt.Errorf("failed to parse version from (%s): %v", path, err)
			}
			versions = append(versions, int(vi))
			return nil
		},
	)
	return
}

func (m *Migrator) SQL(dir string) (versions []int, err error) {
	err = filepath.Walk(
		m.MigrationsDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.Contains(path, ".up.") {
				return nil
			}
			v := strings.Split(filepath.Base(path), ".")[0]
			vi, err := strconv.ParseInt(v, 10, 0)
			if err != nil {
				return fmt.Errorf("failed to parse version from (%s): %v", path, err)
			}
			versions = append(versions, int(vi))
			return nil
		},
	)
	return
}

const (
	none int = iota
	up
	down
)

// Solve will return the direction and a list of migration versions to apply.
func solve(currentVersion, desiredVersion int, versions []int) (int, []int, error) {
	var currentIdx int
	var desiredIdx int
	for idx, v := range versions {
		if v == currentVersion {
			currentIdx = idx + 1
		}
		if v == desiredVersion {
			desiredIdx = idx + 1
		}
		if v == 0 {
			return none, nil, fmt.Errorf("migrate: invalid version at index (%d) must be non-zero", idx)
		}
	}
	if currentIdx == 0 && currentVersion != 0 {
		return none, nil, fmt.Errorf("migrate: could not find current version: %v", currentVersion)
	}
	if desiredIdx == 0 && desiredVersion != 0 {
		return none, nil, fmt.Errorf("migrate: could not find desired version: %v", currentVersion)
	}
	if desiredIdx == currentIdx {
		return none, nil, nil
	}
	// Append the 'initial' version onto the list.
	versions = append([]int{0}, versions...)
	if desiredIdx > currentIdx {
		s := versions[currentIdx+1 : desiredIdx+1]
		return up, s, nil
	}
	if desiredIdx < currentIdx {
		s := versions[desiredIdx+1 : currentIdx+1]
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
		return down, s, nil
	}
	return none, nil, nil
}
