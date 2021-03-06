package coresize

import (
	"io/ioutil"
	"os"
	"path"
)

type CacheFile struct {
	Cache      *Cache
	Cached     bool
	RemotePath string
	LocalPath  string
}

func (c *Cache) NewCacheFile(remotePath string) *CacheFile {
	return &CacheFile{
		Cache:      c,
		Cached:     false,
		RemotePath: remotePath,
		LocalPath:  path.Join(c.Config.CacheFolder, remotePath),
	}
}

func (cf *CacheFile) Image() (ImageFile, error) {
	// Ensure we have the image on local disk
	if !cf.Cached {
		if err := cf.EnsureCached(); err != nil {
			return ImageFile{}, err
		}
	}

	// Compute image hash
	image := NewImageFile(cf.LocalPath)
	if err := image.ComputeHash(); err != nil {
		return ImageFile{}, err
	}

	// Return hash and rendered image
	return image, nil
}

// EnsureCached checks for file on disk, is it isn't there it fetches it from
// S3 and saves it to local disk
func (cf *CacheFile) EnsureCached() error {
	if _, err := os.Stat(cf.LocalPath); os.IsNotExist(err) {
		// fetch file from S3
		bucket := s3BucketFromConfig(cf.Cache.Config)
		fileBytes, err := bucket.Get(cf.RemotePath)
		if err != nil {
			return err
		}

		// make a odirectory for it
		err = ensureFolder(path.Dir(cf.LocalPath))
		if err != nil {
			return err
		}

		// write to disk
		err = ioutil.WriteFile(cf.LocalPath, fileBytes, 0755)
		return err
	} else if err != nil {
		return err
	}

	// if we stated without error that the file is already present locally
	cf.Cached = true
	return nil
}

// FileType proxy to ImageFile
func (cf *CacheFile) FileType() string {
	return NewImageFile(cf.LocalPath).FileType()
}

type Cache struct {
	Config Config
	Files  map[string]*CacheFile
}

func NewCache(c Config) *Cache {
	return &Cache{
		Config: c,
	}
}

func (c *Cache) Setup() error {
	c.Files = map[string]*CacheFile{}

	// fetch a listing of all existing files in s3
	bucket := s3BucketFromConfig(c.Config)
	response, err := bucket.List("", "", "", 500)
	if err != nil {
		return err
	}

	// append all discovered files
	for _, object := range response.Contents {
		c.Files[object.Key] = c.NewCacheFile(object.Key)
	}

	return nil
}

func (c *Cache) Get(filename string) (*CacheFile, bool) {
	f, ok := c.Files[filename]
	return f, ok
}
