<!-- 
file: go.mod
replace github.com/andypangaribuan/project9 => /Users/apangaribuan/repo/github/project9

Find errors not caught by the compilers.
This command vets the package in the current directory.
  $ go vet
Download all dependencies
  $ go mod download
Remove unused dependencies
  $ go mod tidy
Check code format
  $ gofmt -l .
-->

# Sync Service

Lock an event based on channel and key to perform process synchronization.

Use pre-build image from ghcr.io

```text
docker pull ghcr.io/andypangaribuan/ssync:1.0.0
```

Find out the image  
<https://github.com/users/andypangaribuan/packages/container/package/ssync>

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://gitlab.com/treasuryid/service/connect/-/tags).

## Contributions

Feel free to contribute to this project.

If you find a bug or want a feature, but don't know how to fix/implement it, please fill an [`issue`](https://github.com/andypangaribuan/clog/issues).  
If you fixed a bug or implemented a feature, please send a [`pull request`](https://github.com/andypangaribuan/clog/pulls).

## License

MIT License

Copyright (c) 2022 Andy Pangaribuan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
