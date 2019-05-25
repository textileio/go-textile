# go-textile changelog

## 0.2.2 2019-05-24

Updates go-ipfs to v0.4.21-rc3, which contains a fix for the [too many files open bug](https://github.com/ipfs/go-ipfs/issues/6237). This has been causing our cafes to OOM for quite some time.

#### Docs:

- Adds this changelog :)

#### Docker:

- Exposes the profiling server port

#### CLI API:

- Fixes a v0.2.0 regression that did not allow cafe URL with the `cafe add` command
- Moves the init flags under the init command

## 0.2.1 2019-05-23

###  CLI: Fix --version flag (#785)

Fixed v0.2.0 regression where `textile --version` would error with `textile: error: unknown long flag '--version', try --help`.

## 0.2.0 2019-05-23

### Decrypt Blocks + Kingpin CLI (#776)

Brings new CLI commands and HTTP APIs to decrypt file content for you, as well as an improved CLI experience that is better documented and speedier.

Closes #688.

#### Internal:

- Version bump to `0.2.0`
- Renames `FileData` to `FileContent`
- Rewrites CLI from [go-flags](https://github.com/jessevdk/go-flags) to [kingpin](https://github.com/alecthomas/kingpin), advantages:
  - easier for beginners
  - catches errors at compile time, not runtime
  - documentation for CLI args (not just flags and commands)
  - strips quotes for us
  - handles required flags/args for us
  - handles env vars for us
  - handles default vars for us
  - handles types for us
  - way smaller footprint, about 2/3 the size
- CLI handling moved from `textile.go` to `cmd/main.go`
- Moved helpers outside of context, and into `camelCase` as progress for textileio/meta#34
- Simplified http api error sending

#### CLI API:

- Commands are singular with plural alias, for consistency and flexibility, before it was a mix of both, other aliases include
  - `list`, `ls`
  - `remove`, `rm`
  - `ignore`, `remove`, `rm`
  - `delete`, `del`, `remove`, `rm`
  - `unsubscribe`, `unsub`, `remove`, `rm`
  - `search`, `find`
- Descriptions for many commands now updated to be clearer and more consistent
- Arguments for commands are now documented and made explicit
- `textile blocks files <blockid> [--index [index] --path [path] --content]`
- `textile file get <fileid> [--content]`
- `textile profile set [--name [name]] [--avatar [avatar]]` along with these for b/c
  - `textile profile set name <name>`
  - `textile profile set avatar <avatar>`
- `textile docs` now outputs correctly escaped HTML with nice formatting instead of markdown. HTML can be embedded into markdown, and allows tuples of name and multi-line  descriptions via table elements.
- `textile ipfs id` now `textile ipfs peer` with b/c alias
- `textile config var value` will now output a useful error message and a suggestion if `value` was not formatted properly, it also now documents the `textile config` use as a config listing output

#### HTTP API:

- `/blocks/:blockid` redirects to `/blocks/:blockid/meta ` for b/c
- `/blocks/:blockid/meta` returns metadata for the block
- `/blocks/:blockid/files` returns the files for the block
- `/blocks/:blockid/files/:index/:path/meta` returns the metadata for the path at the block
- `/blocks/:blockid/files/:index/:path/content` returns the decrypted content for the path at the block
- `/files/:blockid` redirects to `/blocks/:blockid/files` for b/c
- `/file/:fileid` redirects to `/file/:fileid/meta` for b/c
- `/file/:fileid/data` redirects to `/file/:fileid/content` for b/c
- `/file/:fileid/meta` returns the file metadata
- `/file/:fileid/content` returns the decrypted file content
- Correct error reporting for block file api failures
