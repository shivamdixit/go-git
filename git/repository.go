package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/go-ini/ini"
)

const (
	// gitPath is the directory containing all git metadata
	gitPath = ".git"

	// configFile is the configuration file for a git repo
	configFile = "config"
)

// Repository represents a Git repository
type Repository struct {
	worktree string
	gitDir   string
	config   *ini.File
}

// absPath converts a path relative to git dir to absolute path
func (r *Repository) absPath(p string) string {
	return path.Join(r.gitDir, p)
}

// for a given path 'p' (relative to .git directory), this function
// creates all the intermediate directories, if required
// and finally returns the absolute path of the given path
func (r *Repository) dir(p string, mkdir bool) (string, error) {
	// get the full path for directory
	p = r.absPath(p)

	if f, err := os.Stat(p); !os.IsNotExist(err) {
		if f.IsDir() {
			return p, nil
		}

		return "", fmt.Errorf("'%s' is not a directory", p)
	}

	if !mkdir {
		return p, nil
	}

	err := os.MkdirAll(p, 0755)
	if err != nil {
		return "", err
	}

	return p, nil
}

// for a given path (relative to .git directory), this function
// returns the absolute path and creates intermediate directories
// if required
func (r *Repository) file(p string, mkdir bool) (string, error) {
	// create intermediate directories, if required
	_, err := r.dir(filepath.Dir(p), mkdir)
	if err != nil {
		return "", err
	}

	return r.absPath(p), nil
}

// NewRepository is the constructor of Repository object
func NewRepository(dir string, force bool) (*Repository, error) {
	r := Repository{
		worktree: dir,
		gitDir:   path.Join(dir, gitPath),
	}

	// if force flag is true, we don't do any additional checks
	// this flag is used, for example, to initialize repository
	if force {
		return &r, nil
	}

	// check if git metdata directory exists (.git)
	f, err := os.Stat(r.gitDir)
	if os.IsNotExist(err) || !f.IsDir() {
		return nil, fmt.Errorf("not a git repository: %s", dir)
	}

	// validate git config exists
	confPath := path.Join(r.gitDir, configFile)
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return nil, errors.New("cannot find git config")
	}

	// Load the configuration file
	cfg, err := ini.Load(confPath)
	if err != nil {
		return nil, err
	}

	r.config = cfg

	// Validate the configuration file
	v := cfg.Section("core").Key("repositoryformatversion").String()
	if v != "0" {
		return nil, errors.New("unsupported repository version in config")
	}

	return &r, nil
}

// validate runs all the sanity checks before creating a repository
func validate(repo *Repository) error {
	f, err := os.Stat(repo.worktree)
	// error that is not "directory not exists" error
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// directory doesn't exists
	if os.IsNotExist(err) {
		// create the directory
		e := os.Mkdir(repo.worktree, 0755)
		if e != nil {
			return e
		}
	}

	// file exists, but it is not a directory
	if err == nil && !f.IsDir() {
		return errors.New("not a directory")
	}

	// Either we have created directory, or are already given a
	// dir at this point. We need to validate  if it is empty.
	files, err := ioutil.ReadDir(repo.worktree)
	if err != nil {
		return err
	}

	if len(files) != 0 {
		return errors.New("directory not empty")
	}

	return nil
}

func initializeFile(r *Repository, relPath string, content []byte) error {
	p, err := r.file(relPath, true)
	if err != nil {
		return err
	}

	ioutil.WriteFile(p, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Create initializes an empty git repository for a given path
func Create(dst string) error {
	repo, err := NewRepository(dst, true)
	if err != nil {
		return err
	}

	// run the sanity checks
	err = validate(repo)
	if err != nil {
		return err
	}

	// initialize repo
	repo.dir("branches", true)
	repo.dir("objects", true)
	repo.dir("refs/tags", true)
	repo.dir("refs/head", true)

	// initialize necessary files
	// .git/description
	// .git/HEAD
	// .git/config
	desc := []byte("Unnamed repository; edit this file 'description' to name the repository.\n")
	head := []byte("Unnamed repository; edit this file 'description' to name the repository.\n")
	if err = initializeFile(repo, "description", desc); err != nil {
		return err
	}
	if err = initializeFile(repo, "head", head); err != nil {
		return err
	}
	if err = initializeFile(repo, "config", []byte("")); err != nil {
		return err
	}

	// write the default configuration file
	defaultConfig().SaveTo(repo.absPath("config"))

	return nil
}

// defaultConfig returns the default configuration file for a git repository
func defaultConfig() *ini.File {
	c := ini.Empty()

	sec, _ := c.NewSection("core")
	sec.Key("repositoryformatversion").SetValue("0")
	sec.Key("filemode").SetValue("false")
	sec.Key("bare").SetValue("false")

	return c
}
