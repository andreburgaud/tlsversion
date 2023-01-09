#!/usr/bin/env just --justfile

VERSION := "0.1.0"

alias cr := check-release
alias db := dev-build
alias rel := release
alias ghp := github-push

# Default recipe (this list)
default:
    @just --list

# Update Go dependencies
update:
    go get -u
    go mod tidy -v

# Check release configuration
check-release:
    goreleaser check

# Build a release and publish to GitHub
release:
    goreleaser release --rm-dist

# Local development build
dev-build:
    goreleaser build --rm-dist --single-target --snapshot

# Push and tag changes to github
github-push:
    git push
    git tag -a {{VERSION}} -m 'Version {{VERSION}}'
    git push origin --tags

# Clean build and release artifacts
clean:
    -rm -rf dist