package internal

import "fmt"

func QueryWare(lock *Lockfile, name string) error {
	w, ok := lock.Wares[name]
	if !ok {
		return fmt.Errorf("%s ware %s not found in lockfile", ErrText, name)
	}

	fmt.Printf("%s ware %s\n", QueryText, name)
	fmt.Printf("%s  repo    : %s\n", QueryText, w.Repo)
	fmt.Printf("%s  version : %s\n", QueryText, w.Version)
	fmt.Printf("%s  digest  : %s\n", QueryText, w.Digest)
	fmt.Printf("%s  system  : %t\n", QueryText, w.System)

	return nil
}

func QueryBlueprint(lock *Lockfile, name string) error {
	w, ok := lock.Blueprints[name]
	if !ok {
		return fmt.Errorf("%s blueprint %s not found in lockfile", ErrText, name)
	}

	fmt.Printf("%s blueprint %s\n", QueryText, name)
	fmt.Printf("%s  repo   : %s\n", QueryText, w.Repo)
	fmt.Printf("%s  commit : %s\n", QueryText, w.BuiltCommit)
	fmt.Printf("%s  system : %t\n", QueryText, w.System)

	return nil
}
