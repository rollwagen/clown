# Clown

A Git(lab) group cloner

![gopher](https://user-images.githubusercontent.com/7364201/192334813-aabf43e9-a3a6-4e4f-adca-03c105f7a37e.png)

## Configuration & running

The program requires two config parameters:

* the hostname of the Gitlab server: `GITLAB_HOST`
* a Gitlab authentication token: `GITLAB_TOKEN`

These can be either provided as environment variables or
defined in the configuration file `~/.clown`. Environment
variables take precedence over the config file.

`~/.clown` example content:

```sh
GITLAB_TOKEN=glpat-sddlj890usdfmlwef
GITLAB_HOST=gitlab.company.com  # no protocol prefix and no trailing slash
```

## Terminology

| Bitbucket    | Github        | Gitlab        |
|:--:          |:--:           |:--:           |
| Pull Request | Pull Request  | Merge Request |
| Snipped      | Gist          | Snippet       |
| Repository   | Repository    | Project       |
| Teams        | Organizations | Groups        |
