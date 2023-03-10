#!/usr/bin/env just --justfile

VERSION := "0.3.0"

alias cr := check-release
alias db := dev-build
alias rel := release
alias lrel := local-release
alias ghp := github-push

# Default recipe (this list)
default:
    @just --list

# Update Go dependencies
update:
    go get -u -d ./...
    go mod tidy -v

# Check release configuration
check-release:
    goreleaser check

# Build a release and publish to GitHub
release:
    goreleaser release --rm-dist

# Build a local snapshot release
local-release:
    goreleaser release --rm-dist --snapshot

# Local development build
dev-build:
    goreleaser build --clean --single-target --snapshot

# Push and tag changes to github
github-push:
    git push
    git tag -a {{VERSION}} -m 'Version {{VERSION}}'
    git push origin --tags

# Clean build and release artifacts
clean:
    -rm -rf dist