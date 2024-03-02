#![no_std]

multiversx_sc::imports!();

pub mod crypto;
pub mod events;
pub mod state;

#[multiversx_sc::contract]
pub trait Swap: events::EventsModule + crypto::CryptoModule {
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
    ) {
        self.timeout_duration_1().set(timeout_duration_1);
        self.timeout_duration_2().set(timeout_duration_2);
        self.claim_commitment().set(claim_commitment);
        self.refund_commitment().set(refund_commitment);
        self.claimer().set(claimer);
        self.owner().set(owner);
        self.state().set(state::SwapState::Created);
        self.amount().set(amount);
    }

    /// The owner should `setReady` when
    /// - the claimer has locked the swap funds and the owner verified this
    /// - the duration of timeout_duration_1 has not passed
    #[endpoint(setReady)]
    fn set_ready(&self) {
        let caller = self.blockchain().get_caller();
        require!(
            self.owner().get() == caller,
            "only the owner can set the swap to ready"
        );
        require!(
            self.state().get() == state::SwapState::Created,
            "swap is not in the created state"
        );
        require!(
            self.blockchain().get_block_timestamp() < self.timeout_duration_1().get(),
            "timeout_duration_1 has passed"
        );
        self.state().set(state::SwapState::Ready);
        self.ready_event(self.blockchain().get_block_timestamp());
    }

    /// The claimer should be able to claim when
    /// - the owner has set the swap to `Ready` state, and it's before timeout 1
    /// - the claim window is within timeout 1 and timeout 2
    /// - the provided claim view key is the right private key to the public key on the other chain
    #[endpoint(claim)]
    fn claim(&self, claim_view_key: ManagedBuffer) {
        let caller = self.blockchain().get_caller();
        require!(
            caller == self.claimer().get(),
            "only claimer address can perform this"
        );
        let swap_state = self.state().get();
        require!(
            swap_state == state::SwapState::Ready,
            "cannot perform claim"
        );
        let block_timestamp = self.blockchain().get_block_timestamp();
        require!(
            block_timestamp < self.timeout_duration_1().get()
                && swap_state != state::SwapState::Ready,
            "to early to claim"
        );
        require!(
            block_timestamp >= self.timeout_duration_2().get(),
            "to late to claim"
        );
        self.verify_commitment(&self.claim_commitment().get(), &claim_view_key);
        self.send().direct_egld(&caller, &self.amount().get());
        self.state().set(state::SwapState::Claimed);
        self.claimed_event(&caller, block_timestamp);
    }

    #[endpoint(refund)]
    fn refund(&self, refund_claim_key: ManagedBuffer) {
        let caller = self.blockchain().get_caller();
        require!(
            caller == self.owner().get(),
            "only owner address can perform this"
        );
        let swap_state = self.state().get();
        require!(
            swap_state == state::SwapState::Ready,
            "cannot perform refund"
        );
        let block_timestamp = self.blockchain().get_block_timestamp();
        require!(
            block_timestamp < self.timeout_duration_2().get()
                && (block_timestamp > self.timeout_duration_1().get()
                    || swap_state == state::SwapState::Ready),
            "refund window is overdue"
        );

        // TODO: Verify refund claim key
        self.send().direct_egld(&caller, &self.amount().get());
        self.refund_event(block_timestamp);
    }

    #[view(getState)]
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

    #[storage_mapper("amount")]
    fn amount(&self) -> SingleValueMapper<BigUint>;
}
