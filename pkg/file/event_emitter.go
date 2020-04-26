package file

type EventEmitter interface {
	Close() error
	C() <-chan Event
	Watch(Directory, Path, ShouldWatchFunc, ShouldWalkFunc) error
}
