export function Average(arr: number[]): number {
  return arr.reduce((a, b) => a + b) / arr.length;
}
