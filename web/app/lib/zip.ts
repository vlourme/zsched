export const zip = (a: number[], ...rest: number[][]) =>
  a.map((value, index) => [value, ...rest.map((r) => r[index])]);
