'use strict';

const { Contract } = require('fabric-network');

async function initiateOvertakeProposal(contract, args) {
    const [Dv, Do, Vr, Vo, Tr, Sm] = args;
    const result = await contract.submitTransaction('InitiateOvertakeProposal', Dv, Do, Vr, Vo, Tr, Sm);
    console.log(`Proposal initiated: ${result.toString()}`);
    return result.toString();
}

async function checkProposalStatus(contract, args) {
    const [proposalId] = args;
    const result = await contract.evaluateTransaction('CheckProposalStatus', proposalId);
    console.log(`Proposal status: ${result.toString()}`);
    return result.toString();
}

module.exports = {
    initiateOvertakeProposal,
    checkProposalStatus
};
