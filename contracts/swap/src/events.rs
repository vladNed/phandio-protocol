multiversx_sc::imports!();
multiversx_sc::derive_imports!();



#[multiversx_sc::module]
pub trait EventsModule {
    
    fn generate_claimed_event(&self, claimer: &ManagedAddress, timestamp: u64) {
        self.claimed_event(claimer, timestamp);
    }
    
    #[event("claimed")]
    fn claimed_event(
        &self,
        #[indexed] claimer: &ManagedAddress,
        #[indexed] timestamp: u64
    );
}