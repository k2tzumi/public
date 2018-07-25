package git // import "cirello.io/exp/sdci/pkg/git"
import (
	"log"
	"os"
	"os/exec"

	"cirello.io/errors"
)

// Checkout clones and reset build directory to a given commit.
func Checkout(cloneURL, repoDir, commit string) error {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		os.MkdirAll(repoDir, os.ModePerm&0755)
		cmd := exec.Command("git", "clone", cloneURL, ".")
		cmd.Dir = repoDir
		out, err := cmd.CombinedOutput()
		log.Println("cloning:", string(out))
		if err != nil {
			return errors.E(err, "cannot clone repository")
		}
	}
	cmd := exec.Command("git", "fetch", "--all")
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	log.Println("fetching objects", string(out))
	if err != nil {
		return errors.E(err, "cannot fetch objects")
	}
	cmd = exec.Command("git", "reset", "--hard", commit)
	cmd.Dir = repoDir
	out, err = cmd.CombinedOutput()
	log.Println("reset to", string(out))
	if err != nil {
		return errors.E(err, "cannot reconfigure repository")
	}
	return nil
}
