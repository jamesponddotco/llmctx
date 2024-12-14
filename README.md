# LLMCTX

`llmctx` is a command-line tool that converts the content of a directory
into context for Large Language Models (LLMs). It can output in plain
text or Claude-specific XML format, making it easy to provide codebase
context to LLMs.

See _llmctx(1)_ after installing for more information.

## Installation

### From source

First install the dependencies:

- Go 1.23 or above.
- make.
- [scdoc](https://git.sr.ht/~sircmpwn/scdoc).

Then compile and install:

```bash
make
sudo make install
```

## Usage

To use `llmctx` from the command line, run:

```bash
llmctx
```

This will convert the current directory's content to plain text format
and output to `stdout`.

### Common use cases

Convert a specific directory to Claude format:

```bash
llmctx -i '/path/to/project' -c
```

Save output to a file:

```bash
llmctx -i `/path/to/project` -o `output.txt`
```

Ignore specific files or directories:

```bash
llmctx -i '/path/to/project' -x '.git' -x 'node_modules' -x 'package-lock.json'
```

### Available options

```bash
NAME:
   llmctx - convert the content of a directory into context for LLMs

USAGE:
   llmctx [global options]

VERSION:
   0.1.0

GLOBAL OPTIONS:
   --input value, -i value    the directory path to convert (defaults to current directory)
   --output value, -o value   the output file path (defaults to stdout)
   --claude, -c               output in Claude's XML format (defaults to false)
   --show-hidden, -a          show hidden files and directories (defaults to false)
   --ignore value, -x value   patterns to ignore (can be repeated)
   --ignore-gitignore, -g     ignore .gitignore rules (defaults to false)
   --help, -h                 show help
   --version, -v              print the version
```

See _llmctx(1)_ for more usage information.

## Contributing

Anyone can help make `llmctx` better. Send patches on the [mailing
list](https://lists.sr.ht/~jamesponddotco/llmctx-devel) and report bugs
on the [issue tracker](https://todo.sr.ht/~jamesponddotco/llmctx).

You must sign-off your work using `git commit --signoff`. Follow the
[Linux kernel developer's certificate of
origin](https://www.kernel.org/doc/html/latest/process/submitting-patches.html#sign-your-work-the-developer-s-certificate-of-origin)
for more details.

All contributions are made under [the GPL-2.0 license](LICENSE.md).

## Resources

The following resources are available:

- [Support and general discussions](https://lists.sr.ht/~jamesponddotco/llmctx-discuss).
- [Patches and development related questions](https://lists.sr.ht/~jamesponddotco/llmctx-devel).
- [Instructions on how to prepare patches](https://git-send-email.io/).
- [Feature requests and bug reports](https://todo.sr.ht/~jamesponddotco/llmctx).

---

Released under the [GPL-2.0 license](LICENSE.md).
