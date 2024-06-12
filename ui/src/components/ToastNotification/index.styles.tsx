import * as Toast from "@radix-ui/react-toast";
import { Callout } from "@radix-ui/themes";
import { styled } from "@stitches/react";

export const StyledToastViewport = styled(Toast.Viewport, {
  position: "fixed",
  width: "fit-content",
  height: "fit-content",
  maxWidth: 520,
  bottom: 0,
  right: 16,
});

export const StyledToastWrapper = styled("div", {
  borderRadius: `var(--radius-4)`,
  backgroundColor: `var(--color-panel)`,
});

export const StyledCalloutRoot = styled(Callout.Root, {
  transform: `translate3d(var(--radix-toast-swipe-move-x), 0, 0)`,
});

export const StyledCalloutIcon = styled(Callout.Icon, {
  marginRight: 8,
});

export const StyledCalloutText = styled(Callout.Text, {});

export const StyledToastClose = styled(Toast.Close, {
  gridRowStart: 1,
  gridColumnStart: 2,
  marginLeft: 8,
});
