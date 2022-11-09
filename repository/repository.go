package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type cloneRequest struct {
	Url       string `json:"url"`
	Directory string `json:"directory"`
}

func AddRemote(c *gin.Context) {
	// git remote add origin git@3.34.46.252:remote_fabric.git
}

// CloneCC
// @Summary Clone a repository.
// @Description Clone a repository.
// @Accept json
// @Param body body cloneRequest true "url (https://github.com/arogyaGurkha/GurkhaContracts.git), directory (GurkhaContracts or nil)"
// @Produce json
// @Tags repository
// @Success 200 "successful operation"
// @Router /fabric/repository/clone [get]
func CloneCC(c *gin.Context) {
	var requestBody cloneRequest
	GOPATH := os.Getenv("GOPATH")
	rootPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/", GOPATH)

	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request format."})
		return
	}

	requestBody.Directory = fmt.Sprintf("%s/%s", rootPath, requestBody.Directory)

	r, err := git.PlainClone(requestBody.Directory, false, &git.CloneOptions{
		URL:               requestBody.Url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": commit})
}

func AddChanges(c *gin.Context) {
	// git add .
}

func CommitChanges(c *gin.Context) {
	// git commit
}

func PushChanges(c *gin.Context) {
	// git push origin master
}

func FetchOrigin(c *gin.Context) {

}

// RevertUpdate
// @Summary Revert most recent update.
// @Description Revert most recent update.
// @Produce json
// @Tags repository
// @Success 200 "successful operation"
// @Router /fabric/repository/revert [post]
func RevertUpdate(c *gin.Context) {
	GOPATH := os.Getenv("GOPATH")
	repoPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/GurkhaContracts", GOPATH)

	cmd := exec.Command("git", "reset", "--hard", "HEAD@{1}")
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Path message": err.Error()})
		return
	}
	ref, err := r.Head()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Head": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": ref.Hash().String()})
}

// CheckUpdate
// @Summary Show incoming changes.
// @Description `git log HEAD..origin/main --oneline` is executed through `exec.Command()` to print incoming changes.
// @Produce json
// @Tags repository
// @Success 200 "successful operation"
// @Router /fabric/repository/updates [get]
func CheckUpdate(c *gin.Context) {
	// git log HEAD..origin/main --oneline
	GOPATH := os.Getenv("GOPATH")
	repoPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/GurkhaContracts", GOPATH)

	cmd := exec.Command("git", "fetch", "--all")
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"Fetch Error": errMessage})
		return
	}

	cmd = exec.Command("git", "log", "HEAD..origin/main", "--oneline")
	cmd.Dir = repoPath

	output, err = cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"Logging Error": errMessage})
		return
	}

	outputList := strings.Split(string(output), "\n")

	c.IndentedJSON(http.StatusOK,
		gin.H{"Incoming Updates": len(outputList) - 1, "Updates": outputList[:len(outputList)-1]})
}

// PullOrigin
// @Summary Pull changes from a remote repository.
// @Description Pull changes from a remote repository.
// @Accept json
// @Produce json
// @Tags repository
// @Success 200 "successful operation"
// @Router /fabric/repository/pull [get]
func PullOrigin(c *gin.Context) {
	// git pull origin
	GOPATH := os.Getenv("GOPATH")
	repoPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/GurkhaContracts", GOPATH)

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Path message": err.Error()})
		return
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Worktree message": err.Error()})
		return
	}

	// Pull the latest changes from the origin
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Pull message": err.Error()})
		return
	}

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"git Head message": err.Error()})
		return
	}

	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"Commit Hash message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": commit})
}

// ResetLocal
// @Summary Reset local repository.
// @Description `git fetch`, `git reset --hard`, `git clean -xdf` is executed through `exec.Command()` to reset local repository.
// @Produce json
// @Tags repository
// @Success 200 "successful operation"
// @Router /fabric/repository/reset [post]
func ResetLocal(c *gin.Context) {
	// git fetch
	// git reset --hard
	// git clean -xdf

	GOPATH := os.Getenv("GOPATH")
	repoPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/GurkhaContracts", GOPATH)

	cmd := exec.Command("git", "fetch")
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"Fetch Error": errMessage})
		return
	}

	cmd = exec.Command("git", "reset", "--hard")
	cmd.Dir = repoPath

	output, err = cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"Reset Error": errMessage})
		return
	}

	cmd = exec.Command("git", "clean", "-xdf")
	cmd.Dir = repoPath

	output, err = cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"Clean Error": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Local repository successfully reset."})
}

// GetRefLogs
// @Summary Show the reflog.
// @Description `git reflog` is executed through `exec.Command()` to show the reflogs.
// @Produce json
// @Tags repository
// @Success 200 "successful operation"
// @Router /fabric/repository/logs [get]
func GetRefLogs(c *gin.Context) {
	// git reflog
	GOPATH := os.Getenv("GOPATH")
	repoPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/GurkhaContracts", GOPATH)

	cmd := exec.Command("git", "reflog")
	cmd.Dir = repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	outputList := strings.Split(string(output), "\n")

	c.IndentedJSON(http.StatusOK,
		gin.H{"Log Count": len(outputList) - 1, "Logs": outputList[:len(outputList)-1]})
}
