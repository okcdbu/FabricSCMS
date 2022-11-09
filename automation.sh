#!/bin/bash

DELAY=${10:-"3"}
IPADDR='127.0.0.1'

admin1(){
  ADMIN1=$(curl --location --request POST "$IPADDR:8080/fabric/lifecycle/admin/Org1")
  echo $ADMIN1
}

admin2(){
  ADMIN2=$(curl --location --request POST "$IPADDR:8080/fabric/lifecycle/admin/Org2")
  echo $ADMIN2
}

commitCC(){
  COMMITCC=$(curl --location --request POST "$IPADDR:8080/fabric/lifecycle/commit" \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "cc_name": "basic",
      "cc_sequence": 1,
      "cc_version": "1.0",
      "channel_name": "mychannel"
  }')
  echo $COMMITCC
}

queryInstall(){
  INSTALLID=$(curl --location --request GET "$IPADDR:8080/fabric/lifecycle/install")
  echo $INSTALLID
}

installCC(){
  INSTALLCC=$(curl --location --request POST "$IPADDR:8080/fabric/lifecycle/install/basic.tar.gz")
  echo $INSTALLCC
}

packageCC(){
  PACKAGE=$(curl --location --request POST "$IPADDR:8080/fabric/lifecycle/package" \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "cc_source_name": "asset-transfer-basic",
      "label": "basic_1.0",
      "language": "go",
      "package_name": "basic.tar.gz"
  }')
  echo $PACKAGE
}

approveOrg(){
    APPROVEORG=$(curl --location --request POST "$IPADDR:8080/fabric/lifecycle/approve" \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "cc_name": "basic",
        "cc_sequence": 1,
        "cc_version": "1.0",
        "channel_name": "mychannel",
        "package_ID": "basic_1.0:ff2e5d5f8ff054f8e4f951c592068e863c80b62cb6a9a355d4b51de886273c96"
    }')
    echo $APPROVEORG
}

queryCommitReady(){
     sleep $DELAY
    QUERYCOMMITREADY=$(curl --location --request GET "$IPADDR:8080/fabric/lifecycle/commit/organizations" \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "cc_name": "basic",
    "cc_sequence": 1,
    "cc_version": "1.0",
    "channel_name": "mychannel"
    }')
    echo $QUERYCOMMITREADY
}

queryCommitted(){
     sleep $DELAY
    QUERYCOMMIT=$(curl --location --request GET "$IPADDR:8080/fabric/lifecycle/commit" \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "cc_name": "basic",
    "channel_name": "mychannel"
    }')
    echo $QUERYCOMMIT
}

## Package the CC "asset-transfer-basic"
packageCC

## Set Org1 as the admin
admin1
## Install for Org1
installCC

## Set Org2 as the admin
admin2
## Install for Org2
installCC

## Set Org1 as the admin
admin1
## Check if installed succesfully
queryInstall

## Approve for Org1
admin1
approveOrg

## Check commit readiness from both Orgs
admin1
queryCommitReady

admin2
queryCommitReady

## Approve for Org2
admin2
approveOrg

## Check commit readiness from both Orgs
admin1
queryCommitReady

admin2
queryCommitReady
## Commit the definition
commitCC

## Check committed CC from both Orgs
admin1
queryCommitted

## Check committed CC from both Orgs
admin2
queryCommitted
