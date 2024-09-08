package chaincode

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
)

type invokeCCRequest struct {
	ChannelName string `json:"channel_name"`
	CCName      string `json:"cc_name"`
	Function    string `json:"function"`
}

type invokeFunction struct {
	FunctionName string   `json:"function_name"`
	Args         []string `json:"args"`
}

type ordererInfo struct {
	IP       string `json:"IP"`
	CertPath string `json:"cert_path"`
	Name     string `json:"name"`
}

type peerInfo struct {
	IP       string `json:"IP"`
	CertPath string `json:"cert_path"`
}

// InvokeCC
// @Summary Invoke the specified chaincode.
// @Description `peer chaincode invoke` is executed through `exec.Command()` to invoke the specified chaincode.
// @Accept json
// @Param body body invokeCCRequest true "channel name (mychannel), cc name (basic), function ('{"function":"InitLedger","Args":[]}')"
// @Produce json
// @Tags chaincode
// @Success 200 "successful operation"
// @Router /fabric/chaincode/invoke [post]
func InvokeCC(c *gin.Context) {
	var requestBody invokeCCRequest
	var orderer ordererInfo
	var peer1, peer2 peerInfo
	GOPATH := os.Getwd()
	networkPath := fmt.Sprintf("%s/hyperledger/fabric-samples/test-network", GOPATH)

	orderer.IP = "localhost:7050"
	orderer.Name = "orderer.example.com"
	orderer.CertPath = fmt.Sprintf("%s/organizations/ordererOrganizations/example.com/orderers/"+
		"orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem", networkPath)
	peer1.IP = "localhost:7051"
	peer1.CertPath = fmt.Sprintf("%s/organizations/peerOrganizations/org1.example.com/peers/"+
		"peer0.org1.example.com/tls/ca.crt", networkPath)
	peer2.IP = "localhost:9051"
	peer2.CertPath = fmt.Sprintf("%s/organizations/peerOrganizations/org2.example.com/peers/"+
		"peer0.org2.example.com/tls/ca.crt", networkPath)

	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request format."})
		return
	}

	cmd := exec.Command("peer", "chaincode", "invoke", "-o", orderer.IP,
		"--ordererTLSHostnameOverride", orderer.Name, "-C", requestBody.ChannelName, "-n", requestBody.CCName,
		"--tls", "--cafile", orderer.CertPath, "--peerAddresses", peer1.IP, "--tlsRootCertFiles", peer1.CertPath,
		"--peerAddresses", peer2.IP, "--tlsRootCertFiles", peer2.CertPath, "-c", requestBody.Function)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, string(output))
}

// QueryCC
// @Summary Query using the specified chaincode.
// @Description `peer chaincode invoke` is executed through `exec.Command()` to get endorsed result of chaincode function call and print it. It won't generate transaction.
// @Accept json
// @Param body invokeCCRequest true "channel name (mychannel), cc name (basic)"
// @Produce json
// @Tags chaincode
// @Success 200 "successful operation"
// @Router /fabric/chaincode/query [get]
func QueryCC(c *gin.Context) {
	var requestBody invokeCCRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request format."})
		return
	}

	cmd := exec.Command("peer", "chaincode", "query", "-C", requestBody.ChannelName,
		"-n", requestBody.CCName, "-c", requestBody.Function)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, string(output))
}
