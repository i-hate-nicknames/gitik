# gitik

## Description

Gitik (small and adorable git, similar to "gittie" in English) is my own
implementation of git, done for educational purposes.

## Implementation

Initial plan is to translate python version of mugit (www.leshenko.net/p/ugit),
but eventually my own implementation may diverge from the "source".
Additionally, the plan is to use bare go with standard library only for the key
functionality implementation, but use 3rd party libraries for the things not related
to the version control itself (e.g. parsing arguments).

Eventually, these parts may too be rewritten to get a minimalistic stdlib-only
implementation
