package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bafbi/tableau/internal/domain"
)

const (
	DefaultTableauDir = ".tableau"
	TasksDir          = "tasks"
	MetaFile          = ".meta"
)

type FSRepository struct {
	RootDir string
	DirName string
}

func NewFSRepository(rootDir string) *FSRepository {
	dirName := os.Getenv("TABLEAU_DIR")
	if dirName == "" {
		dirName = DefaultTableauDir
	}
	return &FSRepository{
		RootDir: rootDir,
		DirName: dirName,
	}
}

func (r *FSRepository) Init() error {
	tasksPath := filepath.Join(r.RootDir, r.DirName, TasksDir)
	if err := os.MkdirAll(tasksPath, 0755); err != nil {
		return err
	}
	metaPath := filepath.Join(r.RootDir, r.DirName, MetaFile)
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		return os.WriteFile(metaPath, []byte("0"), 0644)
	}
	return nil
}

func (r *FSRepository) getNextID() (int, error) {
	metaPath := filepath.Join(r.RootDir, r.DirName, MetaFile)
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}
	nextID := id + 1
	if err := os.WriteFile(metaPath, []byte(strconv.Itoa(nextID)), 0644); err != nil {
		return 0, err
	}
	return nextID, nil
}

func (r *FSRepository) Create(task *domain.Task) error {
	id, err := r.getNextID()
	if err != nil {
		return err
	}
	task.ID = id
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	task.UpdatedAt = time.Now()

	slug := strings.ToLower(strings.ReplaceAll(task.Title, " ", "-"))
	filename := fmt.Sprintf("%d-%s.md", task.ID, slug)
	path := filepath.Join(r.RootDir, r.DirName, TasksDir, filename)

	data, err := MarshalTask(task)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (r *FSRepository) List() ([]domain.Task, error) {
	tasksPath := filepath.Join(r.RootDir, r.DirName, TasksDir)
	files, err := os.ReadDir(tasksPath)
	if err != nil {
		return nil, err
	}

	var tasks []domain.Task
	for _, f := range files {
		if filepath.Ext(f.Name()) != ".md" {
			continue
		}
		path := filepath.Join(tasksPath, f.Name())
		file, err := os.Open(path)
		if err != nil {
			continue // Skip unreadable files
		}
		task, err := ParseTask(file)
		file.Close()
		if err != nil {
			continue // Skip unparseable files
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (r *FSRepository) Get(id int) (*domain.Task, error) {
	tasks, err := r.List()
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("task %d not found", id)
}

func (r *FSRepository) Update(task *domain.Task) error {
	tasksPath := filepath.Join(r.RootDir, r.DirName, TasksDir)
	files, err := os.ReadDir(tasksPath)
	if err != nil {
		return err
	}
	
	var oldPath string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), fmt.Sprintf("%d-", task.ID)) {
			oldPath = filepath.Join(tasksPath, f.Name())
			break
		}
	}
	
	if oldPath == "" {
		return fmt.Errorf("task %d not found", task.ID)
	}
	
	task.UpdatedAt = time.Now()
	data, err := MarshalTask(task)
	if err != nil {
		return err
	}
	
	return os.WriteFile(oldPath, data, 0644)
}
