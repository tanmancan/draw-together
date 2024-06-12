export const isError = (e: Error | unknown): e is Error => {
  if (e instanceof Error) return true;
  return false;
};
