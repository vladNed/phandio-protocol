// Code generated by the multiversx-sc build system. DO NOT EDIT.

////////////////////////////////////////////////////
////////////////// AUTO-GENERATED //////////////////
////////////////////////////////////////////////////

// Init:                                 1
// Endpoints:                            5
// Async Callback (empty):               1
// Total number of exported functions:   7

#![no_std]
#![allow(internal_features)]
#![feature(lang_items)]

multiversx_sc_wasm_adapter::allocator!();
multiversx_sc_wasm_adapter::panic_handler!();

multiversx_sc_wasm_adapter::endpoints! {
    swap
    (
        init => init
        setReady => set_ready
        claim => claim
        refund => refund
        getState => state
        getSecretCommitment => secret_commitment
    )
}

multiversx_sc_wasm_adapter::async_callback_empty! {}
