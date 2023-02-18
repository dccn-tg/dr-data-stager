package path

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

// GetPathInfo resolves the PathInfo of the given path.
func GetPathInfo(ctx context.Context, path string) (PathInfo, error) {

	var info PathInfo

	re := regexp.MustCompile(`^(i|irods):`)

	if re.MatchString(path) {

		ipath := strings.TrimSuffix(re.ReplaceAllString(path, ""), "/")

		info.Path = ipath
		info.Type = TypeIrods

		// check if the namespace refers to a file object.
		if _, err := exec.CommandContext(ctx, "imeta", "ls", "-d", ipath).Output(); err == nil {
			info.Mode = 0
			return info, nil
		}

		// check if the namespace refers to a collection object.
		if _, err := exec.CommandContext(ctx, "imeta", "ls", "-C", ipath).Output(); err == nil {
			info.Mode = os.ModeDir
			return info, nil
		}

	} else {

		info.Path = path
		info.Type = TypeFileSystem

		if fi, err := os.Stat(path); err == nil {
			info.Mode = fi.Mode()
			return info, nil
		}
	}

	return info, fmt.Errorf("file or directory not found: %s", path)
}

// GetNumberOfFiles get total number of files at or within the `path`.
func GetNumberOfFiles(ctx context.Context, path PathInfo) (int, error) {

	nf := 0

	var cmds []string

	if path.Mode.IsRegular() {
		return 1, nil
	}

	// The path is a collection, use iquery to get number of files.
	if path.Mode.IsDir() && path.Type == TypeIrods {
		cmds = append(cmds,
			fmt.Sprintf("iquest --no-page '%%s' \"SELECT DATA_NAME WHERE COLL_NAME = '%s'\" | wc -l", path.Path),
			fmt.Sprintf("iquest --no-page '%%s/%%s' \"SELECT COLL_NAME,DATA_NAME WHERE COLL_NAME like '%s/%%'\" | wc -l", path.Path),
		)
	}

	// The path is a filesystem directory, use command `find -type f` to get number of files.
	if path.Mode.IsDir() && path.Type == TypeFileSystem {
		cmds = append(cmds, fmt.Sprintf("find %s -type f | wc -l", path.Path))
	}

	// run commands to get total number of files within the path.
	for _, cmd := range cmds {

		// run the command through bash
		out, err := exec.CommandContext(ctx, "bash", "-c", cmd).Output()
		if err != nil {
			return nf, fmt.Errorf("cannot count files: %s", err)
		}

		// remove the line ending "\n" from out
		n, err := strconv.Atoi(strings.TrimSuffix(string(out), "\n"))
		if err != nil {
			return nf, fmt.Errorf("cannot count files: %s", err)
		}

		nf += n
	}

	return nf, nil
}

// ScanAndSync walks through the files retrieved from the `bufio.Scanner`,
// sync each file from the `srcColl` collection to the `dstColl` collection.
//
// The sync operation is performed in a concurrent way.  The degree of concurrency is
// defined by number of sync workers, `nworkers`.
//
// Files being successfully synced will be returned as a map with key as the filename
// and value as the checksum of the file.
func ScanAndSync(ctx context.Context, src, dst PathInfo, nworkers int) (success chan string, failure chan SyncError) {

	success = make(chan string)
	failure = make(chan SyncError)

	// initiate a source scanner and performs the scan.
	scanner := NewScanner(src)
	dirmaker := NewDirMaker(dst)

	files := scanner.ScanMakeDir(ctx, src.Path, nworkers*8, &dirmaker)

	// create worker group
	var wg sync.WaitGroup
	wg.Add(nworkers)

	// spin off workers
	for i := 1; i <= nworkers; i++ {
		go syncWorker(ctx, i, &wg, src, dst, files, success, failure)
	}

	go func() {
		// wait for all workers to finish.
		wg.Wait()
		// close success and failure channels.
		close(success)
		close(failure)
	}()

	return
}

func syncWorker(
	ctx context.Context,
	id int,
	wg *sync.WaitGroup,
	src, dst PathInfo,
	files chan string,
	success chan string,
	failure chan SyncError) {

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
		case fsrc, ok := <-files:

			// files channel is closed.
			if !ok {
				wg.Done()
				return
			}

			// determin destination file path.
			fdst := path.Join(dst.Path, strings.TrimPrefix(fsrc, src.Path))

			// run irsync
			cmdExec := "irsync"
			cmdArgs := []string{"-K", "-v", fmt.Sprintf("%s%s", fsrcPrefix, fsrc), fmt.Sprintf("%s%s", fdstPrefix, fdst)}

			log.Debugf("exec: %s %s", cmdExec, strings.Join(cmdArgs, " "))

			if _, err := exec.CommandContext(ctx, cmdExec, cmdArgs...).Output(); err != nil {
				failure <- SyncError{
					File:  fsrc,
					Error: fmt.Errorf("%s %s fail: %s", cmdExec, strings.Join(cmdArgs, " "), err),
				}
			} else {
				success <- fsrc
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
func ScanAndRepl(ctx context.Context, coll PathInfo, rescSrc, rescDst string, nworkers int) (success chan string, failure chan ReplError) {

	success = make(chan string)
	failure = make(chan ReplError)

	// initiate a source scanner and performs the scan.
	scanner := NewScanner(coll)

	files := scanner.ScanMakeDir(ctx, coll.Path, 4096, nil)

	// create worker group
	var wg sync.WaitGroup
	wg.Add(nworkers)

	// spin off workers
	for i := 1; i <= nworkers; i++ {
		go replWorker(ctx, i, &wg, rescSrc, rescDst, files, success, failure)
	}

	go func() {
		// wait for all workers to finish.
		wg.Wait()
		// close success and failure channels.
		close(success)
		close(failure)
	}()

	return
}

func replWorker(
	ctx context.Context,
	id int,
	wg *sync.WaitGroup,
	rescSrc, rescDst string,
	files chan string,
	success chan string,
	failure chan ReplError) {

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

			if _, err := exec.CommandContext(ctx, cmdExec, cmdArgs...).Output(); err != nil {
				failure <- ReplError{
					File:  f,
					Error: fmt.Errorf("%s %s fail: %s", cmdExec, strings.Join(cmdArgs, " "), err),
				}
			} else {
				success <- f
			}
		case <-ctx.Done():
			log.Debugf("repl worker aborted")
			// task has been mark to done.
			wg.Done()
			return
		}
	}
}
