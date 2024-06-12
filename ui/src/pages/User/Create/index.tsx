import { Box, Card, Heading } from "@radix-ui/themes";
import CreateUserForm from "../../../components/Forms/CreateUserForm";
import { StyledFormWrapper } from "../../../components/Forms/index.styles";

function UserCreate() {
  return (
    <StyledFormWrapper>
      <Box mb={"4"}>
        <Heading>Create User</Heading>
      </Box>
      <Card>
        <Box>
          <CreateUserForm />
        </Box>
      </Card>
    </StyledFormWrapper>
  );
}

export default UserCreate;
