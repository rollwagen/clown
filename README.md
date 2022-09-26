# clown
A Git(lab) group cloner

## Configuration & Running

The program requires two config parameters:
* the hostname of the Gitlab server: `GITLAB_HOST`
* a Gitlab authentication token: `GITLAB_TOKEN`

These can be either provided as environment variables or
defined in the configuration file "~/.clown". Environment
variables take precedence over the config file.

`~/.clown` example content:

```sh
GITLAB_TOKEN=glpat-sddlj890usdfmlwef
GITLAB_HOST=gitlab.company.com  # no protocol prefix and no trailing slash
```
