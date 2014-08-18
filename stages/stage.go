package stages

type Stage interface {
	Run() bool
}

func InitStage(stageType string) Stage {
	switch stageType {
	case "command":
		return new(CommandStage)
	}
	return nil
}
