multiversx_sc::imports!();

use swap_proxy::ProxyTrait as _;

pub mod swap_proxy {
    multiversx_sc::imports!();

    #[multiversx_sc::proxy]
    pub trait Swap {
        #[init]
        fn init(
            &self,
            timeout_duration_1: u64,
            timeout_duration_2: u64,
            claim_commitment: ManagedBuffer,
            refund_commitment: ManagedBuffer,
            claimer: ManagedAddress,
            owner: ManagedAddress,
            amount: BigUint,
        );
    }
}

#[multiversx_sc::module]
pub trait SwapProxyModule {
    fn deploy_new_swap(
        &self,
        timeout_duration_1: u64,
        timeout_duration_2: u64,
        claim_commitment: ManagedBuffer,
        refund_commitment: ManagedBuffer,
        claimer: ManagedAddress,
        owner: ManagedAddress,
        locked_amount: BigUint,
        from_source: ManagedAddress,
    ) -> (ManagedAddress, ()) {
        self.swap_proxy()
            .init(
                timeout_duration_1,
                timeout_duration_2,
                claim_commitment,
                refund_commitment,
                claimer,
                owner,
                locked_amount.clone(),
            )
            .with_egld_transfer(locked_amount)
            .deploy_from_source(&from_source, CodeMetadata::DEFAULT)
    }

    #[proxy]
    fn swap_proxy(&self) -> swap_proxy::Proxy<Self::Api>;
}
