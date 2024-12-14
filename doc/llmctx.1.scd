llmctx(1)

# NAME

llmctx - convert the content of a directory into context for LLMs

# SYNOPSIS

*llmctx* [OPTIONS]...

# DESCRIPTION

*llmctx* is a command-line utility that converts the content of a directory into
context for Large Language Models (LLMs). It recursively walks through the
directory structure, reading file contents and formatting them in either plain
text or Claude-specific XML format.

By default, *llmctx* ignores hidden files and directories, and respects
_.gitignore_ rules if present in the directory.

# OPTIONS

*-i* _DIRECTORY_, *--input* _DIRECTORY_
	The directory to process. If not specified, defaults to the current
	directory. The directory is walked recursively, and all non-ignored files
	are included in the output.

*-o* _FILE_, *--output* _FILE_
	Write output to _FILE_ instead of standard output. If the file already
	exists, it will be overwritten.

*-c*, *--claude*
	Format output in Claude's XML format, which is specifically designed to take
	advantage of Claude's extended context window. Defaults to false.

*-x*, *--ignore* _PATTERN_
	Specify a pattern of files or directories to ignore. This option can be
	repeated multiple times to ignore different patterns. See *PATTERNS* for
	syntax information.

*-g*, *--ignore-gitignore*
	Ignore any rules specified in _.gitignore_ files. By default, *llmctx*
	respects _.gitignore_ rules if a _.gitignore_ file is present in the input
	directory.

*-a*, *--show-hidden*
	Include hidden files and directories (those starting with a dot) in the
	output. By default, these are ignored.

*-h*, *--help*
	Display help message and exit.

*-v*, *--version*
	Display version information and exit.

# PATTERNS

The *--ignore* option accepts patterns using the following syntax:

- *?* matches any single character except a path separator.
- *\** matches zero or more characters except a path separator.
- *[abc]* matches any character within the brackets.
- *[a-z]* matches any character in the given range.

*Examples of valid patterns:*

- *.git* - matches the .git directory.
- *\*.go* - matches all Go source files.
- *test/\** - matches everything under the test directory.

# EXAMPLES

*1. Convert current directory to Claude format:*
	$ llmctx --claude

*2. Process a specific directory, ignoring certain patterns:*
	$ llmctx --input '/path/to/project' --ignore '\*.log' --ignore 'node_modules'

*3. Include hidden files and save to output file:*
	$ llmctx --show-hidden --output 'context.xml' --claude

*4. Process directory while ignoring .gitignore rules:*
	$ llmctx --ignore-gitignore --input '/path/to/repo'

# BUGS

Report bugs via email to <~jamesponddotco/llmctx@todo.sr.ht> or via the web
interface at https://todo.sr.ht/~jamesponddotco/llmctx.

# SEE ALSO

*gitignore*(5)

# AUTHORS

Maintained by James Pond <james@cipher.host>.
