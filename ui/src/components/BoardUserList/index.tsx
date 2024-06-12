import { Avatar, Box, Card, Flex, Text } from "@radix-ui/themes";
import { Board } from "../../proto-ts/proto/model/board";
import { User } from "../../proto-ts/proto/model/user";

interface IBoardUsersList {
  board: Board;
  currentUserList: User[];
}

function BoardUserList(props: IBoardUsersList) {
  const { board, currentUserList } = props;
  return (
    <Flex gap={"1"} direction={"column"}>
      {currentUserList?.map((boardUser) => {
        const { name, id } = boardUser;
        return (
          <Card key={boardUser?.id?.value}>
            <Flex gap="3" align="center">
              <Avatar size="1" radius="full" fallback={name[0]} />
              <Box>
                <Text as="div" size="2" weight="bold">
                  {name}{" "}
                  {id?.value === board?.owner?.id?.value && (
                    <Text size="1" weight={"light"}>
                      (admin)
                    </Text>
                  )}
                </Text>
              </Box>
            </Flex>
          </Card>
        );
      })}
    </Flex>
  );
}

export default BoardUserList;
