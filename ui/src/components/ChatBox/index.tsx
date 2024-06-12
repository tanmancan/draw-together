import { Box, Callout, Card, Flex } from "@radix-ui/themes";
import { IChatBoxProps } from "./types";
import { useUserState } from "../../lib/users/context";
import { UserState } from "../../lib/users/types";
import { useLayoutEffect, useRef } from "react";
import {
  StyledMessageBody,
  StyledScrollBox,
  StyledUserName,
} from "./index.styles";
import { User } from "../../proto-ts/proto/model/user";

const checkCurrentUser =
  (user: UserState) =>
  (uid?: string): boolean => {
    return user?.id?.value === uid;
  };

const checkBoardUser =
  (boardUsers: User[]) =>
  (uid?: string): User | null => {
    return boardUsers.find((u) => u.id?.value === uid) ?? null;
  };

function ChatBox(props: IChatBoxProps) {
  const { boardUsers, error, messages } = props;
  const user = useUserState();
  const isCurrentUser = checkCurrentUser(user);
  const getBoardUser = checkBoardUser(boardUsers);
  const messageListRef = useRef<HTMLDivElement | null>(null);

  useLayoutEffect(() => {
    if (messageListRef.current) {
      messageListRef.current.scrollTop = messageListRef.current.scrollHeight;
    }
  }, [messages.length, error]);

  return (
    <Card>
      <StyledScrollBox ref={messageListRef}>
        <Flex direction="column" gap="2">
          {messages
            .sort((ma, mb) =>
              (ma.metadata?.createdAt?.seconds ?? 0) >
              (mb.metadata?.createdAt?.seconds ?? 0)
                ? 1
                : -1
            )
            .map((m) => (
              <Box
                key={m.metadata?.id?.value}
                ml={
                  isCurrentUser(m.metadata?.senderId?.value)
                    ? "auto"
                    : undefined
                }
              >
                {!isCurrentUser(m.metadata?.senderId?.value) && (
                  <StyledUserName size="1">
                    {getBoardUser(m.metadata?.senderId?.value)?.name}
                  </StyledUserName>
                )}
                <Callout.Root
                  color={
                    isCurrentUser(m.metadata?.senderId?.value)
                      ? "blue"
                      : undefined
                  }
                  size="1"
                  style={{
                    width: "fit-content",
                  }}
                >
                  <StyledMessageBody>{m.body}</StyledMessageBody>
                </Callout.Root>
              </Box>
            ))}
          {error && error.length > 0 && (
            <Callout.Root color="red">
              <Callout.Text>{error}</Callout.Text>
            </Callout.Root>
          )}
        </Flex>
      </StyledScrollBox>
    </Card>
  );
}

export default ChatBox;
