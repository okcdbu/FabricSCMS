package dashboard

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/okcdbu/FabricSCMS/admin"
	"github.com/okcdbu/FabricSCMS/repository/search"
)

type installCC struct {
	CCName     string `json:"cc_name"`
	CCPath     string `json:"cc_path"`
	CCLanguage string `json:"cc_language"`
}

type transactionResponse struct {
	Peer string `json:"peer"`
}

type assetQueryResponse struct {
	Assets []*Asset `json:""`
}

type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

var (
	peerAddressOrg1 = "localhost:7051"
	peerAddressOrg2 = "localhost:9051"
	GOPATH          = "/home/ubuntu"
	networkPath     = fmt.Sprintf("%s/fabric-samples/test-network", GOPATH)
	scriptPath      = fmt.Sprintf("%s/fabric-samples/test-network/scripts", GOPATH)
	now             = time.Now()
	assetId         = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond()/1e6))
	CCPATHROOT      = "/home/ubuntu/chaincodes"
)

// FileUpload
// @Summary If clients send cc package file, then upload zip file at /Downloads/chaincodes and install system channel.
// @Produce json
// @Tags dashboard
// @Success 200 "successful operation"
// @Router /fabric/dashboard/smart-contracts/file [post]
func FileUpload(c *gin.Context) {
	// get zip file and upload it
	rawData := c.PostForm("data")
	var inputData search.Article
	json.Unmarshal([]byte(rawData), &inputData)
	file, _ := c.FormFile("file")
	c.SaveUploadedFile(file, fmt.Sprintf(`%s/%s.tar.gz`, CCPATHROOT, inputData.Name))
	log.Println("zip file uploaded successfully")

	// save smart constarct info
	inputData.UploadDate = fmt.Sprintf(time.Now().UTC().Format("2006-01-02"))
	log.Println(inputData.UploadDate)
	res, err := search.AddDocumentToES(&inputData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": res})
}

// InstallWithDeployCC
// @Summary Install specified CC using deployCC script.
// @Produce json
// @Tags dashboard
// @Success 200 "successful operation"
// @Router /fabric/dashboard/deployCC [post]
func InstallWithDeployCC(c *gin.Context) {

	var requestBody installCC
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request format."})
		return
	}

	setEnv()
	finalCCPath := fmt.Sprintf("%s/%s", CCPATHROOT, requestBody.CCPath)
	log.Println("Deploying chaincode...")
	log.Println(fmt.Sprintf("CCName :%s, ccPath :%s, finalCCPath : %s", requestBody.CCName, requestBody.CCPath, finalCCPath))
	cmd := exec.Command("bash", "./scripts/deployTarCC.sh", "mychannel", requestBody.CCName, requestBody.CCLanguage)
	cmd.Dir = networkPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		log.Println(errMessage)
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "CC Installed"})
}

// DownloadCC
// @Summary Download specified CC
// @Produce json
// @Tags dashboard
// @Success 200 "successful operation"
// @Router /fabric/dashboard/downloadCC [GET]
func DownloadCC(c *gin.Context) {
	fileName := c.Query("name")
	log.Println(fileName)
	targetPath := filepath.Join(CCPATHROOT, fmt.Sprintf("%s.tar.gz", fileName))
	log.Println(targetPath)
	//This ckeck is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
	if !strings.HasPrefix(filepath.Clean(targetPath), CCPATHROOT) {
		c.String(403, "Look like you attacking me")
		return
	}
	//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(targetPath)

	// c.IndentedJSON(http.StatusOK, gin.H{"message": "CC Download successfully"})
}

// AddDataToES
// @Summary Add document to search index.
// @Description Receive data from UI to upload to the search index. Auto inserts random ID and upload date values.
// @Accept json
// @Param body body search.Article true "Document that needs to be uploaded to the search index."
// @Produce json
// @Tags dashboard
// @Success 200 {object} search.Article
// @Router /fabric/dashboard/smart-contracts [post]
func AddDataToES(c *gin.Context) {

	var searchArticle search.Article
	if err := c.ShouldBindJSON(&searchArticle); err != nil {
		log.Println("add data err")
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	searchArticle.UploadDate = fmt.Sprintf(time.Now().UTC().Format("2006-01-02"))
	log.Println(searchArticle.UploadDate)
	res, err := search.AddDocumentToES(&searchArticle)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": res})
}

func QueryAssets(c *gin.Context) {

	cmd := exec.Command("bash", "queryAsset.sh")
	cmd.Dir = scriptPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(output)
		log.Println(err)
		errMessage := fmt.Sprintf(fmt.Sprint(err) + ": " + string(output))
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": errMessage})
		return
	}
	c.IndentedJSON(http.StatusOK, string(output))
}

func AssetTransfer(c *gin.Context) {
	var transactionRequest admin.TransactionRequest
	if err := c.BindJSON(&transactionRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message1": err})
		return
	}
	log.Println(transactionRequest)

	function := transactionRequest.Function

	switch function {
	case "CreateAsset":
		admin.CreateAsset(admin.ContractPass, transactionRequest)
		log.Println("CreateAsset")
	case "UpdateAsset":
		admin.UpdateAsset(admin.ContractPass, transactionRequest)
		log.Println("UpdateAsset")
	case "TransferAsset":
		admin.TransferAsset(admin.ContractPass, transactionRequest)
		log.Println("TransferAsset")
	default:
		log.Println(function)
		log.Fatalln("function selection error")
	}
}

func createRandomSHA() string {
	data := make([]byte, 10)
	var sha string
	if _, err := rand.Read(data); err == nil {
		sha = fmt.Sprintf("%x", sha256.Sum256(data))
	}
	return sha
}

func setEnv() {
	os.Setenv("CORE_PEER_TLS_ENABLED", "true")
	os.Setenv("CORE_PEER_LOCALMSPID", "Org1MSP")
	os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE",
		fmt.Sprintf("%s/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt", networkPath))
	os.Setenv("CORE_PEER_ADDRESS", "localhost:7051")
	os.Setenv("CORE_PEER_MSPCONFIGPATH", fmt.Sprintf("%s/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp", networkPath))
}
