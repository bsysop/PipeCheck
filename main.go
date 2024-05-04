package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// Define Color constants
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
)

var (
	checkURLs    bool
	checkDomains bool
)

func init() {
	flag.BoolVar(&checkURLs, "urls", false, "Check if input contains valid URLs")
	flag.BoolVar(&checkDomains, "domains", false, "Check if input contains valid domains")
}

func main() {
	flag.Parse() // Parse the command-line flags

	// Initialize a buffer to read stdin
	reader := bufio.NewReader(os.Stdin)
	var input strings.Builder
	var lines []string

	// Read stdin completely
	for {
		part, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, ColorRed+"Error reading stdin: %v\n"+ColorReset, err)
			os.Exit(1)
		}
		if part != "" {
			input.WriteString(part)
			lines = append(lines, strings.TrimSpace(part))
		}
		if err == io.EOF {
			break
		}
	}

	// Get the complete input as a string and measure its size
	data := input.String()
	dataSize := len(data) / 1024 // Convert bytes to kilobytes

	// Perform validations if requested
	var valid bool
	if checkURLs {
		valid = validateURLs(lines)
	} else if checkDomains {
		valid = validateDomains(lines)
	}

	// Only proceed if the input is valid or no validation is requested
	if checkURLs || checkDomains {
		if !valid {
			fmt.Fprintf(os.Stderr, ColorRed+"Input validation failed.\n"+ColorReset)
			os.Exit(1)
		}
	}

	// Output statistics to stderr
	fmt.Fprintf(os.Stderr, ColorGreen+"Input:"+ColorReset+" %d kbytes | %d lines\n\n", dataSize, len(lines))
	if len(lines) >= 2 {
		fmt.Fprintln(os.Stderr, ColorGreen+"First two lines:"+ColorReset)
		fmt.Fprint(os.Stderr, lines[0]+"\n")
		fmt.Fprint(os.Stderr, lines[1]+"\n\n")
	}
	if len(lines) >= 4 {
		fmt.Fprintln(os.Stderr, ColorGreen+"Last two lines:"+ColorReset)
		fmt.Fprint(os.Stderr, lines[len(lines)-2]+"\n")
		fmt.Fprint(os.Stderr, lines[len(lines)-1]+"\n\n")
	}

	// Open /dev/tty for user input to ensure we read from the terminal
	tty, err := os.Open("/dev/tty")
	if err != nil {
		fmt.Fprintf(os.Stderr, ColorGreen+"Error opening /dev/tty: %v\n"+ColorReset, err)
		os.Exit(1)
	}
	defer tty.Close()
	ttyReader := bufio.NewReader(tty)

	// Ask the user for confirmation
	fmt.Fprintf(os.Stderr, ColorGreen+"Do you want to process and send the data forward? (y/n): "+ColorReset)
	response, err := ttyReader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, ColorRed+"Error reading response: %v\n"+ColorReset, err)
		os.Exit(1)
	}

	// Clean response and check user input
	response = strings.TrimSpace(response)
	if response == "y" || response == "Y" {
		// If confirmed, write the original data to stdout
		fmt.Print(data)
	} else {
		// If not confirmed, write a message to stderr and exit
		fmt.Fprintln(os.Stderr, ColorRed+"Operation aborted by user."+ColorReset)
		os.Exit(0)
	}
}

func validateURLs(lines []string) bool {
	invalidCount := 0
	for _, line := range lines {
		if _, err := url.ParseRequestURI(line); err != nil {
			fmt.Fprintf(os.Stderr, ColorYellow+"Invalid URL: %s\n"+ColorReset, line)
			invalidCount++
			if invalidCount >= 10 {
				fmt.Fprintf(os.Stderr, ColorRed+"Aborting execution: Found more than 10 invalid urls\n"+ColorReset)
				return false
			}
		}
	}
	return invalidCount == 0
}

func validateDomains(lines []string) bool {
	domainRegex := regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	invalidCount := 0
	for _, line := range lines {
		if !domainRegex.MatchString(line) {
			fmt.Fprintf(os.Stderr, ColorYellow+"Invalid domain: %s\n"+ColorReset, line)
			invalidCount++
			if invalidCount >= 10 {
				fmt.Fprintf(os.Stderr, ColorRed+"Aborting execution: Found more than 10 invalid domains\n"+ColorReset)
				return false
			}
		}
	}
	return invalidCount == 0
}
