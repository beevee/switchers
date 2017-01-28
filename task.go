package switchers

// GatheringTask is a task to gather team
type GatheringTask struct {
	Text             string
	TimeLimitMinutes int
}

// ActualTask is a task for team that gathered
type ActualTask struct {
	Text             string
	TimeLimitMinutes int
	CorrectAnswers   []string
}
