#![no_std]

multiversx_sc::imports!();


/// Swap Router contract is the contract that deploys individual swap contracts
/// and provides a way to interact with them.
/// It is the main entry point for the swap contracts.
#[multiversx_sc::contract]
pub trait SwapRouter {
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
    ) {
        todo!();
    }

    /// Address of the swap template contract.
    #[storage_mapper("swap_template_address")]
    fn swap_template_address(&self) -> SingleValueMapper<ManagedAddress>;
}
