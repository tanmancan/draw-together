import { TextArea, TextField, Button, Flex } from "@radix-ui/themes";
import { styled } from "@stitches/react";
import {
  BASE_VIEWPORT_WIDTH,
  BASE_VIEWPORT_HEIGHT,
} from "../../../lib/viewport";

export const SIDEBAR_WIDTH = 300;

export const StyledWrapper = styled("div", {
  display: "grid",
  gridTemplateColumns: `1fr 1fr 1fr ${SIDEBAR_WIDTH}px`,
  gridTemplateRows: "auto",
  gridTemplateAreas: `"viewport viewport viewport sidebar"`,
  width: `100vw`,
  height: `100vh`,
  overflow: "hidden",
  "@media (max-width: 768px)": {
    gridTemplateColumns: `1fr 1fr 1fr 0px`,
  },
});

export const StyledViewPortArea = styled("div", {
  gridArea: "viewport",
  boxSizing: "border-box",
  position: "relative",
  height: "100vh",
  background: "var(--gray-2)",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  overflow: "scroll",
  scrollbarWidth: "thin",
  flexFlow: "column",
});

export const StyledSideBarArea = styled("div", {
  gridArea: "sidebar",
  boxSizing: "border-box",
  background: "var(--gray-3)",
  display: "flex",
  height: "100vh",
  flexFlow: "column",
  gap: 16,
  padding: 16,
  borderLeft: `1px solid var(--gray-5)`,
  overflowY: "scroll",
  "@media (max-width: 768px)": {
    gridArea: "viewport",
    zIndex: -1,
    "& .mobile-show": {
      zIndex: 1,
    },
  },
});

export const StyledViewPort = styled("canvas", {
  background: "white",
  border: `1px solid #ccc`,
  borderRadius: 5,
  boxSizing: "border-box",
  position: "relative",
  width: BASE_VIEWPORT_WIDTH,
  height: BASE_VIEWPORT_HEIGHT,
  minWidth: BASE_VIEWPORT_WIDTH,
  minHeight: BASE_VIEWPORT_HEIGHT,
  zIndex: 0,
  margin: "1rem",
});

export const StyledTextArea = styled(TextArea, {});

export const StyledInputWrapper = styled("div", {
  display: "flex",
  flexFlow: "row",
  margin: `8 0 0`,
});

export const StyledTextFieldRoot = styled(TextField.Root, {
  width: `100%`,
});

export const StyledTextFieldInput = styled(TextField.Input, {});

export const StyledSubmitButton = styled(Button, {});

export const StyledCaptionWrapper = styled(Flex, {
  width: BASE_VIEWPORT_WIDTH,
  marginBottom: "1rem",
});
