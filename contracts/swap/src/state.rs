multiversx_sc::imports!();
multiversx_sc::derive_imports!();

#[derive(NestedEncode, NestedDecode, TopEncode, TopDecode, PartialEq, Debug)]
pub enum SwapState {
    Created,
    Ready,
    Claimed,
    Refunded,
}
