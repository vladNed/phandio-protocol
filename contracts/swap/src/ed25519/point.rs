use core::ops::Mul;

pub struct FieldElement(pub [u64; 5]);

pub struct Scalar {
    pub bytes: [u8; 32],
}

pub struct EdwardsPoint {
    pub x: FieldElement,
    pub y: FieldElement,
    pub z: FieldElement,
    pub t: FieldElement,
}

impl<'b> Mul<&'b Scalar> for EdwardsPoint {
    type Output = EdwardsPoint;

    fn mul(self, rhs: &'b Scalar) -> Self::Output {
        todo!()
    }

}