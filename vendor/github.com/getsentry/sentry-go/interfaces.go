package sentry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/getsentry/sentry-go/attribute"
)

const eventType = "event"
const transactionType = "transaction"
const checkInType = "check_in"

var logEvent = struct {
	Type        string
	ContentType string
}{
	"log",
	"application/vnd.sentry.items.log+json",
}

// Level marks the severity of the event.
type Level string

// Describes the severity of the event.
const (
	LevelDebug   Level = "debug"
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
	LevelFatal   Level = "fatal"
)

// SdkInfo contains all metadata about the SDK.
type SdkInfo struct {
	Name         string       `json:"name,omitempty"`
	Version      string       `json:"version,omitempty"`
	Integrations []string     `json:"integrations,omitempty"`
	Packages     []SdkPackage `json:"packages,omitempty"`
}

// SdkPackage describes a package that was installed.
type SdkPackage struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// TODO: This type could be more useful, as map of interface{} is too generic
// and requires a lot of type assertions in beforeBreadcrumb calls
// plus it could just be map[string]interface{} then.

// BreadcrumbHint contains information that can be associated with a Breadcrumb.
type BreadcrumbHint map[string]interface{}

// Breadcrumb specifies an application event that occurred before a Sentry event.
// An event may contain one or more breadcrumbs.
type Breadcrumb struct {
	Type      string                 `json:"type,omitempty"`
	Category  string                 `json:"category,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Level     Level                  `json:"level,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// TODO: provide constants for known breadcrumb types.
// See https://develop.sentry.dev/sdk/event-payloads/breadcrumbs/#breadcrumb-types.

// MarshalJSON converts the Breadcrumb struct to JSON.
func (b *Breadcrumb) MarshalJSON() ([]byte, error) {
	// We want to omit time.Time zero values, otherwise the server will try to
	// interpret dates too far in the past. However, encoding/json doesn't
	// support the "omitempty" option for struct types. See
	// https://golang.org/issues/11939.
	//
	// We overcome the limitation and achieve what we want by shadowing fields
	// and a few type tricks.

	// breadcrumb aliases Breadcrumb to allow calling json.Marshal without an
	// infinite loop. It preserves all fields while none of the attached
	// methods.
	type breadcrumb Breadcrumb

	if b.Timestamp.IsZero() {
		return json.Marshal(struct {
			// Embed all of the fields of Breadcrumb.
			*breadcrumb
			// Timestamp shadows the original Timestamp field and is meant to
			// remain nil, triggering the omitempty behavior.
			Timestamp json.RawMessage `json:"timestamp,omitempty"`
		}{breadcrumb: (*breadcrumb)(b)})
	}
	return json.Marshal((*breadcrumb)(b))
}

type Logger interface {
	// Write implements the io.Writer interface. Currently, the [sentry.Hub] is
	// context aware, in order to get the correct trace correlation. Using this
	// might result in incorrect span association on logs. If you need to use
	// Write it is recommended to create a NewLogger so that the associated context
	// is passed correctly.
	Write(p []byte) (n int, err error)
	// Trace emits a [LogLevelTrace] log to Sentry.
	// Arguments are handled in the manner of [fmt.Print].
	Trace(ctx context.Context, v ...interface{})
	// Debug emits a [LogLevelDebug] log to Sentry.
	// Arguments are handled in the manner of [fmt.Print].
	Debug(ctx context.Context, v ...interface{})
	// Info emits a [LogLevelInfo] log to Sentry.
	// Arguments are handled in the manner of [fmt.Print].
	Info(ctx context.Context, v ...interface{})
	// Warn emits a [LogLevelWarn] log to Sentry.
	// Arguments are handled in the manner of [fmt.Print].
	Warn(ctx context.Context, v ...interface{})
	// Error emits a [LogLevelError] log to Sentry.
	// Arguments are handled in the manner of [fmt.Print].
	Error(ctx context.Context, v ...interface{})
	// Fatal emits a [LogLevelFatal] log to Sentry followed by a call to [os.Exit](1).
	// Arguments are handled in the manner of [fmt.Print].
	Fatal(ctx context.Context, v ...interface{})
	// Panic emits a [LogLevelFatal] log to Sentry followed by a call to panic().
	// Arguments are handled in the manner of [fmt.Print].
	Panic(ctx context.Context, v ...interface{})

	// Tracef emits a [LogLevelTrace] log to Sentry.
	// Arguments are handled in the manner of [fmt.Printf].
	Tracef(ctx context.Context, format string, v ...interface{})
	// Debugf emits a [LogLevelDebug] log to Sentry.
	// Arguments are handled in the manner of [fmt.Printf].
	Debugf(ctx context.Context, format string, v ...interface{})
	// Infof emits a [LogLevelInfo] log to Sentry.
	// Arguments are handled in the manner of [fmt.Printf].
	Infof(ctx context.Context, format string, v ...interface{})
	// Warnf emits a [LogLevelWarn] log to Sentry.
	// Arguments are handled in the manner of [fmt.Printf].
	Warnf(ctx context.Context, format string, v ...interface{})
	// Errorf emits a [LogLevelError] log to Sentry.
	// Arguments are handled in the manner of [fmt.Printf].
	Errorf(ctx context.Context, format string, v ...interface{})
	// Fatalf emits a [LogLevelFatal] log to Sentry followed by a call to [os.Exit](1).
	// Arguments are handled in the manner of [fmt.Printf].
	Fatalf(ctx context.Context, format string, v ...interface{})
	// Panicf emits a [LogLevelFatal] log to Sentry followed by a call to panic().
	// Arguments are handled in the manner of [fmt.Printf].
	Panicf(ctx context.Context, format string, v ...interface{})
	// SetAttributes allows attaching parameters to the log message using the attribute API.
	SetAttributes(...attribute.Builder)
}

// Attachment allows associating files with your events to aid in investigation.
// An event may contain one or more attachments.
type Attachment struct {
	Filename    string
	ContentType string
	Payload     []byte
}

// User describes the user associated with an Event. If this is used, at least
// an ID or an IP address should be provided.
type User struct {
	ID        string            `json:"id,omitempty"`
	Email     string            `json:"email,omitempty"`
	IPAddress string            `json:"ip_address,omitempty"`
	Username  string            `json:"username,omitempty"`
	Name      string            `json:"name,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
}

func (u User) IsEmpty() bool {
	if u.ID != "" {
		return false
	}

	if u.Email != "" {
		return false
	}

	if u.IPAddress != "" {
		return false
	}

	if u.Username != "" {
		return false
	}

	if u.Name != "" {
		return false
	}

	if len(u.Data) > 0 {
		return false
	}

	return true
}

// Request contains information on a HTTP request related to the event.
type Request struct {
	URL         string            `json:"url,omitempty"`
	Method      string            `json:"method,omitempty"`
	Data        string            `json:"data,omitempty"`
	QueryString string            `json:"query_string,omitempty"`
	Cookies     string            `json:"cookies,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
}

var sensitiveHeaders = map[string]struct{}{
	"_csrf":               {},
	"_csrf_token":         {},
	"_session":            {},
	"_xsrf":               {},
	"Api-Key":             {},
	"Apikey":              {},
	"Auth":                {},
	"Authorization":       {},
	"Cookie":              {},
	"Credentials":         {},
	"Csrf":                {},
	"Csrf-Token":          {},
	"Csrftoken":           {},
	"Ip-Address":          {},
	"Passwd":              {},
	"Password":            {},
	"Private-Key":         {},
	"Privatekey":          {},
	"Proxy-Authorization": {},
	"Remote-Addr":         {},
	"Secret":              {},
	"Session":             {},
	"Sessionid":           {},
	"Token":               {},
	"User-Session":        {},
	"X-Api-Key":           {},
	"X-Csrftoken":         {},
	"X-Forwarded-For":     {},
	"X-Real-Ip":           {},
	"XSRF-TOKEN":          {},
}

// NewRequest returns a new Sentry Request from the given http.Request.
//
// NewRequest avoids operations that depend on network access. In particular, it
// does not read r.Body.
func NewRequest(r *http.Request) *Request {
	protocol := schemeHTTP
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = schemeHTTPS
	}
	url := fmt.Sprintf("%s://%s%s", protocol, r.Host, r.URL.Path)

	var cookies string
	var env map[string]string
	headers := map[string]string{}

	if client := CurrentHub().Client(); client != nil && client.options.SendDefaultPII {
		// We read only the first Cookie header because of the specification:
		// https://tools.ietf.org/html/rfc6265#section-5.4
		// When the user agent generates an HTTP request, the user agent MUST NOT
		// attach more than one Cookie header field.
		cookies = r.Header.Get("Cookie")

		headers = make(map[string]string, len(r.Header))
		for k, v := range r.Header {
			headers[k] = strings.Join(v, ",")
		}

		if addr, port, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			env = map[string]string{"REMOTE_ADDR": addr, "REMOTE_PORT": port}
		}
	} else {
		for k, v := range r.Header {
			if _, ok := sensitiveHeaders[k]; !ok {
				headers[k] = strings.Join(v, ",")
			}
		}
	}

	headers["Host"] = r.Host

	return &Request{
		URL:         url,
		Method:      r.Method,
		QueryString: r.URL.RawQuery,
		Cookies:     cookies,
		Headers:     headers,
		Env:         env,
	}
}

