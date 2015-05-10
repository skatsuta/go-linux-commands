package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/k0kubun/pp"
)

const (
	httpVer     = "HTTP/1.1"
	headerDelim = ": "
	maxBodyLen  = 1024 * 1024
	serverName  = "LittleHTTP"
	serverVer   = "1.0"
	usage       = "Usage: %s [-port=n] [-chroot -user=u -group=g] [-debug] <docroot>\n"
)

// HTTPHeaderField represents an HTTP header field.
type HTTPHeaderField map[string]string

// HTTPRequest represents an HTTP request.
type HTTPRequest struct {
	method string
	path   string
	ver    float64
	header HTTPHeaderField
	body   []byte
	length int64
}

type option struct {
	debug  bool
	chroot bool
	user   string
	group  string
	port   int
}

type logger interface {
	debug(format string, arg ...interface{})
	err(format string, arg ...interface{})
}

type httpLogger struct {
	debugMode bool
}

func newLogger(debugMode bool) *httpLogger {
	return &httpLogger{debugMode: debugMode}
}

func (l *httpLogger) debug(format string, arg ...interface{}) {
	if !l.debugMode {
		return
	}

	_, _ = pp.Fprintf(os.Stderr, "[debug] "+format, arg...)
}

func (l *httpLogger) err(format string, arg ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "[error]"+format, arg...)
}

var log logger

func main() {
	opt := &option{}
	flag.BoolVar(&opt.debug, "debug", false, "start on debug mode")
	flag.BoolVar(&opt.chroot, "chroot", false, "change root directory")
	flag.StringVar(&opt.user, "user", "user", "user")
	flag.StringVar(&opt.group, "group", "group", "group")
	flag.IntVar(&opt.port, "port", 8080, "port")
	flag.Parse()

	args := flag.Args()

	log = newLogger(opt.debug)
	log.debug("%v\n", opt)

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		return
	}

	docroot := args[0]

	/*
		if opt.chroot {
			if e := setupEnvironment(docroot); e != nil {
				log.err("setupEnvironment(%v, %v, %v) falied: %v",
					docroot, opt.user, opt.group, e)
			}
			docroot = ""
		}
	*/

	// TODO: installSignalHandlers()

	server, err := listenSocket(opt.port)
	if err != nil {
		log.err("listenSocket(%v) failed: %v", opt.port, err)
		return
	}

	log.debug("server fd: %v\n", server)

	if !opt.debug {
		// openLog()

		/*
			if e := becomeDaemon(); e != nil {
				log.err("becomeDaemon() failed: %v", e)
				return
			}
		*/
	}

	serverMain(server, docroot)
}

func setupEnvironment(docroot string) error {
	if e := syscall.Chroot(docroot); e != nil {
		return fmt.Errorf("chroot(%v) failed: %v", docroot, e)
	}

	return nil
}

const maxBacklog = 5

func listenSocket(port int) (int, error) {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return 0, fmt.Errorf("socket() failed: %v", err)
	}

	log.debug("sock: %v\n", sock)

	sa := &syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte{127, 0, 0, 1},
	}

	if e := syscall.Bind(sock, sa); e != nil {
		_ = syscall.Close(sock)
		return 0, fmt.Errorf("bind to %v:%v failed: %v", sa.Addr, sa.Port, err)
	}

	if e := syscall.Listen(sock, maxBacklog); e != nil {
		_ = syscall.Close(sock)
		return 0, fmt.Errorf("listen(%v, %v) failed: %v", sock, maxBacklog, err)
	}

	return sock, nil
}

func becomeDaemon() error {
	rootDir := "/"
	if e := os.Chdir(rootDir); e != nil {
		return fmt.Errorf("os.Chdir(%v) failed: %v", rootDir, e)
	}

	// replace std* with /dev/null so that avoid error if std* are used
	devNull := "/dev/null"
	os.Stdin = os.NewFile(uintptr(syscall.Stdin), devNull)
	os.Stdout = os.NewFile(uintptr(syscall.Stdout), devNull)
	os.Stderr = os.NewFile(uintptr(syscall.Stderr), devNull)

	return nil
}

