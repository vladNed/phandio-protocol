#![no_std]

multiversx_sc::imports!();

pub mod state;

#[multiversx_sc::contract]
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
    ) {
        self.timeout_duration_1().set(timeout_duration_1);
        self.timeout_duration_2().set(timeout_duration_2);
        self.claim_commitment().set(claim_commitment);
        self.refund_commitment().set(refund_commitment);
        self.claimer().set(claimer);
        self.owner().set(owner);
        self.state().set(state::SwapState::Created);
    }

    /// The owner should `setReady` when
    /// - the claimer has locked the swap funds and the owner verified this
    /// - the duration of timeout_duration_1 has not passed
    #[endpoint(setReady)]
    fn set_ready(&self) {
        let caller = self.blockchain().get_caller();
        require!(self.owner().get() == caller, "only the owner can set the swap to ready");
        require!(self.state().get() == state::SwapState::Created, "swap is not in the created state");
        require!(self.blockchain().get_block_timestamp() < self.timeout_duration_1().get(), "timeout_duration_1 has passed");
        self.state().set(state::SwapState::Ready);
    }

    #[endpoint(claim)]
    fn claim(&self, claimer: ManagedAddress, claim_view_key: ManagedBuffer) {
        require!(self.state().get() == state::SwapState::Ready, "swap is not in the ready state");
        require!(self.claimer().get() == claimer, "claimer is not the same as the one set in the swap");
        todo!("verify the claim commitment");
        self.state().set(state::SwapState::Claimed);
    }

    #[storage_mapper("state")]
    fn state(&self) -> SingleValueMapper<state::SwapState>;

    #[storage_mapper("timeout_duration_1")]
    fn timeout_duration_1(&self) -> SingleValueMapper<u64>;

    #[storage_mapper("timeout_duration_2")]
    fn timeout_duration_2(&self) -> SingleValueMapper<u64>;

    #[storage_mapper("claim_commitment")]
    fn claim_commitment(&self) -> SingleValueMapper<ManagedBuffer>;

    #[storage_mapper("refund_commitment")]
    fn refund_commitment(&self) -> SingleValueMapper<ManagedBuffer>;

    #[storage_mapper("claimer")]
    fn claimer(&self) -> SingleValueMapper<ManagedAddress>;

    #[storage_mapper("owner")]
    fn owner(&self) -> SingleValueMapper<ManagedAddress>;
}