// Mechanism is the mechanism by which an exception was generated and handled.
type Mechanism struct {
	Type             string         `json:"type,omitempty"`
	Description      string         `json:"description,omitempty"`
	HelpLink         string         `json:"help_link,omitempty"`
	Source           string         `json:"source,omitempty"`
	Handled          *bool          `json:"handled,omitempty"`
	ParentID         *int           `json:"parent_id,omitempty"`
	ExceptionID      int            `json:"exception_id"`
	IsExceptionGroup bool           `json:"is_exception_group,omitempty"`
	Data             map[string]any `json:"data,omitempty"`
}

// SetUnhandled indicates that the exception is an unhandled exception, i.e.
// from a panic.
func (m *Mechanism) SetUnhandled() {
	m.Handled = Pointer(false)
}

// Exception specifies an error that occurred.
type Exception struct {
	Type       string      `json:"type,omitempty"`  // used as the main issue title
	Value      string      `json:"value,omitempty"` // used as the main issue subtitle
	Module     string      `json:"module,omitempty"`
	ThreadID   uint64      `json:"thread_id,omitempty"`
	Stacktrace *Stacktrace `json:"stacktrace,omitempty"`
	Mechanism  *Mechanism  `json:"mechanism,omitempty"`
}

// SDKMetaData is a struct to stash data which is needed at some point in the SDK's event processing pipeline
// but which shouldn't get send to Sentry.
type SDKMetaData struct {
	dsc DynamicSamplingContext
}

