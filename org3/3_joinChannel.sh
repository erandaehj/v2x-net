#!/usr/bin/env bash

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/../orderer/crypto-config-ca/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_ORG3_CA=${PWD}/crypto-config-ca/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
export FABRIC_CFG_PATH=${PWD}/../config

export CHANNEL_NAME=mychannel

setGlobalsForPeer0Org3() {
  export CORE_PEER_LOCALMSPID="Org3MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
  export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config-ca/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
  export CORE_PEER_ADDRESS=localhost:9051

}

setGlobalsForPeer1Org3() {
  export CORE_PEER_LOCALMSPID="Org3MSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
  export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config-ca/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
  export CORE_PEER_ADDRESS=localhost:10051

}

fetchChannelBlock() {
  setGlobalsForPeer0Org3
  peer channel fetch 0 $CHANNEL_NAME.block -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
}
joinChannel() {
  setGlobalsForPeer0Org3
  peer channel join -b ./$CHANNEL_NAME.block

  setGlobalsForPeer1Org3
  peer channel join -b ./$CHANNEL_NAME.block
}

updateAnchorPeers() {

  setGlobalsForPeer0Org3
  peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL_NAME -f ./${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA

}
fetchChannelBlock
joinChannel
updateAnchorPeers
