package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
)

// StartFabricWChannel
// @Summary Bring up fabric network with one channel.
// @Description `network.sh up createChannel` is executed through `exec.Command()` to start the network and create channel `mychannel`.
// @Produce json
// @Tags network
// @Success 200 "successful operation"
// @Router /fabric/network/up [post]
func StartFabricWChannel(c *gin.Context) {
	GOPATH := os.Getenv("GOPATH")
	networkPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/test-network", GOPATH)

	cmd := exec.Command("bash", "network.sh", "up", "createChannel")
	cmd.Dir = networkPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Fabric network successfully started."})
}

// StopFabric
// @Summary Bring down the fabric network.
// @Description `network.sh down` is executed through `exec.Command()` to shut down the network.
// @Produce json
// @Tags network
// @Success 200 "successful operation"
// @Router /fabric/network/down [post]
func StopFabric(c *gin.Context) {
	GOPATH := os.Getenv("GOPATH")
	networkPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/test-network", GOPATH)

	cmd := exec.Command("bash", "network.sh", "down")
	cmd.Dir = networkPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Fabric network successfully shut down."})
}