// Contains information about how the name of the transaction was determined.
type TransactionInfo struct {
	Source TransactionSource `json:"source,omitempty"`
}

// The DebugMeta interface is not used in Golang apps, but may be populated
// when proxying Events from other platforms, like iOS, Android, and the
// Web.  (See: https://develop.sentry.dev/sdk/event-payloads/debugmeta/ ).
type DebugMeta struct {
	SdkInfo *DebugMetaSdkInfo `json:"sdk_info,omitempty"`
	Images  []DebugMetaImage  `json:"images,omitempty"`
}

type DebugMetaSdkInfo struct {
	SdkName           string `json:"sdk_name,omitempty"`
	VersionMajor      int    `json:"version_major,omitempty"`
	VersionMinor      int    `json:"version_minor,omitempty"`
	VersionPatchlevel int    `json:"version_patchlevel,omitempty"`
}

type DebugMetaImage struct {
	Type        string `json:"type,omitempty"`         // all
	ImageAddr   string `json:"image_addr,omitempty"`   // macho,elf,pe
	ImageSize   int    `json:"image_size,omitempty"`   // macho,elf,pe
	DebugID     string `json:"debug_id,omitempty"`     // macho,elf,pe,wasm,sourcemap
	DebugFile   string `json:"debug_file,omitempty"`   // macho,elf,pe,wasm
	CodeID      string `json:"code_id,omitempty"`      // macho,elf,pe,wasm
	CodeFile    string `json:"code_file,omitempty"`    // macho,elf,pe,wasm,sourcemap
	ImageVmaddr string `json:"image_vmaddr,omitempty"` // macho,elf,pe
	Arch        string `json:"arch,omitempty"`         // macho,elf,pe
	UUID        string `json:"uuid,omitempty"`         // proguard
}

// EventID is a hexadecimal string representing a unique uuid4 for an Event.
// An EventID must be 32 characters long, lowercase and not have any dashes.
type EventID string

type Context = map[string]interface{}

