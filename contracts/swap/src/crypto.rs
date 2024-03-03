use core::marker::PhantomData;

use curve25519_dalek::{constants::ED25519_BASEPOINT_POINT, scalar::Scalar};
use multiversx_sc::storage::mappers::SingleValueMapper;
use multiversx_sc::{api::ManagedTypeApi, require, types::ManagedBuffer, module};

/// Wrapper over an ed25519 public key which is an edwards point on the
/// curve25519 translated as a managed buffer.
pub(crate) struct ManagedPublicKey<M: ManagedTypeApi> {
    pub data: ManagedBuffer<M>,
}

impl<M: ManagedTypeApi> From<[u8; 32]> for ManagedPublicKey<M> {
    fn from(data: [u8; 32]) -> Self {
        let buf = ManagedBuffer::from(&data);
        Self { data: buf }
    }
}

/// Wrapper over an ed25519 scalar which is a 256-bit integer
pub(crate) struct ManagedPrivateKey<M: ManagedTypeApi> {
    data: Scalar,
    _phantom: PhantomData<M>,
}

impl<M: ManagedTypeApi> From<ManagedBuffer<M>> for ManagedPrivateKey<M> {
    fn from(value: ManagedBuffer<M>) -> Self {
        let mut data = [0u8; 32];
        value.load_to_byte_array(&mut data);
        let scalar = Scalar::from_bytes_mod_order(data);

        Self {
            data: scalar,
            _phantom: PhantomData,
        }
    }
}

impl<M: ManagedTypeApi> ManagedPrivateKey<M> {
    pub fn public_spend_keu(&self) -> ManagedPublicKey<M> {
        ManagedPublicKey::from((self.data * &ED25519_BASEPOINT_POINT).compress().to_bytes())
    }
}


/// Module used to verify ed25519 keys.
#[module]
pub trait CryptoModule {
    /// Verifies that the secret provided on claim is the correct private spend
    /// key for the public spend key provided.
    ///
    /// It derives the ed25519 public spend key from the provided secret and
    /// calculates the hash of the public spend key with keccak256.
    ///
    /// If the hash matches the commitment, the secret is correct.
    fn verify_commitment(&self, commitment: &ManagedBuffer, secret: &ManagedBuffer) {
        let private_spend_key = ManagedPrivateKey::from(secret.clone());
        let public_spend_key = private_spend_key.public_spend_keu();
        let secret_pk_hash = self.hash_key(&public_spend_key);
        require!(&secret_pk_hash == commitment, "Invalid secret provided");
    }

    /// Using keccak256, hashes the public spend key.
    fn hash_key(&self, key: &ManagedPublicKey<Self::Api>) -> ManagedBuffer {
        let key_hash = self.crypto().keccak256(key.data.clone());
        key_hash.as_managed_buffer().clone()
    }
}
