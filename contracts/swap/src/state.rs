multiversx_sc::imports!();
multiversx_sc::derive_imports!();

#[derive(NestedEncode, NestedDecode, TopEncode, TopDecode, PartialEq, Debug, TypeAbi)]
pub enum SwapState {
    Created,  // 0
    Ready,    // 1
    Claimed,  // 2
    Refunded, // 3
}
