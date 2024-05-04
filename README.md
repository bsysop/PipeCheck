# PipeCheck

PipeCheck is a command-line tool written in Go that reads from stdin, validates URLs or domains, and conditionally forwards data to stdout based on user confirmation. It is particularly useful in processing pipelines where validation and user intervention are required before passing data to subsequent commands.

## Features

- **Data Validation**: Supports validation of URLs and domains to ensure data integrity before processing.
- **Interactive Confirmation**: Allows user interaction to confirm whether to proceed with sending data forward after validation.
- **Efficiency**: Designed to handle input efficiently, making it suitable for large datasets.

## Requirements

- Go 1.15 or higher
- A Unix-like operating system; Windows users might need to adjust ANSI color codes or use compatible terminals.

## Installation

You can install PipeCheck directly using the `go install` command:

```bash
go install github.com/bsysop/PipeCheck@latest
```

This command retrieves the latest version of PipeCheck from GitHub, compiles it, and installs the executable in your Go bin directory, making it accessible from anywhere on your system.

## Usage

PipeCheck reads input from stdin, thus can be used in a pipeline as follows:

```bash
cat data.txt | ./PipeCheck [options]
```

### Options

- `-urls`: Validate each line of the input as a URL.
- `-domains`: Validate each line of the input as a domain.

### Examples

Validate URLs from a file and ask for user confirmation to proceed:

```bash
cat urls.txt | ./PipeCheck -urls
```

Validate domains from a file and ask for user confirmation to proceed:

```bash
cat domains.txt | ./PipeCheck -domains
```

## Output

- Prints the size of the input data and number of lines.
- Shows the first and last two lines of input.
- Uses colors to highlight different parts of the output:
  - **Green**: General information and prompts.
  - **Yellow**: Warnings or invalid entries.
  - **Red**: Errors or critical messages.

If the input fails validation (more than 10 invalid entries), the program will output an error and terminate execution.

## Contributing

Contributions to improve PipeCheck are welcome. Please feel free to fork the repository, make changes, and submit pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
