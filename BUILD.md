# gitas

## Build

```bash
goreleaser build --clean
```

## Typical release workflow

```bash
git add --update
```

```bash
git commit -m "fix: change"
```

```bash
git tag "$(svu next)"
git push --tags
goreleaser release --clean
```

## Cookiecutter initiation

```bash
cookiecutter \
  ssh://git@github.com/lukasz-lobocki/go-cookiecutter.git \
  package_name="gitas"
```

### was run with following variables

- package_name: **`gitas`**;
package_short_description: `Multiple git repos management.`

- package_version: `1.3.4`

- author_name: `Lukasz Lobocki`;
open_source_license: `CC0 v1.0 Universal`

- __package_slug: `gitas`

### on

`2023-09-12 10:23:39 +0200`
