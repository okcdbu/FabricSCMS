package lifecycle

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
)

type installCCRequest struct {
	PackageName string `json:"package_name"`
}

type packageCCRequest struct {
	PackageName  string `json:"package_name"`
	Label        string `json:"label"`
	Language     string `json:"language"`
	CCSourceName string `json:"cc_source_name"`
}

type approveCCRequest struct {
	ChannelName string `json:"channel_name"`
	CCName      string `json:"cc_name"`
	CCVersion   string `json:"cc_version"`
	CCSequence  int32  `json:"cc_sequence"`
	PackageID   string `json:"package_ID"`
}

type commitCCRequest struct {
	ChannelName string `json:"channel_name"`
	CCName      string `json:"cc_name"`
	CCVersion   string `json:"cc_version"`
	CCSequence  int32  `json:"cc_sequence"`
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

// PackageCC
// @Summary Package a cc.
// @Description `peer lifecycle chaincode install` is executed through `exec.Command()` to install chaincode on a peer.
// @Accept json
// @Param body body packageCCRequest true "name of the cc to package (e.g. asset-transfer-basic), the language it is written in, and the label and package name for the cc once packaging is done"
// @Produce json
// @Tags lifecycle
// @Success 200 "successful operation"
// @Router /fabric/lifecycle/package [post]
func PackageCC(c *gin.Context) {
	var requestBody packageCCRequest
	var packageLanguage string // fabric uses different language names, go -> golang, js -> node, ts -> node
	GOPATH := os.Getenv("GOPATH")
	rootPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/", GOPATH)
	repoPath := fmt.Sprintf("%s/GurkhaContracts/", rootPath)
	packageStoragePath := fmt.Sprintf("%s/test-network/cc-packages/", rootPath)

	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request format."})
		return
	}

	switch requestBody.Language {
	case "go", "golang":
		packageLanguage = "golang"
	case "java":
		packageLanguage = "java"
	case "javascript", "typescript":
		packageLanguage = "node"
	}

	ccSourcePath := fmt.Sprintf("%s/%s", repoPath, requestBody.CCSourceName)
	exists, err := fileExists(ccSourcePath)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if !exists {
		c.IndentedJSON(http.StatusNotFound,
			gin.H{"message": fmt.Sprintf("Chaincode %s does not exist", requestBody.CCSourceName)})
		return
	}

	ccSourcePath = fmt.Sprintf("%s/chaincode-%s", ccSourcePath, requestBody.Language)
	exists, err = fileExists(ccSourcePath)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if !exists {
		c.IndentedJSON(http.StatusNotFound,
			gin.H{"message": fmt.Sprintf("Chaincode in language %s does not exist", requestBody.Language)})
		return
	}

	cmd := exec.Command("peer", "lifecycle", "chaincode", "package", packageStoragePath+requestBody.PackageName,
		"--path", ccSourcePath, "--lang", packageLanguage, "--label", requestBody.Label)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "CC successfully packaged."})
}

// InstallCC
// @Summary Install a cc.
// @Description `peer lifecycle chaincode install` is executed through `exec.Command()` to install chaincode on a peer.
// @Accept json
// @Param package_name path string true "name of the package to install (e.g. basic.tar.gz)"
// @Produce json
// @Tags lifecycle
// @Success 200 "successful operation"
// @Router /fabric/lifecycle/install/{package_name} [post]
func InstallCC(c *gin.Context) {
	GOPATH := os.Getenv("GOPATH")
	rootPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/", GOPATH)
	packageStoragePath := fmt.Sprintf("%s/test-network/cc-packages/", rootPath)

	packageNameParameter := c.Param("package_name")
	ccPackagePath := fmt.Sprintf("%s/%s", packageStoragePath, packageNameParameter)

	fileExists, err := fileExists(ccPackagePath)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if !fileExists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Package %s does not exist", packageNameParameter)})
		return
	}

	cmd := exec.Command("peer", "lifecycle", "chaincode", "install", ccPackagePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Package successfully installed."})
}

// ApproveCC
// @Summary Approve the cc definition for the current org.
// @Description `peer lifecycle chaincode approveformyorg` is executed through `exec.Command()` to approve a chaincode definition.
// @Accept json
// @Param body body approveCCRequest true "channel name (mychannel), cc name (basic), cc version (1.0), cc sequence (1), package ID (run [GET] /fabric/lifecycle/install)"
// @Produce json
// @Tags lifecycle
// @Success 200 "successful operation"
// @Router /fabric/lifecycle/approve [post]
func ApproveCC(c *gin.Context) {
	var requestBody approveCCRequest
	GOPATH := os.Getenv("GOPATH")
	networkPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/test-network", GOPATH)
	ordererIP := "localhost:7050"
	ordererName := "orderer.example.com"
	ordererCertPath := fmt.Sprintf("%s/organizations/ordererOrganizations/example.com/orderers/"+
		"orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem", networkPath)

	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request format."})
		return
	}

	cmd := exec.Command("peer", "lifecycle", "chaincode", "approveformyorg", "-o", ordererIP,
		"--ordererTLSHostnameOverride", ordererName, "--channelID", requestBody.ChannelName, "--name",
		requestBody.CCName, "--version", requestBody.CCVersion, "--package-id", requestBody.PackageID,
		"--sequence", fmt.Sprint(requestBody.CCSequence), "--tls", "--cafile", ordererCertPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}
	envAdmin := os.Getenv("CORE_PEER_ADMIN")
	successResponseMessage := fmt.Sprintf("CC definition of %s successfully approved for organization %s",
		requestBody.CCName, envAdmin)
	c.IndentedJSON(http.StatusOK, gin.H{"message": successResponseMessage})
}

// CommitCC
// @Summary Commit the chaincode definition on the channel.
// @Description `peer lifecycle chaincode commit` is executed through `exec.Command()` to commit chaincode definition on a channel.
// @Accept json
// @Param body body commitCCRequest true "channel name (mychannel), cc name (basic), cc version (1.0), cc sequence (1)"
// @Produce json
// @Tags lifecycle
// @Success 200 "successful operation"
// @Router /fabric/lifecycle/commit [post]
func CommitCC(c *gin.Context) {
	var requestBody ccApprovalRequest
	var orderer ordererInfo
	var peer1, peer2 peerInfo
	GOPATH := os.Getenv("GOPATH")
	networkPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/test-network", GOPATH)

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

	os.Setenv("CC_NAME", requestBody.ChannelName)
	os.Setenv("CC_SEQUENCE", string(requestBody.CCSequence))
	os.Setenv("CC_VERSION", requestBody.CCVersion)
	os.Setenv("CHANNEL_NAME", requestBody.ChannelName)

	cmd := exec.Command("peer", "lifecycle", "chaincode", "commit", "-o", orderer.IP,
		"--ordererTLSHostnameOverride", orderer.Name, "--channelID", requestBody.ChannelName, "--name", requestBody.CCName,
		"--version", requestBody.CCVersion, "--sequence", fmt.Sprint(requestBody.CCSequence), "--tls",
		"--cafile", orderer.CertPath, "--peerAddresses", peer1.IP, "--tlsRootCertFiles", peer1.CertPath,
		"--peerAddresses", peer2.IP, "--tlsRootCertFiles", peer2.CertPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, string(output))
}

// fileExists checks if the requested file exists in test-network's directory.
func fileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
