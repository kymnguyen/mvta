package loader

import (
	"os"
	"path/filepath"

	"workflow-svc/internal/domain/workflow"

	"gopkg.in/yaml.v3"
)

type YAMLLoader struct {
	workflowDir string
}

func NewYAMLLoader(workflowDir string) *YAMLLoader {
	return &YAMLLoader{workflowDir: workflowDir}
}

func (l *YAMLLoader) LoadAll() ([]*workflow.WorkflowDefinition, error) {
	var workflows []*workflow.WorkflowDefinition

	err := filepath.Walk(l.workflowDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		wf, err := l.LoadFile(path)
		if err != nil {
			return err
		}

		workflows = append(workflows, wf)
		return nil
	})

	return workflows, err
}

func (l *YAMLLoader) LoadFile(path string) (*workflow.WorkflowDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var wf workflow.WorkflowDefinition
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, err
	}

	if err := wf.Validate(); err != nil {
		return nil, err
	}

	return &wf, nil
}
