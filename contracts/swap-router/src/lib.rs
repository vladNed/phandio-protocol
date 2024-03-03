#![no_std]

mod swap_proxy;

multiversx_sc::imports!();

/// Swap Router contract is the contract that deploys individual swap contracts
/// and provides a way to interact with them.
/// It is the main entry point for the swap contracts.
#[multiversx_sc::contract]
pub trait SwapRouter: swap_proxy::SwapProxyModule {
    #[init]
    fn init(&self, swap_template_address: ManagedAddress) {
        self.swap_template_address().set(swap_template_address);
    }

    #[payable("EGLD")]
    #[endpoint(create_swap)]
    fn create_swap(
        &self,
        timeout_duration_1: u64,
        timeout_duration_2: u64,
        claim_commitment: ManagedBuffer,
        refund_commitment: ManagedBuffer,
        claimer: ManagedAddress,
    ) -> ManagedAddress {
        let caller = self.blockchain().get_caller();
        let locked_amount = self.call_value().egld_or_single_esdt();
        let swap_template_address = self.swap_template_address().get();
        // TODO: Hash the claim_commitment and refund_commitment and use them as the swap id.
        let (new_swap_add, _) = self.deploy_new_swap(
            timeout_duration_1,
            timeout_duration_2,
            claim_commitment,
            refund_commitment,
            claimer,
            caller,
            locked_amount.amount,
            swap_template_address,
        );

        new_swap_add
    }

    /// Address of the swap template contract.
    #[storage_mapper("swap_template_address")]
    fn swap_template_address(&self) -> SingleValueMapper<ManagedAddress>;
}
