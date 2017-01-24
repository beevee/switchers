package firebase

import "github.com/beevee/switchers"

// TaskRepository persists task information in Firebase
type TaskRepository struct {
	Repository
}

// GetAllGatheringTasks retrieves all gathering tasks
func (tr *TaskRepository) GetAllGatheringTasks() ([]switchers.GatheringTask, error) {
	ref, err := tr.firebase.Ref("gathering_tasks")
	if err != nil {
		return nil, err
	}

	var tasks []switchers.GatheringTask
	if err = ref.Value(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
