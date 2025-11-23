package vmiddleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/flametest/vita/vhttp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type bodyDumpWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpWriter) Write(b []byte) (int, error) {
	_, err := w.Writer.Write(b)
	if err != nil {
		return 0, err
	}
	return w.ResponseWriter.Write(b)
}

func APILogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		cfg := middleware.DefaultLoggerConfig
		cfg.Format = `{"time":"${time_rfc3339_nano}","id":"${request_id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","latency":"${latency}","latency_human":"${latency_human}",` +
			`"user_agent":"${user_agent}","status":${status},"request_body":${request_body},"error":"${error}",` +
			`"response_body":"${response_body}","bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n"
		re, _ := regexp.Compile("health-check")
		cfg.Skipper = func(context echo.Context) bool {
			if re.Match([]byte(context.Request().URL.Path)) {
				return true
			}
			return false
		}

		return func(c echo.Context) error {
			start := time.Now()

			// record request body
			var reqBody []byte
			var reqBodyStr string
			if c.Request().Body != nil {
				reqBody, _ = io.ReadAll(c.Request().Body)
				// compress the request json body
				var jsonData map[string]interface{}
				reqBodyStr = string(reqBody)
				if err := json.Unmarshal(reqBody, &jsonData); err == nil {
					if compact, err := json.Marshal(jsonData); err == nil {
						reqBodyStr = string(compact)
					} else {
						reqBodyStr = string(reqBody)
					}
				}
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody))

			// record response body
			respBody := new(bytes.Buffer)
			origWriter := c.Response().Writer
			c.Response().Writer = &bodyDumpWriter{Writer: respBody, ResponseWriter: origWriter}

			var errValue string
			e := next(c)
			if e != nil {
				errValue = e.Error()
			}
			duration := time.Since(start)
			output := cfg.Format
			output = strings.ReplaceAll(output, "${time_rfc3339_nano}", time.Now().Format(time.RFC3339Nano))
			output = strings.ReplaceAll(output, "${request_id}", c.Request().Header.Get(vhttp.HeaderXRequestID))
			output = strings.ReplaceAll(output, "${remote_ip}", c.RealIP())
			output = strings.ReplaceAll(output, "${host}", c.Request().Host)
			output = strings.ReplaceAll(output, "${method}", c.Request().Method)
			output = strings.ReplaceAll(output, "${uri}", c.Request().RequestURI)
			output = strings.ReplaceAll(output, "${user_agent}", c.Request().UserAgent())
			output = strings.ReplaceAll(output, "${status}", strconv.Itoa(c.Response().Status))
			output = strings.ReplaceAll(output, "${request_body}", strconv.Quote(reqBodyStr))
			output = strings.ReplaceAll(output, "${error}", errValue)
			output = strings.ReplaceAll(output, "${latency}", duration.String())
			output = strings.ReplaceAll(output, "${latency_human}", duration.String())
			output = strings.ReplaceAll(output, "${bytes_in}", strconv.FormatInt(c.Request().ContentLength, 10))
			output = strings.ReplaceAll(output, "${bytes_out}", strconv.FormatInt(c.Response().Size, 10))
			output = strings.ReplaceAll(output, "${response_body}", respBody.String())

			if cfg.Output == nil {
				_, err := c.Logger().Output().Write([]byte(output))
				if err != nil {
					return err
				}
				return e
			}
			_, err := cfg.Output.Write([]byte(output))
			if err != nil {
				return err
			}
			return e
		}
	}
}
