package menu

type EventType int

const (
	EventTypeCreate EventType = iota
	EventTypeUpdate
	EventTypeDelete
)

type EventRuntime int

const (
	// EventRuntimeStandard indicates that the action is a standard action
	EventRuntimeStandardd EventRuntime = iota
	// EventRuntimePre indicates that the action is a pre action
	EventRuntimePre
	// EventRuntimePost indicates that the action is a post action
	EventRuntimePost
)

type EventOrigin int

const (
	// EventOriginEnvironment indicates that the an environment was changed (created, updated, deleted)
	EventOriginEnvironment EventOrigin = iota
	// EventOriginStage indicates that the a stage was changed (created, updated, deleted)
	EventOriginStage
	// EventOriginCluster indicates that the a cluster was changed (created, updated, deleted)
	EventOriginCluster
	// EventOriginAddon indicates that an addon was changed (created, updated, deleted)
	EventOriginAddon
)

type Event struct {
	Type        EventType
	Origin      EventOrigin
	Runtime     EventRuntime
	Environment string
	Stage       string
	Cluster     string
}

func newPreCreateEvent(origin EventOrigin, environment, stage, cluster string) Event {
	return Event{
		Type:        EventTypeCreate,
		Origin:      origin,
		Runtime:     EventRuntimePre,
		Environment: environment,
		Stage:       stage,
		Cluster:     cluster,
	}
}

func newPostCreateEvent(origin EventOrigin, environment, stage, cluster string) Event {
	return Event{
		Type:        EventTypeCreate,
		Origin:      origin,
		Runtime:     EventRuntimePost,
		Environment: environment,
		Stage:       stage,
		Cluster:     cluster,
	}
}

func newPreUpdateEvent(origin EventOrigin, environment, stage, cluster string) Event {
	return Event{
		Type:        EventTypeUpdate,
		Origin:      origin,
		Runtime:     EventRuntimePre,
		Environment: environment,
		Stage:       stage,
		Cluster:     cluster,
	}
}

func newPostUpdateEvent(origin EventOrigin, environment, stage, cluster string) Event {
	return Event{
		Type:        EventTypeUpdate,
		Origin:      origin,
		Runtime:     EventRuntimePost,
		Environment: environment,
		Stage:       stage,
		Cluster:     cluster,
	}
}

func newPreDeleteEvent(origin EventOrigin, environment, stage, cluster string) Event {
	return Event{
		Type:        EventTypeDelete,
		Origin:      origin,
		Runtime:     EventRuntimePre,
		Environment: environment,
		Stage:       stage,
		Cluster:     cluster,
	}
}

func newPostDeleteEvent(origin EventOrigin, environment, stage, cluster string) Event {
	return Event{
		Type:        EventTypeDelete,
		Origin:      origin,
		Runtime:     EventRuntimePost,
		Environment: environment,
		Stage:       stage,
		Cluster:     cluster,
	}
}
