package sentry

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

type Option func(hook *logrus_sentry.SentryHook) error

type Condition func(res *http.Response, url string) bool

func SetupSentry(dsn string, opts ...Option) error {
	hook, err := logrus_sentry.NewSentryHook(dsn, []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	})
	if err != nil {
		return fmt.Errorf("failed to create new sentry hook: %w", err)
	}

	hook.Timeout = 0
	hook.StacktraceConfiguration.Enable = true
	hook.StacktraceConfiguration.IncludeErrorBreadcrumb = true
	hook.StacktraceConfiguration.Context = 10
	hook.StacktraceConfiguration.SendExceptionType = true
	hook.StacktraceConfiguration.SwitchExceptionTypeAndMessage = true

	for _, o := range opts {
		if err = o(hook); err != nil {
			return err
		}
	}

	log.AddHook(hook)

	return nil
}

func WithDefaultLoggerName(name string) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetDefaultLoggerName(name)

		return nil
	}
}

func WithEnvironment(env string) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetEnvironment(env)

		return nil
	}
}

func WithHTTPContext(h *raven.Http) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetHttpContext(h)

		return nil
	}
}

func WithIgnoreErrors(errs ...string) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		if err := hook.SetIgnoreErrors(errs...); err != nil {
			return fmt.Errorf("failed to set ignore errors: %w", err)
		}

		return nil
	}
}

func WithIncludePaths(p []string) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetIncludePaths(p)

		return nil
	}
}

func WithRelease(release string) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetRelease(release)

		return nil
	}
}

func WithSampleRate(rate float32) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		if err := hook.SetSampleRate(rate); err != nil {
			return fmt.Errorf("failed to set sample rate: %w", err)
		}

		return nil
	}
}

func WithTagsContext(t map[string]string) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetTagsContext(t)

		return nil
	}
}

func WithUserContext(u *raven.User) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetUserContext(u)

		return nil
	}
}

func WithServerName(serverName string) Option {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetServerName(serverName)

		return nil
	}
}

//nolint:gochecknoglobals
var SentryErrorHandler = func(res *http.Response, url string) error {
	statusCode := res.StatusCode
	// Improve ways to identify if worth logging the error
	if statusCode != http.StatusOK && statusCode != http.StatusNotFound {
		log.WithFields(log.Fields{
			"tags": raven.Tags{
				{Key: "status_code", Value: strconv.Itoa(res.StatusCode)},
				{Key: "host", Value: res.Request.URL.Host},
				{Key: "path", Value: res.Request.URL.Path},
				{Key: "body", Value: getBody(res)},
			},
			"url":         url,
			"fingerprint": []string{"client_errors"},
		}).Error("Client Errors")
	}

	return nil
}

// GetSentryErrorHandler initializes sentry logger for http response errors
// Responses to be logged are defined via passed conditions
func GetSentryErrorHandler(conditions ...Condition) func(res *http.Response, url string) error {
	return func(res *http.Response, url string) error {
		for _, condition := range conditions {
			if condition(res, url) {
				log.WithFields(log.Fields{
					"tags": raven.Tags{
						{Key: "status_code", Value: strconv.Itoa(res.StatusCode)},
						{Key: "host", Value: res.Request.URL.Host},
						{Key: "path", Value: res.Request.URL.Path},
						{Key: "body", Value: getBody(res)},
					},
					"url":         url,
					"fingerprint": []string{"client_errors"},
				}).Error("Client Errors")

				break
			}
		}

		return nil
	}
}

//nolint:errcheck
func getBody(res *http.Response) string {
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	_ = res.Body.Close() //  must close
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes)
}

//nolint:gochecknoglobals
var (
	// ConditionAnd returns true only when all conditions are satisfied
	ConditionAnd = func(conditions ...Condition) Condition {
		return func(res *http.Response, url string) bool {
			result := true
			for _, condition := range conditions {
				if !condition(res, url) {
					result = false

					break
				}
			}

			return result
		}
	}

	// ConditionOr return true when any of conditions is satisfied
	ConditionOr = func(conditions ...Condition) Condition {
		return func(res *http.Response, url string) bool {
			for _, condition := range conditions {
				if condition(res, url) {
					return true
				}
			}

			return false
		}
	}

	ConditionNotStatusOk = func(res *http.Response, _ string) bool {
		return res.StatusCode < 200 || res.StatusCode > 299
	}

	ConditionNotStatusBadRequest = func(res *http.Response, _ string) bool {
		return res.StatusCode != http.StatusBadRequest
	}

	ConditionNotStatusNotFound = func(res *http.Response, _ string) bool {
		return res.StatusCode != http.StatusNotFound
	}
)

//nolint:bodyclose
func DefaultSentryErrorHandler() func(res *http.Response, url string) error {
	// return GetSentryErrorHandler(ConditionAnd(sentry.ConditionNotStatusOk, ConditionNotStatusBadRequest))
	return GetSentryErrorHandler(ConditionAnd(ConditionNotStatusOk))
}
