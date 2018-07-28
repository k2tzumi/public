# SDCI experiment

## Design goals

SDCI experiment implement a continuous integration service in a single binary
capable of running either standalone or distributed.

* No retries - on failure, you should have the option to try again.
* No worker filter - all workers attached to a repository should be able to run its CD/CI steps.
* No global queue - every worker is attached to one queue only, and every build target is its own queue.