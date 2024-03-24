package loghandler

import (
	"context"
	"errors"

	"github.com/andyinabox/xboxcrelay/pkg/xboxc"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func HandleInput(ctx context.Context, state *xboxc.State) error {
	return errors.New("not implemented")
}
