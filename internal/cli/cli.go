package cli

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"tlsversion/internal/client"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

var (
	Version = "1.0_dev"
)

const (
	defaultTimeout = 2
)

// getExe return the executable name without any path
func getExe() string {
	path, _ := os.Executable()
	exe := filepath.Base(path)
	return exe
}

const textUsage = `
%[1]s display TLS versions (TLSv1.0, TLSv1.1, TLSv1.2, TLSv1.3) supported by a server

%[2]s:
  %[1]s [FLAGS] <host_names> | <hosts_file>

%[3]s:
  -h, --help       displays %[1]s help
  -v, --version    displays %[1]s version
  -f, --file       reads input from a file (one host per line)
  -t, --timeout    connection timeout in seconds (default: 2)
  -d, --debug      outputs debug information

%[4]s:
  %[1]s --help
  %[1]s --version
  %[1]s --file some_hosts.txt
  %[1]s google.com
  %[1]s example.com:8443
  %[1]s google.com example.com:8443
  %[1]s --timeout 5 google.com example.com:8443

`

// Usage displays the usage of the command with all its sub commands
func Usage() {
	bold := color.New(color.Bold).SprintFunc()
	yellow := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	fmt.Printf(textUsage, bold(getExe()), yellow("USAGE"), yellow("FLAGS"), yellow("EXAMPLES"))
}

// Command holds the options and argument of the CLI
type Command struct {
	help    bool
	version bool
	debug   bool
	file    string
	timeout int
	hosts   []string
}

// ParseOptions parse CLI options and return a populated Command
func ParseOptions() (*Command, error) {

	// Overwrite the default help to show the overall tool usage rather than the usage for the top flags
	// To test it, execute the app with a non-valid option
	flag.Usage = func() {
		Usage()
	}

	var cmd = Command{
		help:    false,
		version: false,
		timeout: defaultTimeout,
		file:    "",
	}

	flag.BoolVar(&cmd.help, "help", false, "help")
	flag.BoolVar(&cmd.help, "h", false, "help")
	flag.BoolVar(&cmd.version, "version", false, "version")
	flag.BoolVar(&cmd.version, "v", false, "version")
	flag.BoolVar(&cmd.debug, "debug", false, "debug")
	flag.BoolVar(&cmd.debug, "d", false, "debug")
	flag.StringVar(&cmd.file, "file", "", "file")
	flag.StringVar(&cmd.file, "f", "", "file")
	flag.IntVar(&cmd.timeout, "timeout", defaultTimeout, "timeout")
	flag.IntVar(&cmd.timeout, "t", defaultTimeout, "timeout")

	flag.Parse()

	var err error
	if cmd.file != "" {
		cmd.hosts, err = readFile(cmd.file)
		if err != nil {
			return nil, err
		}
	}

	cmd.hosts = flag.Args()

	return &cmd, nil
}

// resultToString converts a result from the TLS validation to a string
func resultToString(r client.Result) string {
	switch r {
	case client.NotSupported:
		return "N"
	case client.Supported:
		return "Y"
	default:
		return "-"
	}
}

// TlsSupport container keeping track of each supported TLS version for a host
type TlsSupport struct {
	Host         string
	Tls10Support string
	Tls11Support string
	Tls12Support string
	Tls13Support string
	Error        string
}

// verifyHost verifies the TLS versions (1.0, 1.1, 1.2, 1.3) for a host
func verifyHost(host string, timeout int, ch chan<- *TlsSupport) {
	tlsVersions := []uint16{tls.VersionTLS10, tls.VersionTLS11, tls.VersionTLS12, tls.VersionTLS13}

	var support = &TlsSupport{
		Host:         host,
		Tls10Support: "-",
		Tls11Support: "-",
		Tls12Support: "-",
		Tls13Support: "-",
		Error:        "-",
	}

	for _, v := range tlsVersions {
		s, err := client.SupportedTls(host, v, timeout)
		if s == client.ConnectionFailed {
			support.Error = err.Error()
		}
		switch v {
		case tls.VersionTLS10:
			support.Tls10Support = resultToString(s)
		case tls.VersionTLS11:
			support.Tls11Support = resultToString(s)
		case tls.VersionTLS12:
			support.Tls12Support = resultToString(s)
		case tls.VersionTLS13:
			support.Tls13Support = resultToString(s)
		}
	}
	ch <- support
}

// processHosts triggers the process of verifying the TLS version for the hosts
func (cmd Command) processHosts() {

	ch := make(chan *TlsSupport)
	defer close(ch)

	for _, host := range cmd.hosts {
		go verifyHost(host, cmd.timeout, ch)
	}

	var versions []*TlsSupport

	for range cmd.hosts {
		ts := <-ch
		versions = append(versions, ts)
	}

	headerFmt := color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	columnFmt := color.New(color.FgHiYellow, color.Bold).SprintfFunc()

	tbl := table.New("Host", "TLS1.0", "TLS1.1", "TLS1.2", "TLS1.3", "Error")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	// Sort the collection by host name
	sort.Slice(versions, func(i, j int) bool { return versions[i].Host < versions[j].Host })

	for _, v := range versions {
		tbl.AddRow(v.Host, v.Tls10Support, v.Tls11Support, v.Tls12Support, v.Tls13Support, v.Error)
	}

	tbl.Print()
}

// Execute the command from the properties of Command
func (cmd Command) Execute() error {

	if cmd.help {
		Usage()
		return nil
	}

	if cmd.version {
		green := color.New(color.FgHiGreen, color.Bold)
		_, err := green.Printf("%s %s\n", getExe(), Version)
		if err != nil {
			return err
		}
		return nil
	}

	if len(cmd.hosts) == 0 && cmd.file == "" {
		Usage()
		return fmt.Errorf("not enough options or arguments")
	}

	var err error
	// A file was passed as argument
	if len(cmd.file) > 0 {
		cmd.hosts, err = readFile(cmd.file)
		if err != nil {
			return err
		}
		cmd.processHosts()
		return nil
	}

	// Process arguments (each argument is expected to be a host)
	cmd.processHosts()

	return nil
}

// readFile reads a file with one host per line and return an array of hosts
// skips empty lines and return an array of lines stripped from the EOL
func readFile(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	var hosts []string

	fs := bufio.NewScanner(f)
	fs.Split(bufio.ScanLines)

	for fs.Scan() {
		line := fs.Text()
		host := strings.TrimSpace(line)
		if len(host) == 0 {
			continue
		}
		if strings.HasPrefix(host, "#") {
			continue
		}
		hosts = append(hosts, line)
	}

	return hosts, nil
}
