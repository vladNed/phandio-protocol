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

    #[upgrade]
    fn upgrade(&self, swap_template_address: ManagedAddress) {
        self.swap_template_address().set(swap_template_address);
    }

    #[endpoint(createSwap)]
    #[payable("EGLD")]
    fn create_swap(
        &self,
        timeout_duration_1: u64,
        timeout_duration_2: u64,
        claim_commitment: ManagedBuffer,
        refund_commitment: ManagedBuffer,
        claimer: ManagedAddress,
    ) -> ManagedAddress {
        let caller = self.blockchain().get_caller();
        let locked_amount = self.call_value().egld_value().clone_value();
        let swap_template_address = self.swap_template_address().get();
        // TODO: Hash the claim_commitment and refund_commitment and use them as the swap id.
        let new_swap_add = self.deploy_new_swap(
            timeout_duration_1,
            timeout_duration_2,
            claim_commitment,
            refund_commitment,
            claimer,
            caller,
            &locked_amount,
            swap_template_address,
        );

        new_swap_add
    }

    fn deploy_new_swap(
        &self,
        timeout_duration_1: u64,
        timeout_duration_2: u64,
        claim_commitment: ManagedBuffer,
        refund_commitment: ManagedBuffer,
        claimer: ManagedAddress,
        owner: ManagedAddress,
        locked_amount: &BigUint,
        from_source: ManagedAddress,
    ) -> ManagedAddress {
        let mut args = ManagedArgBuffer::new();
        args.push_arg(timeout_duration_1);
        args.push_arg(timeout_duration_2);
        args.push_arg(claim_commitment);
        args.push_arg(refund_commitment);
        args.push_arg(claimer);
        args.push_arg(owner);

        let gas_left = self.blockchain().get_gas_left();
        let (addr, _) = self.send_raw().deploy_from_source_contract(
            gas_left,
            &locked_amount,
            &from_source,
            CodeMetadata::READABLE | CodeMetadata::PAYABLE,
            &args,
        );

        addr
    }

    /// Address of the swap template contract.
    #[storage_mapper("swap_template_address")]
    fn swap_template_address(&self) -> SingleValueMapper<ManagedAddress>;
}
