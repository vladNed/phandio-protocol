multiversx_sc::imports!();
multiversx_sc::derive_imports!();

#[multiversx_sc::module]
pub trait EventsModule {
    #[event("claimed")]
    fn claimed_event(&self, #[indexed] claimer: &ManagedAddress, #[indexed] timestamp: u64);

    #[event("ready")]
    fn ready_event(&self, #[indexed] timestamp: u64);

    #[event("refunded")]
    fn refund_event(&self, #[indexed] timestamp: u64);
}
