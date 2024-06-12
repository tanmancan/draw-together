import { Button, Heading } from "@radix-ui/themes";
import { styled } from "@stitches/react";
import { Link } from "react-router-dom";
import { RoutesEnum } from "../../../lib/router/routes";

const StyledWrapper = styled("div", {
  width: `100vw`,
  height: `100vh`,
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  flexFlow: "column",
  gap: 16,
});
function PageNotFound() {
  return (
    <StyledWrapper>
      <Heading>404: Not Found</Heading>
      <Button asChild>
        <Link to={RoutesEnum.PLAY_GET_STARTED}>Get Started</Link>
      </Button>
    </StyledWrapper>
  );
}

export default PageNotFound;
