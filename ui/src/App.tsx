import { Outlet } from "react-router-dom";
import LayoutDefault from "./layouts/LayoutDefault";
import { globalCss } from "@stitches/react";

const globalStyles = globalCss({
  body: {
    margin: 0,
  },
});

function App() {
  globalStyles();

  return (
    <LayoutDefault>
      <Outlet />
    </LayoutDefault>
  );
}

export default App;
