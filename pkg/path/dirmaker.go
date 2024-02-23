package path

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"
	"github.com/Donders-Institute/dr-data-stager/pkg/dr"
	log "github.com/sirupsen/logrus"
)

// NewDirMaker determines the path type and returns a corresponding
// implementation of the DirMaker interface.
func NewDirMaker(path PathInfo, config config.Configuration) DirMaker {
	switch path.Type {
	case TypeIrods:
		return IrodsCollectionMaker{
			base:   path.Path,
			config: config.Dr,
		}
	default:
		return FileSystemDirMaker{
			base: path.Path,
		}
	}
}

// DirMaker defines interface for implementing directory creation in local filesystem and iRODS.
type DirMaker interface {
	Mkdir(ctx context.Context, path string) error
}

// FileSystemDirMaker implements the DirMaker for local filesystem.
type FileSystemDirMaker struct {
	// Base is the top-level directory.
	base string
}

// Mkdir ensures the directory referred by the path is created.
func (m FileSystemDirMaker) Mkdir(ctx context.Context, path string) error {

	if !strings.HasPrefix(path, m.base) {
		path = filepath.Join(m.base, path)
	}

	log.Debugf("creating directory %s", path)

	return os.MkdirAll(path, 0775)
}

// IrodsCollectionMaker implements the DirMaker for iRODS, using the `imkdir` system call.
type IrodsCollectionMaker struct {
	// Base is the top-level collection.
	base   string
	config dr.Config
}

// Mkdir ensures the iRODS collection referred by the path is created.
func (m IrodsCollectionMaker) Mkdir(ctx context.Context, coll string) error {

	// trim the possible leading `i:` used in the syntax of "irsync" for referring to iRODS namespace.
	coll = strings.TrimPrefix(coll, "i:")

	if !strings.HasPrefix(coll, m.base) {
		coll = filepath.Join(m.base, coll)
	}

	log.Debugf("creating collection %s", coll)

	// try to initialize the filesystem if it is not available
	cfg := m.config

	// credential in current context is not the same as the service-account credential
	// defined in the config.
	if u := ctx.Value(dr.KeyCredential).(dr.Credential).Username(); u != cfg.IrodsUser {
		cfg.IrodsUser = u
		cfg.IrodsPass = ctx.Value(dr.KeyCredential).(dr.Credential).Password()
	}

	fsys, err := dr.NewFileSystem("collMaker", cfg)
	if err != nil {
		return err
	}

	if err := fsys.MakeDir(coll, true); err != nil {
		return fmt.Errorf("cannot create %s: %s", coll, err)
	}

	return nil
}
