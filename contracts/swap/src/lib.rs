#![no_std]

multiversx_sc::imports!();

pub mod crypto;
pub mod events;
pub mod state;

const CONTRACT_CURRENT_VERSION : &str = "1.0.0";

#[multiversx_sc::contract]
pub trait Swap: events::EventsModule + crypto::CryptoModule {
    #[init]
    #[payable("EGLD")]
    fn init(
        &self,
        timeout_duration_1: u64,
        timeout_duration_2: u64,
        claim_commitment: ManagedBuffer,
        refund_commitment: ManagedBuffer,
        claimer: ManagedAddress,
        owner: ManagedAddress,
    ) {
        self.state().set(state::SwapStateData::new(
            timeout_duration_1,
            timeout_duration_2,
            claim_commitment,
            refund_commitment,
            claimer,
            owner,
            self.call_value().egld_value().clone_value(),
        ));
        self.version().set(ManagedBuffer::from(CONTRACT_CURRENT_VERSION));
    }

    #[upgrade]
    fn upgrade(&self) {}

    /// The owner should `setReady` when
    /// - the claimer has locked the swap funds and the owner verified this
    /// - the duration of timeout_duration_1 has not passed
    #[endpoint(setReady)]
    fn set_ready(&self) {
        let caller = self.blockchain().get_caller();
        let state_handler = self.state();
        let mut swap_state = state_handler.get();
        require!(
            swap_state.owner == caller,
            "only the owner can set the swap to ready"
        );
        require!(
            swap_state.state == state::SwapState::Created,
            "swap is not in the created state"
        );
        require!(
            self.blockchain().get_block_timestamp() < swap_state.timeout_duration_1,
            "timeout_duration_1 has passed"
        );

        swap_state.state = state::SwapState::Ready;
        state_handler.set(swap_state);
        self.ready_event(self.blockchain().get_block_timestamp());
    }

    /// The claimer should be able to claim when
    /// - the owner has set the swap to `Ready` state, and it's before timeout 1
    /// - the claim window is within timeout 1 and timeout 2
    /// - the provided claim view key is the right private key to the public key on the other chain
    #[endpoint(claim)]
    fn claim(&self, claim_view_key: ManagedBuffer) {
        let caller = self.blockchain().get_caller();
        let state_handler = self.state();
        let mut swap_state = state_handler.get();
        require!(
            caller == swap_state.claimer,
            "only claimer address can perform this"
        );
        require!(
            swap_state.state == state::SwapState::Ready,
            "cannot perform claim"
        );
        let block_timestamp = self.blockchain().get_block_timestamp();
        require!(
            block_timestamp < swap_state.timeout_duration_1
                && swap_state.state != state::SwapState::Ready,
            "to early to claim"
        );
        require!(
            block_timestamp >= swap_state.timeout_duration_2,
            "to late to claim"
        );
        let secret_commitment_handler = self.secret_commitment();
        require!(
            secret_commitment_handler.is_empty(),
            "secret already claimed"
        );
        self.verify_commitment(&swap_state.claim_commitment, &claim_view_key);
        secret_commitment_handler.set(claim_view_key);
        self.send().direct_egld(&caller, &swap_state.amount);
        swap_state.state = state::SwapState::Claimed;
        state_handler.set(swap_state);
        self.claimed_event(&caller, block_timestamp);
    }

    #[endpoint(refund)]
    fn refund(&self, refund_claim_key: ManagedBuffer) {
        let caller = self.blockchain().get_caller();
        let state_handler = self.state();
        let mut swap_state = state_handler.get();
        require!(
            caller == swap_state.owner,
            "only owner address can perform this"
        );
        require!(
            swap_state.state == state::SwapState::Ready,
            "cannot perform refund"
        );
        let block_timestamp = self.blockchain().get_block_timestamp();
        require!(
            block_timestamp < swap_state.timeout_duration_2
                && (block_timestamp > swap_state.timeout_duration_1
                    || swap_state.state == state::SwapState::Ready),
            "refund window is overdue"
        );

        self.verify_commitment(&swap_state.refund_commitment, &refund_claim_key);
        self.send().direct_egld(&caller, &swap_state.amount);
        swap_state.state = state::SwapState::Refunded;
        state_handler.set(swap_state);
        self.refund_event(block_timestamp);
    }

    /// Represents the details of the swap both parties agreed on.
    #[view(getState)]
    #[storage_mapper("state")]
    fn state(&self) -> SingleValueMapper<state::SwapStateData<Self::Api>>;

    /// The commitment of the secret provided on claim.
    /// This should be empty until a valid claim is completed.
    #[view(getSecretCommitment)]
    #[storage_mapper("secret_commitment")]
    fn secret_commitment(&self) -> SingleValueMapper<ManagedBuffer>;

    #[view(getVersion)]
    #[storage_mapper("version")]
    fn version(&self) -> SingleValueMapper<ManagedBuffer>;
}
