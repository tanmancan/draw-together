// import.meta.env.BASE_URL ??
export const BASE_PROTOCOL = window.location.protocol;

export const BASE_HOST = window.location.hostname;

export const BASE_PORT = window.location.port;

export const BASE_URL = `${location.protocol}://${location.hostname}${
  window.location.port ? `:${window.location.port}` : ""
}`;

export const getBaseUrl = () => {
  console.log(BASE_URL);
  return BASE_URL;
};
