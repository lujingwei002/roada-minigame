package cli

import "errors"

// Errors that could be occurred during message handling.
var (
	//ErrSessionOnNotify    = errors.New("current session working on notify mode")
	ErrCloseClosedSession = errors.New("close closed session")
	//ErrInvalidRegisterReq = errors.New("invalid register request")
	// ErrBrokenPipe represents the low-level connection has broken.
	ErrBrokenPipe = errors.New("broken low-level pipe")
	ErrHandShake  = errors.New("handshake failed")
	// ErrBufferExceed indicates that the current session buffer is full and
	// can not receive more data.
	ErrBufferExceed       = errors.New("session send buffer exceed")
	ErrCloseClosedGroup   = errors.New("close closed group")
	ErrClosedGroup        = errors.New("group closed")
	ErrMemberNotFound     = errors.New("member not found in the group")
	ErrSessionDuplication = errors.New("session has existed in the current group")
)
