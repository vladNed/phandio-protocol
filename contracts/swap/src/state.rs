multiversx_sc::imports!();
multiversx_sc::derive_imports!();

#[derive(NestedEncode, NestedDecode, TopEncode, TopDecode, PartialEq, Debug, TypeAbi)]
pub enum SwapState {
    Created,  // 0
    Ready,    // 1
    Claimed,  // 2
    Refunded, // 3
}

#[derive(TopDecode, TopEncode, Debug, TypeAbi, PartialEq)]
pub struct SwapStateData<M: ManagedTypeApi> {
    pub state: SwapState,
    pub timeout_duration_1: u64,
    pub timeout_duration_2: u64,
    pub claim_commitment: ManagedBuffer<M>,
    pub refund_commitment: ManagedBuffer<M>,
    pub claimer: ManagedAddress<M>,
    pub owner: ManagedAddress<M>,
    pub amount: BigUint<M>,
}

impl<M: ManagedTypeApi> SwapStateData<M> {
    pub fn new(
        timeout_duration_1: u64,
        timeout_duration_2: u64,
        claim_commitment: ManagedBuffer<M>,
        refund_commitment: ManagedBuffer<M>,
        claimer: ManagedAddress<M>,
        owner: ManagedAddress<M>,
        amount: BigUint<M>,
    ) -> Self {
        Self {
            state: SwapState::Created,
            timeout_duration_1,
            timeout_duration_2,
            claim_commitment,
            refund_commitment,
            claimer,
            owner,
            amount,
        }
    }
}
