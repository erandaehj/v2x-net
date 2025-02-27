/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

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
