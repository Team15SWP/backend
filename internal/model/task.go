package model

type Statistics struct {
	Easy   int32 `json:"easy"`
	Medium int32 `json:"medium"`
	Hard   int32 `json:"hard"`
	Total  int32 `json:"total"`
}

type GeneratedTask struct {
	ID               int64        `json:"id"`
	UserID           int64        `json:"user_id"`
	TaskName         string       `json:"Task_name"`
	TaskDescription  string       `json:"Task_description"`
	SampleInputCases []InputCases `json:"Sample_input_cases"`
	Hints            Hints        `json:"Hints"`
	Solution         string       `json:"Full_solution"`
	Difficulty       string       `json:"Difficulty"`
	Solved           bool         `json:"Solved"`
}

type InputCases struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

type Hints struct {
	Hint1 string `json:"Hint1"`
	Hint2 string `json:"Hint2"`
	Hint3 string `json:"Hint3"`
}
