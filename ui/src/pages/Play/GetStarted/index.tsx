import {
  Box,
  Button,
  Card,
  Heading,
  Separator,
  TextField,
} from "@radix-ui/themes";
import CreateBoardForm from "../../../components/Forms/CreateBoardForm";
import { ChangeEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { RoutesEnum } from "../../../lib/router/routes";
import {
  StyledFieldWrapper,
  StyledFormWrapper,
} from "../../../components/Forms/index.styles";

function PageGetStarted() {
  const navigate = useNavigate();
  const [gotoID, setGotoID] = useState("");
  const gotoIDHandler = () => {
    if (!gotoID) return;
    navigate(`${RoutesEnum.PLAY_BOARD}/${gotoID}`);
  };
  return (
    <StyledFormWrapper>
      <Box mb={"4"}>
        <Heading>Get Started</Heading>
      </Box>
      <Card>
        <Box>
          <StyledFieldWrapper>
            <TextField.Root>
              <TextField.Input
                name="boardId"
                id="boardId"
                placeholder="Enter board ID"
                value={gotoID}
                onChange={(e: ChangeEvent<HTMLInputElement>) =>
                  setGotoID(e.target.value)
                }
                required
                max={20}
              />
            </TextField.Root>
            <Button onClick={gotoIDHandler}>Join Existing Board</Button>
          </StyledFieldWrapper>
        </Box>
        <Separator size={"4"} />
        <Box>
          <CreateBoardForm />
        </Box>
      </Card>
    </StyledFormWrapper>
  );
}

export default PageGetStarted;
