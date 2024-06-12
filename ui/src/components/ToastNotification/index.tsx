import * as Toast from "@radix-ui/react-toast";
import {
  StyledCalloutIcon,
  StyledCalloutRoot,
  StyledCalloutText,
  StyledToastClose,
  StyledToastViewport,
  StyledToastWrapper,
} from "./index.styles";
import { useCallback, useEffect, useState } from "react";
import {
  INotificationPayload,
  NOTIFICATION_EVENT,
  NotificationSeverityEnum,
} from "./events";
import { Box, Button, Flex, IconButton, Text } from "@radix-ui/themes";
import {
  CrossCircledIcon,
  ExclamationTriangleIcon,
  InfoCircledIcon,
} from "@radix-ui/react-icons";

const DEFAULT_DURATION = 3000;

const getSeverityIcon = (severity?: NotificationSeverityEnum) => {
  switch (severity) {
    case NotificationSeverityEnum.ERROR:
    case NotificationSeverityEnum.WARNING:
      return <ExclamationTriangleIcon />;
    case NotificationSeverityEnum.INFO:
    default:
      return <InfoCircledIcon />;
  }
};

function ToastNotification() {
  const [messages, setMessage] = useState<INotificationPayload[]>([]);

  const onOpenChangeHandler = useCallback((open: boolean) => {
    if (!open) {
      setMessage((mList) => {
        mList.pop();
        return [...mList];
      });
    }
  }, []);

  const onNotificationHandler = useCallback(
    (e: CustomEventInit<INotificationPayload>) => {
      if (e.detail) {
        setMessage((mList) => {
          const newMsg = [...mList];
          if (e?.detail) {
            newMsg.push(e.detail);
          }
          return newMsg;
        });
      }
    },
    []
  );

  useEffect(() => {
    window.addEventListener(NOTIFICATION_EVENT, onNotificationHandler);
    () => {
      window.removeEventListener(NOTIFICATION_EVENT, onNotificationHandler);
    };
  }, [onNotificationHandler]);

  const clearAllMessages = () => {
    setMessage([]);
  };

  const lastMessage = messages.at(-1);
  const lastMessageDuration = lastMessage?.duration
    ? lastMessage.duration
    : DEFAULT_DURATION;

  return (
    <Toast.Provider
      duration={
        lastMessage?.severity === NotificationSeverityEnum.ERROR
          ? 300000
          : lastMessageDuration
      }
    >
      <Toast.Root
        asChild
        open={messages.length > 0}
        onOpenChange={onOpenChangeHandler}
      >
        <Box>
          <StyledToastWrapper>
            <StyledCalloutRoot
              size="1"
              color={lastMessage?.severity}
              highContrast
            >
              <StyledCalloutIcon>
                {getSeverityIcon(lastMessage?.severity)}
              </StyledCalloutIcon>
              <StyledCalloutText>
                <Toast.Title asChild>
                  <Text size="2" weight="bold">
                    {lastMessage?.title}
                  </Text>
                </Toast.Title>
              </StyledCalloutText>
              <StyledCalloutText>
                <Toast.Description asChild>
                  <Text size="2">{lastMessage?.body}</Text>
                </Toast.Description>
              </StyledCalloutText>
              {messages.length > 1 && (
                <Flex align="baseline" gap="3">
                  <StyledCalloutText>
                    <Text size="1" color="gray">
                      1 of {messages.length}
                    </Text>
                  </StyledCalloutText>
                  <Button onClick={clearAllMessages} size="1" variant="ghost">
                    Clear All
                  </Button>
                </Flex>
              )}
              <StyledToastClose asChild>
                <IconButton size="1" variant="ghost">
                  <CrossCircledIcon />
                </IconButton>
              </StyledToastClose>
            </StyledCalloutRoot>
          </StyledToastWrapper>
        </Box>
      </Toast.Root>

      <StyledToastViewport data-name="toast-viewport" />
    </Toast.Provider>
  );
}

export default ToastNotification;
