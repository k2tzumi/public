package git // import "cirello.io/exp/sdci/pkg/git"
import (
	"context"
	"log"
	"os"
	"os/exec"

	"cirello.io/errors"
)

// Checkout clones and reset build directory to a given commit.
func Checkout(ctx context.Context, cloneURL, repoDir, commit string) error {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		os.MkdirAll(repoDir, os.ModePerm&0755)
		cmd := exec.CommandContext(ctx, "git", "clone", cloneURL, ".")
		cmd.Dir = repoDir
		out, err := cmd.CombinedOutput()
		log.Println("cloning:", string(out))
		if err != nil {
			return errors.E(err, "cannot clone repository")
		}
	}
	cmd := exec.CommandContext(ctx, "git", "fetch", "--all")
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	log.Println("fetching objects", string(out))
	if err != nil {
		return errors.E(err, "cannot fetch objects")
	}
	cmd = exec.CommandContext(ctx, "git", "reset", "--hard", commit)
	cmd.Dir = repoDir
	out, err = cmd.CombinedOutput()
	log.Println("reset to", string(out))
	return errors.E(err, "cannot reconfigure repository")
}