func serverMain(server int, docroot string) {
	for {
		sock, _, err := syscall.Accept(server)
		if err != nil {
			log.err("accept(%v) failed: %v\n", server, err)
			continue
		}

		log.debug("accept %v\n", sock)

		go func(sock int) {
			s := os.NewFile(uintptr(sock), "socket")

			if e := service(s, s, docroot); e != nil {
				log.err("service() failed: %v\n", e)
			}

			if e := syscall.Close(sock); e != nil {
				log.err("close(%v) failed: %v\n", sock, e)
			}
		}(sock)
	}
}

func service(in, out *os.File, docroot string) error {
	req, err := readRequest(in)
	if err != nil {
		return fmt.Errorf("readRequest() failed: %v", err)
	}

	err = respondTo(req, out, docroot)
	if err != nil {
		return fmt.Errorf("respondTo() failed: %v", err)
	}

	return nil
}

func readRequest(in *os.File) (*HTTPRequest, error) {
	r := bufio.NewReader(in)

	// read request line
	line, err := r.ReadString('\n')
	if err == io.EOF {
		return nil, nil
	}

	log.debug("req header: %v\n", line)

	// parse request line to HTTPRequest
	req, err := parseRequestLine(line)
	if err != nil {
		return nil, fmt.Errorf("failed to read request: %v", err)
	}

	// iterate to read request header fields
	h := make(HTTPHeaderField)
	for {
		line, err := r.ReadString('\n')
		// finish reading if EOF or empty line
		if line == "\n" || line == "\r\n" {
			break
		}

		// parse header field
		key, val, err := parseHeaderField(line)
		if err != nil {
			return nil, fmt.Errorf("failed to read request: %v", err)
		}

		h[key] = val
	}
	req.header = h

	// set content length
	// TODO: do this in constructor
	l, err := req.contentLength()
	if err != nil {
		return nil, fmt.Errorf("failed to parse content length: %v", err)
	}
	req.length = l

	if req.length == 0 {
		return req, nil
	}

	if req.length > maxBodyLen {
		return nil, fmt.Errorf("request body too long: %v", req.length)
	}

	// set content body
	// TODO: do this in constructor
	b := make([]byte, req.length)
	if _, e := r.Read(b); e != nil {
		return nil, fmt.Errorf("failed to read request body: %v", b)
	}
	req.body = b

	return req, err
}

func parseRequestLine(line string) (*HTTPRequest, error) {
	// trim trailing \r, \n
	line = strings.Trim(line, "\r\n")

	// split request line
	p := strings.Split(line, " ")
	if len(p) != 3 { // METHOD Request-URI HTTP-Version (RFC2616)
		return nil, fmt.Errorf("parse error on request line: %v", line)
	}

	// check supported HTTP version
	if p[2] != httpVer {
		return nil, fmt.Errorf("not supported HTTP version: %v", p[2])
	}

	ver, err := parseHTTPVer(p[2])
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP version: %v", p[2])
	}

	// p consists of ["METHOD", "/path/to/file", "HTTP/ver"]
	req := &HTTPRequest{
		method: strings.ToUpper(p[0]), // HTTP method is upper case
		path:   p[1],
		ver:    ver,
	}

	return req, nil
}

func parseHTTPVer(ver string) (float64, error) {
	l := len("HTTP/")
	return strconv.ParseFloat(ver[l:], 64)
}

// parseHeaderField parses request header to (key, val, err).
func parseHeaderField(line string) (key, val string, err error) {
	// trim trailing \r, \n
	line = strings.Trim(line, "\r\n")

	// split header field
	p := strings.Split(line, headerDelim)
	if len(p) != 2 { // key:value style
		err = fmt.Errorf("parse error on request header field: %v", line)
		return
	}

	key, val = p[0], p[1]
	return
}

