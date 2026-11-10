package internal

import (
	"github.com/gofrs/flock"
	"github.com/indium114/slag"
)

func Lock() error {
	fileLock := flock.New("/var/lock/wares.lock")
	slag.Wait("acquiring lock\n")
	err := fileLock.Lock()
	if err != nil {
		return err
	}
	return nil
}

func Unlock() {
	fileLock := flock.New("/var/lock/wares.lock")
	fileLock.Unlock()
}
