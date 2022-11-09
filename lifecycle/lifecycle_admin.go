package lifecycle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

type currentAdmin struct {
	Admin string `json:"admin"`
}

// SetAdmin
// @Summary Set an org as the admin.
// @Description Use terminal environmental variables to set the admin for peer cli container. Only Org1 and Org2 are supported.
// @Accept json
// @Param organization path string true "organization to be set as admin (Org1 and Org2 supported)"
// @Produce json
// @Tags lifecycle
// @Success 200 {object} currentAdmin
// @Router /fabric/lifecycle/admin/{organization} [post]
func SetAdmin(c *gin.Context) {
	var admin currentAdmin
	GOPATH := os.Getenv("GOPATH")
	networkPath := fmt.Sprintf("%s/src/github.com/hyperledger/fabric-samples/test-network/", GOPATH)
	os.Setenv("CORE_PEER_TLS_ENABLED", "true")

	organization := c.Param("organization")

	if !(organization == "Org1" || organization == "Org2") {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Only Org1 or Org2 supported."})
		return
	}

	os.Setenv("CORE_PEER_ADMIN", organization)
	os.Setenv("CORE_PEER_LOCALMSPID", fmt.Sprintf("%sMSP", organization))
	os.Setenv("CORE_PEER_TLS_ROOTCERT_FILE",
		fmt.Sprintf("%s/organizations/peerOrganizations/%s.example.com/peers/peer0.%s.example.com/tls/ca.crt",
			networkPath, strings.ToLower(organization), strings.ToLower(organization)))
	os.Setenv("CORE_PEER_MSPCONFIGPATH",
		fmt.Sprintf("%s/organizations/peerOrganizations/%s.example.com/users/Admin@%s.example.com/msp",
			networkPath, strings.ToLower(organization), strings.ToLower(organization)))

	if organization == "Org1" {
		os.Setenv("CORE_PEER_ADDRESS", "localhost:7051")
	} else if organization == "Org2" {
		os.Setenv("CORE_PEER_ADDRESS", "localhost:9051")
	}

	admin.Admin = organization
	c.IndentedJSON(http.StatusOK, admin)
}

// GetAdmin
// @Summary Get the current admin org.
// @Description Use terminal environmental variables to get the admin for peer cli container. Only Org1 and Org2 are supported.
// @Accept json
// @Produce json
// @Tags lifecycle
// @Success 200 {object} currentAdmin
// @Router /fabric/lifecycle/admin [get]
func GetAdmin(c *gin.Context) {
	var admin currentAdmin
	os.Setenv("CORE_PEER_TLS_ENABLED", "true")

	envAdmin := os.Getenv("CORE_PEER_ADMIN")

	if envAdmin == "Org1" || envAdmin == "Org2" {
		admin.Admin = envAdmin
		c.IndentedJSON(http.StatusOK, admin)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error getting current admin. Please check if an admin has been set."})
}
