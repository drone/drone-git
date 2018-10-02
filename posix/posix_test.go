package posix

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommits(t *testing.T) {
	remote := "/tmp/remote/greeting"

	base, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(base)

	for i, test := range tests {
		local := filepath.Join(base, fmt.Sprint(i))
		err = os.MkdirAll(local, 0777)
		if err != nil {
			t.Error(err)
			return
		}

		bin, err := filepath.Abs("clone-commit")
		if err != nil {
			t.Error(err)
			return
		}

		cmd := exec.Command(bin)
		cmd.Dir = local
		cmd.Env = []string{
			fmt.Sprintf("DRONE_COMMIT_BRANCH=%s", test.branch),
			fmt.Sprintf("DRONE_COMMIT_SHA=%s", test.commit),
			fmt.Sprintf("DRONE_WORKSPACE=%s", local),
			fmt.Sprintf("DRONE_REMOTE_URL=%s", remote),
		}

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Error(err)
			t.Log(string(out))
			return
		}

		commit, err := getCommit(local)
		if err != nil {
			t.Error(err)
			return
		}

		branch, err := getBranch(local)
		if err != nil {
			t.Error(err)
			return
		}

		if want, got := test.commit, commit; got != want {
			t.Errorf("Want commit %s, got %s", want, got)
		}

		if want, got := test.branch, branch; got != want {
			t.Errorf("Want branch %s, got %s", want, got)
		}

		file := filepath.Join(local, test.file)
		out, err = ioutil.ReadFile(file)
		if err != nil {
			t.Error(err)
			return
		}

		if want, got := test.text, string(out); want != got {
			t.Errorf("Want file content %q, got %q", want, got)
		}
	}
}

func TestTags(t *testing.T) {
	remote := "/tmp/remote/greeting"

	base, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(base)

	for i, test := range tests {
		local := filepath.Join(base, fmt.Sprint(i))
		err = os.MkdirAll(local, 0777)
		if err != nil {
			t.Error(err)
			return
		}

		bin, err := filepath.Abs("clone-tag")
		if err != nil {
			t.Error(err)
			return
		}

		cmd := exec.Command(bin)
		cmd.Dir = local
		cmd.Env = []string{
			fmt.Sprintf("DRONE_TAG=%s", test.tag),
			fmt.Sprintf("DRONE_COMMIT_SHA=%s", test.commit),
			fmt.Sprintf("DRONE_WORKSPACE=%s", local),
			fmt.Sprintf("DRONE_REMOTE_URL=%s", remote),
		}

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Error(err)
			t.Log(string(out))
			return
		}

		commit, err := getCommit(local)
		if err != nil {
			t.Error(err)
			return
		}

		if want, got := test.commit, commit; got != want {
			t.Errorf("Want commit %s, got %s", want, got)
		}

		file := filepath.Join(local, test.file)
		out, err = ioutil.ReadFile(file)
		if err != nil {
			t.Error(err)
			return
		}

		if want, got := test.text, string(out); want != got {
			t.Errorf("Want file content %q, got %q", want, got)
		}
	}
}

func TestPullRequest(t *testing.T) {
	remote := "https://github.com/octocat/Spoon-Knife.git"

	base, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(base)

	for i, test := range testsPR {
		local := filepath.Join(base, fmt.Sprint(i))
		err = os.MkdirAll(local, 0777)
		if err != nil {
			t.Error(err)
			return
		}

		bin, err := filepath.Abs("clone-pull-request")
		if err != nil {
			t.Error(err)
			return
		}

		cmd := exec.Command(bin)
		cmd.Dir = local
		cmd.Env = []string{
			fmt.Sprintf("DRONE_COMMIT_REF=refs/pull/%s/head", test.pull_request),
			fmt.Sprintf("DRONE_COMMIT_BRANCH=%s", test.branch),
			fmt.Sprintf("DRONE_COMMIT_SHA=%s", test.commit),
			fmt.Sprintf("DRONE_WORKSPACE=%s", local),
			fmt.Sprintf("DRONE_REMOTE_URL=%s", remote),
			fmt.Sprintf("GIT_AUTHOR_NAME=%s", "drone"),
			fmt.Sprintf("GIT_AUTHOR_EMAIL=%s", "drone@localhost"),
			fmt.Sprintf("GIT_COMMITTER_NAME=%s", "drone"),
			fmt.Sprintf("GIT_COMMITTER_EMAIL=%s", "drone@localhost"),
		}

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Error(err)
			t.Log(string(out))
			return
		}

		head_commit, err := getCommit(local)
		if err != nil {
			t.Error(err)
			return
		}

		branch, err := getBranch(local)
		if err != nil {
			t.Error(err)
			return
		}

		// The local repo HEAD must be the PR destination branch's local checkout
		if want, got := test.branch, branch; got != want {
			t.Errorf("Want branch %s, got %s", want, got)
		}

		if test.require_ff {
			// Require a fast-forward merge if the PR source branch is a descendant of the PR destination branch
			if want, got := test.commit, head_commit; got != want {
				t.Errorf("Want commit %s, got %s", want, got)
			}
		} else {
			// Otherwise the PR branch's HEAD must still be an ancestor of the local repo HEAD
			if !checkAncestry(local, test.commit, head_commit) {
				t.Errorf("PR branch HEAD %s is not an ancestor of local HEAD %s", test.commit, head_commit)
			}
		}

		file := filepath.Join(local, "directory/file.txt")
		out, err = ioutil.ReadFile(file)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func getBranch(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func getCommit(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func checkAncestry(path string, ancestor string, descendant string) (bool) {
	cmd := exec.Command("git", "merge-base", "--is-ancestor", ancestor, descendant)
	cmd.Dir = path
	_, err := cmd.CombinedOutput()
	return err == nil
}

var tests = []struct {
	branch string
	commit string
	tag    string
	file   string
	text   string
}{
	{
		commit: "9cd29dca0a98f76df94d66493ee54788a18190a0",
		branch: "master",
		tag:    "v1.0.0",
		file:   "hello.txt",
		text:   "hi world\n",
	},
	{
		commit: "bbdf5d4028a6066431f59fcd8d83afff610a55ae",
		branch: "master",
		tag:    "v1.1.0",
		file:   "hello.txt",
		text:   "hello world\n",
	},
	{
		commit: "553af1ca53c9ad54b096d7ff1416f6c4d1e5049f",
		branch: "fr",
		tag:    "v2.0.0",
		file:   "hello.txt",
		text:   "salut monde\n",
	},
	{
		commit: "94b4a1710d1581b8b00c5f7b077026eae3c07646",
		branch: "fr",
		tag:    "v2.1.0",
		file:   "hello.txt",
		text:   "bonjour monde\n",
	},
}

var testsPR = []struct {
	pull_request string
	branch       string
	commit       string
	require_ff   bool
}{
	{
		pull_request: "14596",
		branch:       "master",
		commit:       "26923a8f37933ccc23943de0d4ebd53908268582",
		require_ff:   true,
	},
	{
		pull_request: "16667",
		branch:       "master",
		commit:       "e8d96e48cffd34479c3ecd0410d98c70ee9bbd7f",
		require_ff:   false,
	},
}
