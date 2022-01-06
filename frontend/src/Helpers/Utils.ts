export function Average(arr: number[]): number {
  return arr.length > 0 ? arr.reduce((a, b) => a + b) / arr.length : 0;
}
export function Max(arr: number[]): number {
  return arr.reduce((a, b) => Math.max(a, b));
}
export function Min(arr: number[]): number {
  return arr.reduce((a, b) => Math.min(a, b));
}
