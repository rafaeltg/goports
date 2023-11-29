package logging

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
)

type Strings []string

func (s Strings) LogValue() slog.Value {
	return slog.StringValue(
		fmt.Sprintf("[%v]", strings.Join(s, ",")),
	)
}

type Float64s []float64

func (fs Float64s) LogValue() slog.Value {
	strs := make([]string, 0, len(fs))
	for _, f := range fs {
		strs = append(strs, fmt.Sprintf("%f", f))
	}

	return slog.StringValue(
		fmt.Sprintf("[%v]", strings.Join(strs, ",")),
	)
}

func Objects[T slog.LogValuer](values []T) any {
	return objects[T](values)
}

type objects[T slog.LogValuer] []T

func (os objects[T]) LogValue() slog.Value {
	b, _ := json.Marshal(os)
	return slog.StringValue(string(b))
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
