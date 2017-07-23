# Original code for phpfmt

This code is no longer functional and is kept here for historical purposes.

This is the original code for phpfmt - the open-source project I sold to a
startup in US. The long story short is that they wanted to create an automated
PHP source code modifier, composed of modules - which was the perfect use-case
for phpfmt's internal machinery. More importantly, they wanted no one else to
use it, so they decided to take it offline.

They failed, and recently I reached out for them asking whether I was allowed to
reopen the original source code to the public. Although I am still not allowed
to work on it, I can make it public again in the same state as just before the
transaction.

Nowadays, you really should be using [PHP-CS-Fixer](https://github.com/FriendsOfPHP/PHP-CS-Fixer). However there are some aspects
that is worth checking in this implementation.

First, it was written as if it was being written in Plan 9's C variant. `src/fmt.src.php`,
for instance, import explicitly all its dependencies at the top of the file in
the exact order they are supposed to be imported. Also, no imported library
would ever have an import internally. `src/build.php` would actually check that.

Also, note that all dependencies were vendored as if they would be vendored in
Go. The rationale was that I wouldn't trust the state of PHP's ecosystem, with
libraries changing or versions vanishing away. This was a design decision that
the reality proved to be wise during the outages of Github.

Also inspired by Go, it used CSP-like channels for concurrency on *nix systems.
Often, it would perform better than PHP-CS-Fixer for that alone.

And lastly, it used pure array lists of tokens, instead of using PHP SPL data
structures. At a certain point, I actually compared the performance of phpfmt
with and without PHP SPL datastructures, and not using them was always faster.

You might note the use of complete `if () ... else ... ` blocks, where cleaner
and more elegant code would not make them necessary. The reason was to avoid
allocation costs - it would consistenly save around 5% of wall-clock time not
having them.

At the moment I sold the rights of phpfmt, I was in the process of proving that
it was modular and composable enough to be decomposed into smaller components,
and that's why `src/autopreincrement.src.php` was created.

phpfmt used semver initially to control its version. Later it swapped to CoreOS
versioning style. The version deployed here (`825.0`) represented the number of
days since phpfmt inception.

The last but not the least, phpfmt is clearly divided into two parts: `src/Core`
and `src/Additionals`. `src/Core` was a set of normalizing transformations -
they aimed to produce an uniform representation of the source code on top of
which `src/Additionals` transformation could safely assume context while
enhancing them.
