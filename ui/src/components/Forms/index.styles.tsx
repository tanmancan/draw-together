import { Flex } from "@radix-ui/themes";
import { styled } from "@stitches/react";

export const StyledFormWrapper = styled("div", {
  width: `100vw`,
  height: `100vh`,
  display: `flex`,
  justifyContent: `center`,
  alignItems: `center`,
  flexFlow: `column`,
});

export const StyledFieldWrapper = styled(Flex, {
  width: `100%`,
  margin: 24,
  gap: 6,
});
