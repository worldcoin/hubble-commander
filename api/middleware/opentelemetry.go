package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Worldcoin/hubble-commander/o11y"
	"github.com/felixge/httpsnoop"
	log "github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func extractMethod(r *http.Request) *string {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("API: failed to read request body: %s", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var decoded payload

	err = json.Unmarshal(body, &decoded)
	if err != nil {
		// TODO: this might be a batch request, we should try to interpret those
		log.Warn("received JSON request without interpretable method, not creating a span")
		return nil
	}

	return &decoded.Method
}

func OpenTelemetryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := extractMethod(r)

		if method == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx, span := otel.Tracer("rpc.call").Start(r.Context(), *method)
		defer span.End()
		r = r.WithContext(ctx)

		intercepted := make([]byte, 0)
		w = httpsnoop.Wrap(w, httpsnoop.Hooks {
			Write: func(original httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				return func(b []byte) (int, error) {
					intercepted = append(intercepted, b...)
					return original(b)
				}
			},
		})

		next.ServeHTTP(w, r)

		var response jsonrpc.RPCResponse
		err := json.Unmarshal(intercepted, &response)
		if err != nil {
			// TODO: should we annotate the span with an error here, to make
			//       these easier to search for?
			log.WithFields(o11y.TraceFields(ctx)).Warn("received uninterpretable JSON reponse, not annotating span with status")
			return
		}

		if response.Error != nil {
			span.RecordError(response.Error)
			span.SetStatus(codes.Error, response.Error.Error())
			span.SetAttributes(attribute.Int("errorCode", response.Error.Code))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	})
}
