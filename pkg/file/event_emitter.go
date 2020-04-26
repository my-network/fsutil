package file

type EventEmitter interface {
	Close() error
	C() <-chan Event
	Watch(Object, ShouldWatchFunc, ShouldWalkFunc) error
}
