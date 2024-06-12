export const asyncSleep = async (delay: number) => {
  return new Promise((r) => window.setTimeout(r, delay));
};
