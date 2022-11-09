package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
)

func AddOrg(c *gin.Context) {
	GOPATH := os.Getenv("GOPATH")
	networkPath := fmt.Sprintf("%s/src/github.com/jhl8109/fabric-samples/test-network/addOrg3", GOPATH)

	cmd := exec.Command("bash", "addOrg3.sh", "up")
	cmd.Dir = networkPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "A new organization successfully started."})
}
