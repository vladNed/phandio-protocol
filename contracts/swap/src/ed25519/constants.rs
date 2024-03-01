use super::point::{EdwardsPoint, FieldElement};

pub const ED25519_BASE_POINT: EdwardsPoint = EdwardsPoint {
  x: FieldElement([
      1738742601995546,
      1146398526822698,
      2070867633025821,
      562264141797630,
      587772402128613,
  ]),
  y: FieldElement([
      1801439850948184,
      1351079888211148,
      450359962737049,
      900719925474099,
      1801439850948198,
  ]),
  z: FieldElement([1, 0, 0, 0, 0]),
  t: FieldElement([
      1841354044333475,
      16398895984059,
      755974180946558,
      900171276175154,
      1821297809914039,
  ]),
};