package internal

import "github.com/indium114/slag"

func QueryWare(lock *Lockfile, name string) error {
	w, ok := lock.Wares[name]
	if !ok {
		return slag.Err("ware %s not found in lockfile", name)
	}

	slag.Query("ware %s\n", name)
	slag.Query(" repo    : %s\n", w.Repo)
	slag.Query(" version : %s\n", w.Version)
	slag.Query(" digest  : %s\n", w.Digest)
	slag.Query(" system  : %t\n", w.System)

	return nil
}

func QueryBlueprint(lock *Lockfile, name string) error {
	w, ok := lock.Blueprints[name]
	if !ok {
		return slag.Err("blueprint %s not found in lockfile", name)
	}

	slag.Query("blueprint %s\n", name)
	slag.Query(" repo   : %s\n", w.Repo)
	slag.Query(" commit : %s\n", w.BuiltCommit)
	slag.Query(" system : %t\n", w.System)

	return nil
}