// Event is the fundamental data structure that is sent to Sentry.
type Event struct {
	Breadcrumbs []*Breadcrumb          `json:"breadcrumbs,omitempty"`
	Contexts    map[string]Context     `json:"contexts,omitempty"`
	Dist        string                 `json:"dist,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	EventID     EventID                `json:"event_id,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
	Fingerprint []string               `json:"fingerprint,omitempty"`
	Level       Level                  `json:"level,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Platform    string                 `json:"platform,omitempty"`
	Release     string                 `json:"release,omitempty"`
	Sdk         SdkInfo                `json:"sdk,omitempty"`
	ServerName  string                 `json:"server_name,omitempty"`
	Threads     []Thread               `json:"threads,omitempty"`
	Tags        map[string]string      `json:"tags,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Transaction string                 `json:"transaction,omitempty"`
	User        User                   `json:"user,omitempty"`
	Logger      string                 `json:"logger,omitempty"`
	Modules     map[string]string      `json:"modules,omitempty"`
	Request     *Request               `json:"request,omitempty"`
	Exception   []Exception            `json:"exception,omitempty"`
	DebugMeta   *DebugMeta             `json:"debug_meta,omitempty"`
	Attachments []*Attachment          `json:"-"`

	// The fields below are only relevant for transactions.

	Type            string           `json:"type,omitempty"`
	StartTime       time.Time        `json:"start_timestamp"`
	Spans           []*Span          `json:"spans,omitempty"`
	TransactionInfo *TransactionInfo `json:"transaction_info,omitempty"`

	// The fields below are only relevant for crons/check ins

	CheckIn       *CheckIn       `json:"check_in,omitempty"`
	MonitorConfig *MonitorConfig `json:"monitor_config,omitempty"`

	// The fields below are only relevant for logs
	Logs []Log `json:"items,omitempty"`

	// The fields below are not part of the final JSON payload.

	sdkMetaData SDKMetaData
}

// SetException appends the unwrapped errors to the event's exception list.
//
// maxErrorDepth is the maximum depth of the error chain we will look
// into while unwrapping the errors. If maxErrorDepth is -1, we will
// unwrap all errors in the chain.
func (e *Event) SetException(exception error, maxErrorDepth int) {
	if exception == nil {
		return
	}

	err := exception

	for i := 0; err != nil && (i < maxErrorDepth || maxErrorDepth == -1); i++ {
		// Add the current error to the exception slice with its details
		e.Exception = append(e.Exception, Exception{
			Value:      err.Error(),
			Type:       reflect.TypeOf(err).String(),
			Stacktrace: ExtractStacktrace(err),
		})

		// Attempt to unwrap the error using the standard library's Unwrap method.
		// If errors.Unwrap returns nil, it means either there is no error to unwrap,
		// or the error does not implement the Unwrap method.
		unwrappedErr := errors.Unwrap(err)

		if unwrappedErr != nil {
			// The error was successfully unwrapped using the standard library's Unwrap method.
			err = unwrappedErr
			continue
		}

		cause, ok := err.(interface{ Cause() error })
		if !ok {
			// We cannot unwrap the error further.
			break
		}

		// The error implements the Cause method, indicating it may have been wrapped
		// using the github.com/pkg/errors package.
		err = cause.Cause()
	}

	// Add a trace of the current stack to the most recent error in a chain if
	// it doesn't have a stack trace yet.
	// We only add to the most recent error to avoid duplication and because the
	// current stack is most likely unrelated to errors deeper in the chain.
	if e.Exception[0].Stacktrace == nil {
		e.Exception[0].Stacktrace = NewStacktrace()
	}

	if len(e.Exception) <= 1 {
		return
	}

	// event.Exception should be sorted such that the most recent error is last.
	slices.Reverse(e.Exception)

	for i := range e.Exception {
		e.Exception[i].Mechanism = &Mechanism{
			IsExceptionGroup: true,
			ExceptionID:      i,
			Type:             "generic",
		}
		if i == 0 {
			continue
		}
		e.Exception[i].Mechanism.ParentID = Pointer(i - 1)
	}
}

// TODO: Event.Contexts map[string]interface{} => map[string]EventContext,
// to prevent accidentally storing T when we mean *T.
// For example, the TraceContext must be stored as *TraceContext to pick up the
// MarshalJSON method (and avoid copying).
// type EventContext interface{ EventContext() }

// MarshalJSON converts the Event struct to JSON.
func (e *Event) MarshalJSON() ([]byte, error) {
	// We want to omit time.Time zero values, otherwise the server will try to
	// interpret dates too far in the past. However, encoding/json doesn't
	// support the "omitempty" option for struct types. See
	// https://golang.org/issues/11939.
	//
	// We overcome the limitation and achieve what we want by shadowing fields
	// and a few type tricks.
	if e.Type == transactionType {
		return e.transactionMarshalJSON()
	}

	if e.Type == checkInType {
		return e.checkInMarshalJSON()
	}
	return e.defaultMarshalJSON()
}

