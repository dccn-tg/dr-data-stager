package main

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Donders-Institute/dr-data-stager/internal/worker/config"
	"github.com/Donders-Institute/dr-data-stager/pkg/dr"
	"github.com/cyverse/go-irodsclient/fs"
	ifs "github.com/cyverse/go-irodsclient/irods/fs"

	ppath "github.com/Donders-Institute/dr-data-stager/pkg/path"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

// syncOutput registers the error message of a particular file sync error.
type syncOutput struct {
	File  string
	Error error
}

// scanAndSync walks through the files retrieved from the `bufio.Scanner`,
// sync each file from the `srcColl` collection to the `dstColl` collection.
//
// The sync operation is performed in a concurrent way.  The degree of concurrency is
// defined by number of sync workers, `nworkers`.
//
// Files being successfully synced will be returned as a map with key as the filename
// and value as the checksum of the file.
func scanAndSync(ctx context.Context, config config.Configuration, src, dst ppath.PathInfo, nworkers int) (processed chan syncOutput) {

	processed = make(chan syncOutput)

	// initiate a source scanner and performs the scan.
	scanner := ppath.NewScanner(src)
	dirmaker := ppath.NewDirMaker(dst, config)

	files := scanner.ScanMakeDir(ctx, nworkers*8, &dirmaker)

	// create worker group
	var wg sync.WaitGroup
	wg.Add(nworkers)

	// spin off workers
	for i := 1; i <= nworkers; i++ {
		go syncWorker(ctx, &wg, src, dst, files, processed)
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
	wg *sync.WaitGroup,
	src, dst ppath.PathInfo,
	files chan string,
	processed chan syncOutput) {

	// determin the basedir of the source
	srcbase := src.Path
	if src.Mode.IsRegular() {
		srcbase = filepath.Dir(src.Path)
	}

	for {
		select {
		case fsrc, more := <-files:

			// files channel is closed.
			if !more {
				wg.Done()
				return
			}

			// construct the destination path of this particular source `fsrc`
			var fdst string
			if src.Mode.IsRegular() && !dst.Mode.IsDir() {
				// destination isn't a directory, then it should be used as the destination file path.
				fdst = dst.Path
			} else {
				// `dst` is an existing directory/collection, construct the
				// destination file path of this particular file.
				fdst = path.Join(dst.Path, strings.TrimPrefix(fsrc, srcbase))
			}

			switch {
			case src.Type == ppath.TypeIrods && dst.Type == ppath.TypeFileSystem:

				psrc, _ := ppath.GetPathInfo(ctx, fmt.Sprintf("i:%s", fsrc))
				pdst, _ := ppath.GetPathInfo(ctx, fdst)

				if pdst.SameAs(ctx, psrc) {
					log.Debugf("skip transfer: %s == %s\n", fsrc, fdst)
					processed <- syncOutput{
						File:  fsrc,
						Error: nil,
					}
					continue
				}

				// get file from irods
				log.Debugf("irods get: %s -> %s\n", fsrc, fdst)
				processed <- syncOutput{
					File:  fsrc,
					Error: ctx.Value(dr.KeyFilesystem).(*fs.FileSystem).DownloadFileParallel(fsrc, "", fdst, 0, nil),
				}
			case src.Type == ppath.TypeFileSystem && dst.Type == ppath.TypeIrods:

				pdst, _ := ppath.GetPathInfo(ctx, fmt.Sprintf("i:%s", fdst))
				psrc, _ := ppath.GetPathInfo(ctx, fsrc)

				if pdst.SameAs(ctx, psrc) {
					log.Debugf("skip transfer: %s == %s\n", fsrc, fdst)
					processed <- syncOutput{
						File:  fsrc,
						Error: nil,
					}
					continue
				}

				// put file to irods
				log.Debugf("irods put: %s -> %s\n", fsrc, fdst)

				err := ctx.Value(dr.KeyFilesystem).(*fs.FileSystem).UploadFileParallel(fsrc, fdst, "", 0, false, nil)

				// trigger checksum calculation
				if err == nil {
					if conn, err := ctx.Value(dr.KeyFilesystem).(*fs.FileSystem).GetMetadataConnection(); err == nil {

						defer ctx.Value(dr.KeyFilesystem).(*fs.FileSystem).ReturnMetadataConnection(conn)

						if chksum, err := ifs.GetDataObjectChecksum(conn, fdst, ""); err != nil {
							log.Errorf("cannot request checksum: %s\n", err)
						} else {
							log.Debugf("%s (%s)\n", fdst, chksum.GetChecksumString())

							// TODO: compare checksum to confirm the file is correctly uploaded
						}
					}
				}

				processed <- syncOutput{
					File:  fsrc,
					Error: err,
				}

			default:
				// both source/destination has the same type
				processed <- syncOutput{
					File:  fsrc,
					Error: fmt.Errorf("not supported"),
				}
			}
		case <-ctx.Done():
			log.Debugf("sync worker aborted")
			// task has been mark to done.
			wg.Done()
			return
		}
	}
}