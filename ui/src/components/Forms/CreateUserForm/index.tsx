import * as Form from "@radix-ui/react-form";
import { Button, TextField } from "@radix-ui/themes";
import { handleCreateUser } from "../../../lib/users/handlers";
import { FormEvent, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  isUserLoaded,
  useUserState,
  useUserStateAction,
} from "../../../lib/users/context";
import { addUserAction } from "../../../lib/users/action";
import { RoutesEnum } from "../../../lib/router/routes";
import { StyledFieldWrapper } from "../index.styles";

function CreateUserForm() {
  const user = useUserState();
  const dispatch = useUserStateAction();
  const navigate = useNavigate();
  const [submitError, setSubmitError] = useState("");

  useEffect(() => {
    if (isUserLoaded(user)) {
      navigate(RoutesEnum.PLAY_GET_STARTED);
    }
  }, [user, navigate]);

  const handleCreateUserSubmit = async (e: FormEvent<HTMLFormElement>) => {
    try {
      e.preventDefault();
      const data = Object.fromEntries(new FormData(e.currentTarget));
      const { userName } = data;

      const newUser = await handleCreateUser(userName.toString());
      if (!newUser?.id) {
        setSubmitError("Error creating user.");
      } else {
        dispatch(addUserAction(newUser));
        navigate(RoutesEnum.PLAY_GET_STARTED);
      }
    } catch (e: unknown) {
      const err = e as unknown as Error;
      setSubmitError(err.toString());
    }
  };

  return (
    <Form.Root className="FormRoot" onSubmit={handleCreateUserSubmit}>
      <StyledFieldWrapper>
        <Form.Field className="FormField" name="userName">
          <TextField.Root>
            <Form.Control asChild>
              <TextField.Input
                className="Input"
                type="text"
                required
                placeholder="Enter name"
                max={20}
              />
            </Form.Control>
          </TextField.Root>

          <Form.Message className="FormMessage" match="valueMissing">
            Please enter your name
          </Form.Message>
          <Form.Message className="FormMessage" match="typeMismatch">
            Please enter your name
          </Form.Message>
          {submitError !== "" && <Form.Message>{submitError}</Form.Message>}
        </Form.Field>
        <Form.Submit asChild>
          <Button className="Button">Create User</Button>
        </Form.Submit>
      </StyledFieldWrapper>
    </Form.Root>
  );
}

export default CreateUserForm;
