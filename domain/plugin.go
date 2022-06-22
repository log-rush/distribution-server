package domain

type LogPlugin interface {
	HandleLog(streamId string, log Log)
}