func (e *Event) defaultMarshalJSON() ([]byte, error) {
	// event aliases Event to allow calling json.Marshal without an infinite
	// loop. It preserves all fields while none of the attached methods.
	type event Event

	// errorEvent is like Event with shadowed fields for customizing JSON
	// marshaling.
	type errorEvent struct {
		*event

		// Timestamp shadows the original Timestamp field. It allows us to
		// include the timestamp when non-zero and omit it otherwise.
		Timestamp json.RawMessage `json:"timestamp,omitempty"`

		// The fields below are not part of error events and only make sense to
		// be sent for transactions. They shadow the respective fields in Event
		// and are meant to remain nil, triggering the omitempty behavior.

		Type            json.RawMessage `json:"type,omitempty"`
		StartTime       json.RawMessage `json:"start_timestamp,omitempty"`
		Spans           json.RawMessage `json:"spans,omitempty"`
		TransactionInfo json.RawMessage `json:"transaction_info,omitempty"`
	}

	x := errorEvent{event: (*event)(e)}
	if !e.Timestamp.IsZero() {
		b, err := e.Timestamp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		x.Timestamp = b
	}
	return json.Marshal(x)
}

func (e *Event) transactionMarshalJSON() ([]byte, error) {
	// event aliases Event to allow calling json.Marshal without an infinite
	// loop. It preserves all fields while none of the attached methods.
	type event Event

	// transactionEvent is like Event with shadowed fields for customizing JSON
	// marshaling.
	type transactionEvent struct {
		*event

		// The fields below shadow the respective fields in Event. They allow us
		// to include timestamps when non-zero and omit them otherwise.

		StartTime json.RawMessage `json:"start_timestamp,omitempty"`
		Timestamp json.RawMessage `json:"timestamp,omitempty"`
	}

	x := transactionEvent{event: (*event)(e)}
	if !e.Timestamp.IsZero() {
		b, err := e.Timestamp.MarshalJSON()
		if err != nil {
			return nil, err
		}
		x.Timestamp = b
	}
	if !e.StartTime.IsZero() {
		b, err := e.StartTime.MarshalJSON()
		if err != nil {
			return nil, err
		}
		x.StartTime = b
	}
	return json.Marshal(x)
}

func (e *Event) checkInMarshalJSON() ([]byte, error) {
	checkIn := serializedCheckIn{
		CheckInID:     string(e.CheckIn.ID),
		MonitorSlug:   e.CheckIn.MonitorSlug,
		Status:        e.CheckIn.Status,
		Duration:      e.CheckIn.Duration.Seconds(),
		Release:       e.Release,
		Environment:   e.Environment,
		MonitorConfig: nil,
	}

	if e.MonitorConfig != nil {
		checkIn.MonitorConfig = &MonitorConfig{
			Schedule:      e.MonitorConfig.Schedule,
			CheckInMargin: e.MonitorConfig.CheckInMargin,
			MaxRuntime:    e.MonitorConfig.MaxRuntime,
			Timezone:      e.MonitorConfig.Timezone,
		}
	}

	return json.Marshal(checkIn)
}

// NewEvent creates a new Event.
func NewEvent() *Event {
	return &Event{
		Contexts: make(map[string]Context),
		Extra:    make(map[string]interface{}),
		Tags:     make(map[string]string),
		Modules:  make(map[string]string),
	}
}

// Thread specifies threads that were running at the time of an event.
type Thread struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Stacktrace *Stacktrace `json:"stacktrace,omitempty"`
	Crashed    bool        `json:"crashed,omitempty"`
	Current    bool        `json:"current,omitempty"`
}

// EventHint contains information that can be associated with an Event.
type EventHint struct {
	Data               interface{}
	EventID            string
	OriginalException  error
	RecoveredException interface{}
	Context            context.Context
	Request            *http.Request
	Response           *http.Response
}

type Log struct {
	Timestamp  time.Time            `json:"timestamp,omitempty"`
	TraceID    TraceID              `json:"trace_id,omitempty"`
	Level      LogLevel             `json:"level"`
	Severity   int                  `json:"severity_number,omitempty"`
	Body       string               `json:"body,omitempty"`
	Attributes map[string]Attribute `json:"attributes,omitempty"`
}

type Attribute struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
}
