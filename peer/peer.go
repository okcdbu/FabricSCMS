package peer

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strings"
)

// peerVersion represents version information of peer
type peerVersion struct {
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
}

// GetPeerVersion checks the current version of peer binary
// @Summary Get the current peer binary version
// @Description `peer version` is executed through `exec.Command()` to return the current peer version.
// @Produce json
// @Tags peer
// @Success 200 {object} peerVersion "successful operation"
// @Router /fabric/peer/ [get]
func GetPeerVersion(c *gin.Context) {
	var versionResponse peerVersion

	cmd := exec.Command("peer", "version")
	output, _ := cmd.Output()

	outputList := strings.Split(string(output), "\n")
	version := strings.SplitAfter(outputList[1], ":")[1][1:]      // "Version: 2.4.0" -> "2.4.0"
	architecture := strings.SplitAfter(outputList[4], ":")[1][1:] // "OS/Arch: darwin/amd64" -> "darwin/amd64"

	versionResponse.Version = version
	versionResponse.Architecture = architecture

	c.IndentedJSON(http.StatusOK, versionResponse)
}
