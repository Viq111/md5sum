# md5sum
## A md5sum clone in Go

This clone basically have the same functionalities as the original [md5sum](http://linux.die.net/man/1/md5sum) tool but also implement the `--recursive` flag allowing user to compute md5 hashes of entire directories

Command-line usage:
```
usage: md5sum [--help] [--check] [--recursive] [--quiet] [--status] file

Print or check MD5 (128-bit) checksums.

positional arguments:
  file                  file to compute md5 from

optional arguments:
  --help                show this help message and exit
  --check               read MD5 sums from the FILEs and check them
                        (default: False)
  --recursive           Apply md5sum on all files in the directory
                        (default: False)
  --quiet               don't print OK for each successfully verified file
                        (default: False)
  --status              don't output anything, status code shows successfully
                        (default: False)
```
