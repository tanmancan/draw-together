import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import basicSsl from "@vitejs/plugin-basic-ssl";

const CSP_NONCE =
  "7079816de27d50c2535a1266c6f9256e6473ceb182fa418fc15b78f882afdc7633b4bdffa617e38f329a8afbe2a0faf3361b783c637740f640ab0c9b6e960126";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    basicSsl(),
    {
      name: "csp-nonce-script",
      enforce: "post",
      transformIndexHtml(html: string) {
        const regex = /<script(.*?)/gi;
        const replacement = `<script nonce="${CSP_NONCE}"$1`;
        return html.replace(regex, replacement);
      },
    },
  ],
  server: {
    headers: {
      "Content-Security-Policy": `default-src 'self';script-src 'self' 'nonce-${CSP_NONCE}';style-src 'self' 'unsafe-inline';child-src 'none';img-src 'self' data:;font-src 'self' fonts.gstatic.com;connect-src 'self' https://localhost:8443 wss://localhost:8443;object-src 'none';frame-ancestors 'none';frame-src 'none';worker-src 'none';`,
    },
  },
});
