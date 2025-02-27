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

// Required dependencies for Caliper
const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class EvaluateSafetyWorkload extends WorkloadModuleBase {

    // Initialize the workload (optional) before starting the benchmark
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        this.workerIndex = workerIndex;
        this.sutAdapter = sutAdapter;
        this.sutContext = sutContext;
    }

    // Submit the transaction to the smart contract
    async submitTransaction() {
        // Randomized parameters for the EvaluateSafety function
        const relativeSpeed = (Math.random() * 10 + 5).toFixed(2);   // Vr (5-15)
        const oncomingSpeed = (Math.random() * 20 + 20).toFixed(2);  // Vo (20-40)
        const visibilityDistance = (Math.random() * 500 + 200).toFixed(2); // Dv (200-700)
        const overtakingDistance = (Math.random() * 50 + 30).toFixed(2);   // Do (30-80)
        const reactionTime = (Math.random() * 1.5 + 0.5).toFixed(2); // Tr (0.5-2)
        const safetyMargin = (Math.random() * 10 + 5).toFixed(2);    // Sm (5-15)

        // Construct the transaction arguments
        const args = {
            relativeSpeed: parseFloat(relativeSpeed),
            oncomingSpeed: parseFloat(oncomingSpeed),
            visibilityDistance: parseFloat(visibilityDistance),
            overtakingDistance: parseFloat(overtakingDistance),
            reactionTime: parseFloat(reactionTime),
            safetyMargin: parseFloat(safetyMargin),
        };

        // Submit the transaction invoking the EvaluateSafety function
        try {
            await this.sutAdapter.sendRequests({
                contractId: 'overtake_chaincode',    // Name of your smart contract
                contractFunction: 'EvaluateSafety',  // Function to invoke
                contractArguments: [JSON.stringify(args)],  // Pass the parameters as a JSON string
                readOnly: false  // This is an invoke operation, not a query
            });
        } catch (error) {
            console.error(`Failed to submit EvaluateSafety transaction: ${error}`);
        }
    }

    // Cleanup function after workload execution
    async cleanupWorkloadModule() {
        // Perform any necessary cleanup tasks (if needed)
    }
}

module.exports.createWorkloadModule = () => {
    return new EvaluateSafetyWorkload();
};
