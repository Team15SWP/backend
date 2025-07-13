package service

import (
	"context"
	"encoding/json"
	"fmt"

	"study_buddy/internal/config"
	"study_buddy/internal/model"
	"study_buddy/pkg/llm"
)

var _ Service = (*TaskService)(nil)

type TaskService struct {
	repo      TaskProvider
	statsRepo StatsProvider
	openAi    *config.OpenAI
	prompts   *config.Prompts
	LLM       llm.Client
}

func NewTaskService(repo TaskProvider, statsRepo StatsProvider, llmClient llm.Client, openAi *config.OpenAI, prompts *config.Prompts) *TaskService {
	return &TaskService{
		repo:      repo,
		statsRepo: statsRepo,
		openAi:    openAi,
		prompts:   prompts,
		LLM:       llmClient,
	}
}

type Service interface {
	GenerateTask(ctx context.Context, userId int64, topic, difficulty string) (*model.GeneratedTask, error)
	EvaluateCodeForTask(ctx context.Context, task, code string) (*Question, error)
	GetStatistics(ctx context.Context, userId int64) (*model.Statistics, error)
}

type TaskProvider interface {
	CreateTask(ctx context.Context, task *model.GeneratedTask) error
}

type StatsProvider interface {
	GetStatisticsData(ctx context.Context, userId int64) (*model.Statistics, error)
}

func (t *TaskService) GenerateTask(ctx context.Context, userId int64, topic, difficulty string) (*model.GeneratedTask, error) {
	prompt := fmt.Sprintf(t.prompts.GenerateTask, topic, difficulty)
	response, err := t.LLM.Complete(ctx, prompt)
	if err != nil {
		return nil, err
	}
	var task *model.GeneratedTask
	err = json.Unmarshal([]byte(response), &task)
	if err == nil && task != nil {
		task.Difficulty = difficulty
		task.UserID = userId
		task.Solved = false
		err = t.repo.CreateTask(ctx, task)
		fmt.Println(task)
		if err != nil {
			return nil, fmt.Errorf("t.repo.CreateTask: %w", err)
		}
	}
	return task, nil
}

type EvaluateCodeResponse struct {
	Feedback string `json:"feedback"`
}

type Question struct {
	Request  string `json:"question"`
	Verdict  string `json:"correct"`
	Feedback string `json:"feedback"`
}

func (t *TaskService) EvaluateCodeForTask(ctx context.Context, task, code string) (*Question, error) {
	prompt := fmt.Sprintf(t.prompts.CheckCodeForTask, task, code)

	response, err := t.LLM.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("t.LLM.Complete: %w", err)
	}

	var feedback *Question
	err = json.Unmarshal([]byte(response), &feedback)
	if err == nil {
	}

	fmt.Println(feedback)

	return feedback, nil
}

func (t *TaskService) GetStatistics(ctx context.Context, userId int64) (*model.Statistics, error) {
	response, err := t.statsRepo.GetStatisticsData(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("t.statsRepo: %w", err)
	}
	return response, nil
}
