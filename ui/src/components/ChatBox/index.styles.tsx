import { Box, Callout, Text } from "@radix-ui/themes";
import { styled } from "@stitches/react";

export const StyledUserName = styled(Text, {
  display: "inline-block",
  maxWidth: 80,
  overflow: "hidden",
  textOverflow: "ellipsis",
  whiteSpace: "nowrap",
});

export const StyledScrollBox = styled(Box, {
  height: `40vh`,
  overflowY: "scroll",
  margin: `calc(var(--card-padding) * -1)`,
  padding: `var(--card-padding)`,
  scrollbarWidth: "thin",
});

export const StyledMessageBody = styled(Callout.Text, {
  wordBreak: "break-all",
});
