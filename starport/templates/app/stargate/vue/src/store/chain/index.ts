import CosmosCosmosSdkIbcApplicationsTransferV1 from './cosmos/cosmos-sdk/ibc.applications.transfer.v1'
import CosmosCosmosSdkCosmosVestingV1Beta1 from './cosmos/cosmos-sdk/cosmos.vesting.v1beta1'
import CosmosCosmosSdkCosmosStakingV1Beta1 from './cosmos/cosmos-sdk/cosmos.staking.v1beta1'
import CosmosCosmosSdkIbcCoreChannelV1 from './cosmos/cosmos-sdk/ibc.core.channel.v1'
import CosmosCosmosSdkCosmosEvidenceV1Beta1 from './cosmos/cosmos-sdk/cosmos.evidence.v1beta1'
import CosmosCosmosSdkCosmosGovV1Beta1 from './cosmos/cosmos-sdk/cosmos.gov.v1beta1'
import CosmosCosmosSdkIbcCoreClientV1 from './cosmos/cosmos-sdk/ibc.core.client.v1'
import CosmosCosmosSdkIbcCoreConnectionV1 from './cosmos/cosmos-sdk/ibc.core.connection.v1'
import CosmosCosmosSdkCosmosCrisisV1Beta1 from './cosmos/cosmos-sdk/cosmos.crisis.v1beta1'
import CosmosCosmosSdkCosmosBankV1Beta1 from './cosmos/cosmos-sdk/cosmos.bank.v1beta1'
import CosmosCosmosSdkCosmosDistributionV1Beta1 from './cosmos/cosmos-sdk/cosmos.distribution.v1beta1'
import CosmosCosmosSdkCosmosSlashingV1Beta1 from './cosmos/cosmos-sdk/cosmos.slashing.v1beta1'


export default { 
  CosmosCosmosSdkIbcApplicationsTransferV1: load(CosmosCosmosSdkIbcApplicationsTransferV1, 'chain/cosmos/cosmos-sdk/ibc.applications.transfer.v1'),
  CosmosCosmosSdkCosmosVestingV1Beta1: load(CosmosCosmosSdkCosmosVestingV1Beta1, 'chain/cosmos/cosmos-sdk/cosmos.vesting.v1beta1'),
  CosmosCosmosSdkCosmosStakingV1Beta1: load(CosmosCosmosSdkCosmosStakingV1Beta1, 'chain/cosmos/cosmos-sdk/cosmos.staking.v1beta1'),
  CosmosCosmosSdkIbcCoreChannelV1: load(CosmosCosmosSdkIbcCoreChannelV1, 'chain/cosmos/cosmos-sdk/ibc.core.channel.v1'),
  CosmosCosmosSdkCosmosEvidenceV1Beta1: load(CosmosCosmosSdkCosmosEvidenceV1Beta1, 'chain/cosmos/cosmos-sdk/cosmos.evidence.v1beta1'),
  CosmosCosmosSdkCosmosGovV1Beta1: load(CosmosCosmosSdkCosmosGovV1Beta1, 'chain/cosmos/cosmos-sdk/cosmos.gov.v1beta1'),
  CosmosCosmosSdkIbcCoreClientV1: load(CosmosCosmosSdkIbcCoreClientV1, 'chain/cosmos/cosmos-sdk/ibc.core.client.v1'),
  CosmosCosmosSdkIbcCoreConnectionV1: load(CosmosCosmosSdkIbcCoreConnectionV1, 'chain/cosmos/cosmos-sdk/ibc.core.connection.v1'),
  CosmosCosmosSdkCosmosCrisisV1Beta1: load(CosmosCosmosSdkCosmosCrisisV1Beta1, 'chain/cosmos/cosmos-sdk/cosmos.crisis.v1beta1'),
  CosmosCosmosSdkCosmosBankV1Beta1: load(CosmosCosmosSdkCosmosBankV1Beta1, 'chain/cosmos/cosmos-sdk/cosmos.bank.v1beta1'),
  CosmosCosmosSdkCosmosDistributionV1Beta1: load(CosmosCosmosSdkCosmosDistributionV1Beta1, 'chain/cosmos/cosmos-sdk/cosmos.distribution.v1beta1'),
  CosmosCosmosSdkCosmosSlashingV1Beta1: load(CosmosCosmosSdkCosmosSlashingV1Beta1, 'chain/cosmos/cosmos-sdk/cosmos.slashing.v1beta1'),
  
}

function load(mod, fullns) {
    return function init(store) {
        const fullnsLevels = fullns.split('/')
        for (let i = 1; i < fullnsLevels.length; i++) {
            let ns = fullnsLevels.slice(0, i)
            if (!store.hasModule(ns)) {
                store.registerModule(ns, { namespaced: true })
            }
        }
        store.registerModule(fullnsLevels, mod)
        store.subscribe((mutation) => {
            if (mutation.type == 'chain/common/env/INITIALIZE_WS_COMPLETE') {
                store.dispatch(fullns+ '/init', null, {
                    root: true
                })
            }
        })
    }
}
