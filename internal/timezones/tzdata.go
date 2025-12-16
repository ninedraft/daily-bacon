package timezones

import (
	"errors"
	"io/fs"
	"iter"
	"log"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
	"testing/fstest"

	_ "embed"

	"github.com/ninedraft/itermore"
)

func All() []string {
	return allTimezones()
}

var allTimezones = sync.OnceValue(func() (timezones []string) {
	ls := func(fsys fs.FS) iter.Seq[string] {
		return func(yield func(string) bool) {
			if walkErr := fs.WalkDir(fsys, ".", func(fpath string, d fs.DirEntry, err error) error {
				if err != nil || !d.IsDir() || path.Ext(d.Name()) != "" {
					return err
				}

				if !yield(fpath) {
					return fs.SkipAll
				}

				return nil
			}); walkErr != nil && !errors.Is(walkErr, fs.SkipAll) {
				log.Printf("[WARNING] walking timezones: %v", walkErr)
			}
		}
	}

	zoneinfo := func(dir string) (fs.FS, func() error) {
		zoneDir, err := os.OpenRoot(dir)
		if err != nil {
			log.Printf("[WARNING] loading timezone info: %v", err)
			return fstest.MapFS{}, func() error { return nil }
		}

		return zoneDir.FS(), zoneDir.Close
	}

	shareZoneInfo, closeShareZoneInfo := zoneinfo("/usr/share/zoneinfo")
	defer func() {
		if err := closeShareZoneInfo(); err != nil {
			log.Printf("[WARNING] closing /usr/share/zoneinfo: %v", err)
		}
	}()

	shareLibZoneInfo, closeShareLibZoneInfo := zoneinfo("/usr/share/lib/zoneinfo")
	defer func() {
		if err := closeShareLibZoneInfo(); err != nil {
			log.Printf("[WARNING] closing /usr/share/lib/zoneinfo: %v", err)
		}
	}()

	locations := itermore.Chain(
		tzdataEmbeded(),
		ls(shareZoneInfo),
		ls(shareLibZoneInfo),
	)

	timezones = slices.Compact(slices.Sorted(locations))

	return timezones
})

//go:generate go run ./gen/tzdatals timezones.txt
//go:embed timezones.txt
var timezonesList string

func tzdataEmbeded() iter.Seq[string] {
	return strings.FieldsSeq(timezonesList)
}
