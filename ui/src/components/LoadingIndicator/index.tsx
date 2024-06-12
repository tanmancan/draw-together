import { Flex } from "@radix-ui/themes";
import { keyframes, styled } from "@stitches/react";

const keySpin = keyframes({
  "0%": {
    transform: `rotate(0deg)`,
  },
  "100%": {
    transform: `rotate(360deg)`,
  },
});

const StyledSpinner = styled("div", {
  animation: `${keySpin} 1s infinite`,
  borderStyle: `solid`,
  borderColor: `var(--accent-9) var(--accent-9) transparent`,
  borderWidth: 10,
  borderRadius: `50%`,
  width: 50,
  height: 50,
});

function LoadingIndicator() {
  return (
    <Flex
      style={{
        width: `100vw`,
        height: `100vh`,
        overflow: "hidden",
      }}
      justify="center"
      align="center"
    >
      <StyledSpinner />
    </Flex>
  );
}

export default LoadingIndicator;
