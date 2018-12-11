package poplar

import (
	"os"

	"github.com/evergreen-ci/sink/pail"
	"github.com/pkg/errors"
)

func (a *TestArtifact) hasConversion() bool {
	return a.ConvertBSON2FTDC || a.ConvertCSV2FTDC || a.ConvertGzip
}

func (a *TestArtifact) Convert() error {
	if !a.hasConversion() {
		return nil
	}

	if a.LocalFile == "" {
		return errors.New("cannot specify a conversion on a remote file")
	}

	if _, err := os.Stat(a.LocalFile); os.IsNotExist(err) {
		return errors.New("cannot convert non existant file")
	}

	switch {
	case a.ConvertBSON2FTDC:
		fn, err := a.bsonToFTDC(a.LocalFile)
		if err != nil {
			return errors.Wrap(err, "problem converting file")
		}
		a.LocalFile = fn
		fallthrough
	case a.ConvertCSV2FTDC:
		// DO MAGIC
		fn, err := a.csvToFTDC(a.LocalFile)
		if err != nil {
			return errors.Wrap(err, "problem converting file")
		}
		a.LocalFile = fn
		fallthrough
	case a.ConvertGzip:
		fn, err := a.gzip(a.LocalFile)
		if err != nil {
			return errors.Wrap(err, "problem writing file")
		}
		a.LocalFile = fn
	}

	return errors.New("file conversion not implemented")
}

func (a *TestArtifact) Upload() error {
	if a.LocalFile == "" {
		return errors.New("cannot upload unspecified file")
	}

	if _, err := os.Stat(a.LocalFile); os.IsNotExist(err) {
		return errors.New("cannot upload file that does not exist ")
	}

	if a.Bucket == "" {
		return errors.New("cannot upload file, no bucket specified")
	}

	pail.NewS3Bucket()

	return errors.New("upload operation not supported")
}
