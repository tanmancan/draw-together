import * as Form from "@radix-ui/react-form";
import { Button, Text, TextField } from "@radix-ui/themes";
import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { handleCreateBoard } from "../../../lib/boards/handlers";
import { RoutesEnum } from "../../../lib/router/routes";
import { StyledFieldWrapper } from "../index.styles";
import { isError } from "../../../lib/helpers/errors";

function CreateBoardForm() {
  const [submitError, setSubmitError] = useState("");
  const navigate = useNavigate();
  const handleCreateBoardSubmit = async (e: FormEvent<HTMLFormElement>) => {
    try {
      e.preventDefault();
      const data = Object.fromEntries(new FormData(e.currentTarget));
      const { boardName } = data;

      const newBoard = await handleCreateBoard(boardName.toString());
      if (!newBoard) {
        setSubmitError("Error creating board.");
      } else {
        navigate(`${RoutesEnum.PLAY_BOARD}/${newBoard.id?.value}`);
      }
    } catch (e) {
      if (isError(e)) {
        setSubmitError(e.toString());
      }
    }
  };

  return (
    <Form.Root className="FormRoot" onSubmit={handleCreateBoardSubmit}>
      <StyledFieldWrapper>
        <Form.Field className="FormField" name="boardName">
          <TextField.Root>
            <Form.Control asChild>
              <TextField.Input
                className="Input"
                type="text"
                placeholder="Enter board name"
                required
                max={20}
              />
            </Form.Control>
          </TextField.Root>

          <Form.Message className="FormMessage" match="valueMissing">
            <Text>Please enter a name</Text>
          </Form.Message>
          <Form.Message className="FormMessage" match="typeMismatch">
            <Text>Please enter a name</Text>
          </Form.Message>
          {submitError !== "" && <Form.Message>{submitError}</Form.Message>}
        </Form.Field>
        <Form.Submit asChild>
          <Button className="Button">Create New Board</Button>
        </Form.Submit>
      </StyledFieldWrapper>
    </Form.Root>
  );
}

export default CreateBoardForm;