func (req *HTTPRequest) contentLength() (int64, error) {
	key := "Content-Length"

	val, err := req.headerValue(key)
	if err != nil {
		return 0, nil
	}

	l, e := strconv.ParseInt(val, 10, 64)
	if e != nil {
		return 0, fmt.Errorf("parse error on ``%v: %v`: %v", key, val, e)
	}

	return l, nil
}

func (req *HTTPRequest) headerValue(name string) (string, error) {
	val, exists := req.header[name]
	if !exists {
		return "", fmt.Errorf("header `%v` not exist", name)
	}

	return val, nil
}

func respondTo(req *HTTPRequest, out *os.File, docroot string) error {
	var err error

	switch req.method {
	case "GET", "HEAD":
		err = doFileResponse(req, out, docroot)
	case "POST":
		methodNotAllowed(req, out)
	default:
		notImplemented(req, out)
	}

	return err
}

func doFileResponse(req *HTTPRequest, out *os.File, docroot string) error {
	fi, err := getFileInfo(docroot, req.path)
	if err != nil {
		notFound(req, out)
		return fmt.Errorf("getFileInfo(%s, %s) failed: %v", docroot, req.path, err)
	}

	outputCommonHeaderFields(req, out, "200 OK")
	fmt.Fprintf(out, "Content-Length: %d\r\n", fi.Size())
	fmt.Fprintf(out, "Content-Type: %s\r\n\r\n", guessContentType(fi))

	// HEAD method responses only headers
	if req.method == "HEAD" {
		return nil
	}

	// read file contents
	fspath := buildFSPath(docroot, req.path)
	err = printFileContents(fspath, out)
	if err != nil {
		return fmt.Errorf("doFileResponse() failed: %v", err)
	}

	return nil
}

func printFileContents(path string, out *os.File) error {
	// Note: suppose the size of the requested file is small
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}

	if _, e := out.Write(buf); e != nil {
		return fmt.Errorf("failed to write: %v", err)
	}

	return nil
}

func notFound(req *HTTPRequest, out *os.File) {
	outputCommonHeaderFields(req, out, "404 Not Found")
	fmt.Fprintf(out, "Content-Type: text/plain\r\n\r\n")
	fmt.Fprintf(out, "File not found\r\n")
}

func methodNotAllowed(req *HTTPRequest, out *os.File) {
	outputCommonHeaderFields(req, out, "405 Method Not Allowed")
	fmt.Fprintf(out, "Content-Type: text/plain\r\n\r\n")
	fmt.Fprintf(out, "The requested method %s is not allowed\r\n", req.method)
}

func notImplemented(req *HTTPRequest, out *os.File) {
	outputCommonHeaderFields(req, out, "501 Not Implemented")
	fmt.Fprintf(out, "Content-Type: text/plain\r\n\r\n")
	fmt.Fprintf(out, "The requested method %s is not implemented\r\n", req.method)
}

func outputCommonHeaderFields(req *HTTPRequest, out *os.File, status string) {
	fmt.Fprintf(out, "%s %s\r\n", httpVer, status)
	fmt.Fprintf(out, "Date: %s\r\n", time.Now().Format(time.RFC1123))
	fmt.Fprintf(out, "Server: %s/%s\r\n", serverName, serverVer)
	fmt.Fprintf(out, "Connection: close\r\n")
}

func guessContentType(fi os.FileInfo) string {
	ext := filepath.Ext(fi.Name())

	log.debug("file: %v   ext: %v\n", fi.Name(), ext)

	switch ext {
	case ".json":
		return "application/json"
	case ".js":
		return "application/javascript"
	case ".htm", ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".jpeg", "jpg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "text/plain"
	}
}

func getFileInfo(docroot, urlpath string) (os.FileInfo, error) {
	path := buildFSPath(docroot, urlpath)
	return os.Lstat(path)
}

func buildFSPath(docroot, urlpath string) string {
	return fmt.Sprintf("%s/%s", docroot, urlpath)
}

func installSignalHandlers() {
	// TODO
	return
}
