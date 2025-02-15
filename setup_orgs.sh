#!/bin/bash

# Check if Docker service is running
if ! systemctl is-active --quiet docker; then
    echo "Docker service is not running. Starting it now..."
    sudo systemctl start docker
else
    echo "Docker service is already running."
fi

git checkout local

export PATH=${PWD}/bin:$PATH

echo -e "======Setting up Org1======\n"
(cd org1; ./1_enrollOrg1AdminAndUsers.sh; ./2_generateMSPOrg1.sh)
 
echo -e "\n======Setting up Org2======\n"
(cd org2; ./1_enrollOrg2AdminAndUsers.sh; ./2_generateMSPOrg2.sh )
 
echo -e "\n======Setting up Org3======\n"
(cd org3; ./1_enrollOrg3AdminAndUsers.sh; ./2_generateMSPOrg3.sh )

echo -e "\n======Setting up Orderer======\n"
(cd orderer; ./1_enrollAdminAndMSP.sh; ./2_artifact.sh)
 
# echo -e "\n======Creating Channel======\n"
# (cd org1; ./3_createChannel.sh)
 
# echo -e "\n======Org 2 Joining Channel======\n"

# (cd org2; ./3_joinChannel.sh)

# echo -e "\n======Org 3 Joining Channel======\n"

# (cd org3; ./3_joinChannel.sh)