'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const { initiateOvertakeProposal, checkProposalStatus } = require('./overtakeOperations');

class OvertakeWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        this.sutAdapter = sutAdapter;
        this.contractID = 'overtakingcc';  // Chaincode name
        this.channelID = 'mychannel';      // Channel name
    }

    async submitTransaction() {
        const args = ['10.0', '5.0', '20.0', '30.0', '2.0', '1.5'];
        await this.sutAdapter.sendRequests({
            contractId: this.contractID,
            contractFunction: 'InitiateOvertakeProposal',
            contractArguments: args,
            readOnly: false
        });
    }

    async cleanupWorkloadModule() {
        console.log('Cleanup workload module');
    }
}

module.exports.createWorkloadModule = () => {
    return new OvertakeWorkload();
};
