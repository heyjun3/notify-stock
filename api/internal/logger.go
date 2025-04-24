package notifystock

import (
	"log/slog"
	"os"
)

var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
