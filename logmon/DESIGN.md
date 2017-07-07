# Design considerations and potential improvements.

This tool is not brilliantly designed, however some important precautions have
been made, which by themselves point to some improvement opportunities. The
following list is not exhaustive though.

- Take a smarter approach to simplify and reduce the number of responsibilities for monitor.Reader - currently it is doing too much: reading logs, storing log state info and handling alerts
- Evaluate ways to introduce new log format parsers into monitor.Reader (currently tied to xojoc.pw/logparse)
- The monitor package relies internally in a mutex that works as sort of global lock. Although it simplified a lot the development of the initial version, there is surely a contention problem
- The UI itself could use some considerations, for instance: it relies on a timer to update itself, including alerts. There could be some way to let monitor alerts to trigger an UI update, instead of having to wait a screen refresh cycle
