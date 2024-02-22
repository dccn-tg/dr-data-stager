package path

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"
	log "github.com/sirupsen/logrus"
)

// GetPathInfo resolves the PathInfo of the given path.
func GetPathInfo(ctx context.Context, path string) (PathInfo, error) {

	var info PathInfo

	re := regexp.MustCompile(`^(i|irods):`)

	var err error

	if re.MatchString(path) {

		ipath := strings.TrimSuffix(re.ReplaceAllString(path, ""), "/")

		info.Path = ipath
		info.Type = TypeIrods

		// check if the namespace refers to a file object.
		if _, err = exec.CommandContext(ctx, "imeta", "ls", "-d", ipath).Output(); err == nil {
			info.Mode = 0
			return info, nil
		}

		// check if the namespace refers to a collection object.
		if _, err = exec.CommandContext(ctx, "imeta", "ls", "-C", ipath).Output(); err == nil {
			info.Mode = os.ModeDir
			return info, nil
		}

		err = fmt.Errorf("data or collection not found: %s", info.Path)

	} else {

		info.Path = path
		info.Type = TypeFileSystem

		if fi, err := os.Stat(path); err == nil {
			info.Mode = fi.Mode()
			return info, nil
		}

		err = fmt.Errorf("file or directory not found: %s", info.Path)
	}

	return info, err
}

// ScanAndSync walks through the files retrieved from the `bufio.Scanner`,
// sync each file from the `srcColl` collection to the `dstColl` collection.
//
// The sync operation is performed in a concurrent way.  The degree of concurrency is
// defined by number of sync workers, `nworkers`.
//
// Files being successfully synced will be returned as a map with key as the filename
// and value as the checksum of the file.
func ScanAndSync(ctx context.Context, config config.Configuration, src, dst PathInfo, nworkers int) (processed chan SyncError) {

	processed = make(chan SyncError)

	// initiate a source scanner and performs the scan.
	scanner := NewScanner(src)
	dirmaker := NewDirMaker(dst, config)

	files := scanner.ScanMakeDir(ctx, nworkers*8, &dirmaker)

	// create worker group
	var wg sync.WaitGroup
	wg.Add(nworkers)

	// spin off workers
	for i := 1; i <= nworkers; i++ {
		go syncWorker(ctx, i, &wg, src, dst, files, processed)
	}

	go func() {
		// wait for all workers to finish.
		wg.Wait()
		// close processed channels.
		close(processed)
	}()

	return
}

func syncWorker(
	ctx context.Context,
	id int,
	wg *sync.WaitGroup,
	src, dst PathInfo,
	files chan string,
	processed chan SyncError) {

	fsrcPrefix := ""
	fdstPrefix := ""

	if src.Type == TypeIrods {
		fsrcPrefix = "i:"
	}

	if dst.Type == TypeIrods {
		fdstPrefix = "i:"
	}

	for {
		select {
		case fsrc, more := <-files:

			// files channel is closed.
			if !more {
				wg.Done()
				return
			}

			// determin destination file path.
			fdst := path.Join(dst.Path, strings.TrimPrefix(fsrc, src.Path))

			// run irsync
			cmdExec := "irsync"
			cmdArgs := []string{"-K", "-v", fmt.Sprintf("%s%s", fsrcPrefix, fsrc), fmt.Sprintf("%s%s", fdstPrefix, fdst)}

			log.Debugf("exec: %s %s", cmdExec, strings.Join(cmdArgs, " "))

			_, err := exec.CommandContext(ctx, cmdExec, cmdArgs...).Output()
			processed <- SyncError{
				File:  fsrc,
				Error: err,
			}
		case <-ctx.Done():
			log.Debugf("sync worker aborted")
			// task has been mark to done.
			wg.Done()
			return
		}
	}
}

// ScanAndRepl walks through the files retrieved from the `bufio.Scanner`,
// replicate each file from the `rescSrc` iRODS resrouce to the `rescDst` iORDS resource.
//
// The sync operation is performed in a concurrent way.  The degree of concurrency is
// defined by number of sync workers, `nworkers`.
//
// Files being successfully synced will be returned as a map with key as the filename
// and value as the checksum of the file.
func ScanAndRepl(ctx context.Context, coll PathInfo, rescSrc, rescDst string, nworkers int) (processed chan ReplError) {

	processed = make(chan ReplError)

	// initiate a source scanner and performs the scan.
	scanner := NewScanner(coll)

	files := scanner.ScanMakeDir(ctx, 4096, nil)

	// create worker group
	var wg sync.WaitGroup
	wg.Add(nworkers)

	// spin off workers
	for i := 1; i <= nworkers; i++ {
		go replWorker(ctx, i, &wg, rescSrc, rescDst, files, processed)
	}

	go func() {
		// wait for all workers to finish.
		wg.Wait()
		// close success and failure channels.
		close(processed)
	}()

	return
}

func replWorker(
	ctx context.Context,
	id int,
	wg *sync.WaitGroup,
	rescSrc, rescDst string,
	files chan string,
	processed chan ReplError) {

	for {
		select {
		case f, ok := <-files:
			// files channel is closed.
			if !ok {
				wg.Done()
				return
			}

			// run irepl
			//cmdExec := "irepl"
			//cmdArgs := []string{"-v", "-S", rescSrc, "-R", rescDst, f}
			// run irule
			cmdExec := "irule"
			cmdArgs := []string{"rdmReplicateData(*obj, list('" + rescDst + "'))", "*obj=" + f, "null"}

			log.Debugf("exec: %s %s", cmdExec, strings.Join(cmdArgs, " "))

			_, err := exec.CommandContext(ctx, cmdExec, cmdArgs...).Output()

			processed <- ReplError{
				File:  f,
				Error: err,
			}
		case <-ctx.Done():
			log.Debugf("repl worker aborted")
			// task has been mark to done.
			wg.Done()
			return
		}
	}
}
