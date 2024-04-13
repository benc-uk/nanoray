package raytrace

type RaytraceError string

func (e RaytraceError) Error() string { return string(e) }

const (
	ErrInvalidRadius = RaytraceError("invalid radius")
)
