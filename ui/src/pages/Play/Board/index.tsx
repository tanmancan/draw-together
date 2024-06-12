import {
  KeyboardEventHandler,
  useCallback,
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
} from "react";
import WsBaseClient from "../../../lib/websocket";
import { useLoaderData } from "react-router-dom";
import { useUserState } from "../../../lib/users/context";
import {
  BASE_VIEWPORT_HEIGHT,
  BASE_VIEWPORT_WIDTH,
  ViewPort,
} from "../../../lib/viewport";
import {
  StyledWrapper,
  StyledViewPortArea,
  StyledViewPort,
  StyledSideBarArea,
  StyledInputWrapper,
  StyledTextFieldRoot,
  StyledTextFieldInput,
  StyledSubmitButton,
  StyledCaptionWrapper,
} from "./index.styles";
import {
  handleGetBoardMessages,
  handleSendChatMessage,
} from "../../../lib/chat/handlers";
import { EventChatMessage } from "../../../proto-ts/proto/model/chat";
import ChatBox from "../../../components/ChatBox";
import { Board, EventBoardUpdate } from "../../../proto-ts/proto/model/board";
import ViewPortControls from "../../../components/ViewPortControls";
import BoardUserList from "../../../components/BoardUserList";
import {
  IDrawingUpdateEventPayload,
  IPointerUpdateEventPayload,
  ViewPortEventsEnum,
} from "../../../lib/viewport/event";
import { handlePointerUpdate } from "../../../lib/pointer/handlers";
import { User } from "../../../proto-ts/proto/model/user";
import { UUID } from "../../../proto-ts/proto/model/common";
import {
  EventDrawingDetect,
  EventDrawingUpdate,
  GetBoardDrawingsResponse,
  ImageData as ProtoImageData,
} from "../../../proto-ts/proto/model/drawing";
import {
  handleDrawingDetect,
  handleDrawingUpdate,
} from "../../../lib/boards/handlers";
import { WebsocketInputEventsEnum } from "../../../lib/websocket/events";
import { EventPointerUpdate } from "../../../proto-ts/proto/model/pointer";
import {
  NotificationSeverityEnum,
  showNotification,
} from "../../../components/ToastNotification/events";
import { Blockquote, Button, IconButton } from "@radix-ui/themes";
import { CrossCircledIcon } from "@radix-ui/react-icons";

const ENABLE_DETECT = false;

let pointerTid: number | null = null;

const checkBoardUser =
  (currentBoardUsers: User[]) =>
  (userID?: UUID): User | null => {
    const found = currentBoardUsers.find((u) => {
      return u.id?.value === userID?.value;
    });

    return found ?? null;
  };

function PageBoard() {
  const user = useUserState();
  const wsClient = useRef(new WsBaseClient());
  const viewPort = useRef(new ViewPort());
  const [board, boardDrawings]: [Board, GetBoardDrawingsResponse | null] =
    useLoaderData() as unknown as [Board, GetBoardDrawingsResponse | null];
  const { id, boardUsers } = board;
  const { value: boardID } = id as { value: string };
  const [currentBoardUsers, setCurrentBoardUsers] =
    useState<User[]>(boardUsers);
  const [messageList, setMessageList] = useState<EventChatMessage[]>([]);
  const [messageListError, setMessageListError] = useState("");
  const [message, setMessage] = useState("");
  const [detectCaption, setDetectCaption] = useState("");
  const [loadingCaption, setLoadingCaption] = useState(false);
  const viewPortRef = useRef<HTMLCanvasElement | null>(null);
  const getBoardUser = checkBoardUser(currentBoardUsers);

  const pointerUpdateEventHandler = useCallback(
    (e: CustomEventInit<IPointerUpdateEventPayload>) => {
      const { pointerStack } = e.detail ?? {};
      if (!pointerStack) return;
      if (pointerTid) {
        window.clearTimeout(pointerTid);
      }
      pointerTid = window.setTimeout(async () => {
        await handlePointerUpdate(pointerStack, board.id?.value ?? "");
      }, 100);
    },
    [board.id?.value]
  );

  const drawingUpdateEventHandler = useCallback(
    (e: CustomEventInit<IDrawingUpdateEventPayload>) => {
      const { imageData } = e.detail ?? {};
      if (!imageData) return;

      const proto: ProtoImageData = ProtoImageData.create();
      proto.data = imageData;

      (async () => {
        await handleDrawingUpdate(boardID, proto);
      })();
    },
    [boardID]
  );

  const onChatMessageInput = useCallback(
    (e: CustomEventInit<EventChatMessage>) => {
      const chatMessage = e.detail;
      if (chatMessage) {
        setMessageList((m) => [...m, chatMessage]);
      }
    },
    []
  );

  const onPointerUpdateInput = useCallback(
    (e: CustomEventInit<EventPointerUpdate>) => {
      const pointerUpdate = e.detail;
      const { metadata, pointerPositions = [] } = pointerUpdate ?? {};
      const { senderId } = metadata ?? {};

      if (senderId?.value === user?.id?.value) return;

      const pointerUser = getBoardUser(senderId);
      if (pointerUser) {
        viewPort?.current.updateExternalPointer(pointerUser, pointerPositions);
      }
    },
    [getBoardUser, user]
  );

  const onDrawingUpdateInput = useCallback(
    (e: CustomEventInit<EventDrawingUpdate>) => {
      const drawingUpdate = e.detail ?? {};
      const { metadata, imageData } = drawingUpdate ?? {};
      const { senderId } = metadata ?? {};

      if (senderId?.value === user?.id?.value) return;

      const drawingUser = getBoardUser(senderId);
      if (drawingUser && senderId?.value) {
        const { data } = imageData ?? {};
        if (data) {
          const blob = new Blob([data]);
          (async () => {
            const imageBitmap = await createImageBitmap(blob);
            viewPort?.current.updateExternalUserCanvas(
              senderId?.value,
              imageBitmap
            );
          })();
        }
      }
    },
    [getBoardUser, user]
  );

  const onDrawingDetectInput = useCallback(
    (e: CustomEventInit<EventDrawingDetect>) => {
      const drawingDetect = e.detail;
      const { description } = drawingDetect ?? {};
      if (description && description.length > 0) {
        setDetectCaption(description);
      }
    },
    []
  );

  const onBoardUpdate = useCallback((e: CustomEventInit<EventBoardUpdate>) => {
    const boardUpdate = e.detail ?? {};
    const { board } = boardUpdate ?? {};
    const { boardUsers } = board ?? {};
    if (boardUsers) {
      setCurrentBoardUsers((currentUsers) => {
        let newUser: User | null = null;
        for (let i = 0; i < boardUsers.length; i++) {
          for (let j = 0; j < currentUsers.length; j++) {
            if (boardUsers[i].id?.value !== currentUsers[j].id?.value) {
              newUser = boardUsers[i];
            }
          }
        }

        if (newUser) {
          showNotification(
            "new user joined the board",
            `user: ${newUser.name} was added to the board`,
            NotificationSeverityEnum.INFO
          );
        }

        return boardUsers;
      });
    }
  }, []);

  useLayoutEffect(() => {
    const viewPortCurrent = viewPortRef.current;
    (async () => {
      if (viewPortCurrent && boardID && user?.id?.value) {
        const uID = user?.id?.value;
        let imageBitmap: ImageBitmap | undefined = undefined;
        if (boardDrawings) {
          const imageData = boardDrawings.drawings?.[uID];
          const { data } = imageData ?? {};
          if (data) {
            const blob = new Blob([data]);
            imageBitmap = await createImageBitmap(blob);
          }
        }
        viewPort?.current.init(
          viewPortCurrent,
          user.id.value,
          boardID,
          imageBitmap
        );

        currentBoardUsers.forEach((bu) => {
          const { id } = bu;
          const uID = id?.value;
          if (uID) {
            viewPort?.current.addExternalUserLayer(uID);
            if (boardDrawings) {
              const imageData = boardDrawings.drawings?.[uID];
              const { data } = imageData ?? {};
              if (data) {
                const blob = new Blob([data]);
                (async () => {
                  const imageBitmap = await createImageBitmap(blob);
                  viewPort?.current.updateExternalUserCanvas(uID, imageBitmap);
                })();
              }
            }
          }
        });

        viewPortCurrent.addEventListener(
          ViewPortEventsEnum.POINTER_UPDATE_EVENT,
          pointerUpdateEventHandler
        );

        viewPortCurrent.addEventListener(
          ViewPortEventsEnum.DRAWING_UPDATE_EVENT,
          drawingUpdateEventHandler
        );
      }
    })();

    return () => {
      viewPortCurrent?.removeEventListener(
        ViewPortEventsEnum.POINTER_UPDATE_EVENT,
        pointerUpdateEventHandler
      );
      viewPortCurrent?.removeEventListener(
        ViewPortEventsEnum.DRAWING_UPDATE_EVENT,
        drawingUpdateEventHandler
      );
    };
  }, [
    boardID,
    user?.id?.value,
    drawingUpdateEventHandler,
    pointerUpdateEventHandler,
    boardDrawings,
    currentBoardUsers,
  ]);

  useEffect(() => {
    if (user?.id?.value && boardID) {
      wsClient?.current.init(
        { boardID: boardID },
        undefined,
        undefined,
        undefined,
        undefined,
        false
      );
    }

    window.addEventListener(
      WebsocketInputEventsEnum.CHAT_MESSAGE,
      onChatMessageInput
    );
    window.addEventListener(
      WebsocketInputEventsEnum.POINTER_UPDATE,
      onPointerUpdateInput
    );
    window.addEventListener(
      WebsocketInputEventsEnum.DRAWING_UPDATE,
      onDrawingUpdateInput
    );
    window.addEventListener(
      WebsocketInputEventsEnum.BOARD_UPDATE,
      onBoardUpdate
    );
    window.addEventListener(
      WebsocketInputEventsEnum.DRAWING_DETECT,
      onDrawingDetectInput
    );

    return () => {
      window.removeEventListener(
        WebsocketInputEventsEnum.CHAT_MESSAGE,
        onChatMessageInput
      );
      window.removeEventListener(
        WebsocketInputEventsEnum.POINTER_UPDATE,
        onPointerUpdateInput
      );
      window.removeEventListener(
        WebsocketInputEventsEnum.DRAWING_UPDATE,
        onDrawingUpdateInput
      );
      window.removeEventListener(
        WebsocketInputEventsEnum.BOARD_UPDATE,
        onBoardUpdate
      );
      window.removeEventListener(
        WebsocketInputEventsEnum.DRAWING_DETECT,
        onDrawingDetectInput
      );
    };
  }, [
    boardID,
    user?.id?.value,
    onChatMessageInput,
    onPointerUpdateInput,
    onDrawingUpdateInput,
    onBoardUpdate,
    onDrawingDetectInput,
  ]);

  useEffect(() => {
    (async () => {
      const chatMsgCollection = await handleGetBoardMessages(boardID);

      setMessageList(chatMsgCollection);
    })();
  }, [boardID]);

  const sendMessage = (message: string) => {
    if (message.trim().length > 0 && user?.id) {
      (async () => {
        try {
          if (!board?.id) return;
          const success = await handleSendChatMessage(boardID, [message]);
          if (success) {
            setMessage("");
            setMessageListError("");
            return;
          }

          setMessageListError("error message not sent");
        } catch (error) {
          const { message } = error as unknown as { message: string };
          setMessageListError(`error: ${message}`);
        }
      })();

      setMessage("");
    }
  };

  const handleSendMessage = () => {
    sendMessage(message);
  };

  const handleOnKeyDown: KeyboardEventHandler<HTMLInputElement> = (e) => {
    if (e.key === "Enter" && message.length > 0) {
      sendMessage(message);
    }
  };

  const handleDetect = () => {
    setDetectCaption("");
    setLoadingCaption(true);

    (async () => {
      await handleDrawingDetect(boardID);
      setLoadingCaption(false);
    })();
  };

  const handleCloseCaption = () => {
    setDetectCaption("");
  };

  return (
    <StyledWrapper>
      <StyledViewPortArea>
        <StyledViewPort
          id="viewport"
          ref={viewPortRef}
          width={BASE_VIEWPORT_WIDTH}
          height={BASE_VIEWPORT_HEIGHT}
        />
        {ENABLE_DETECT && detectCaption && detectCaption.length > 0 && (
          <StyledCaptionWrapper align="center">
            <Blockquote>{detectCaption}</Blockquote>
            <IconButton
              onClick={handleCloseCaption}
              ml={"auto"}
              variant="ghost"
            >
              <CrossCircledIcon />
            </IconButton>
          </StyledCaptionWrapper>
        )}
        {ENABLE_DETECT && (
          <Button disabled={loadingCaption} onClick={handleDetect}>
            AI Guess Drawing
          </Button>
        )}
        <ViewPortControls viewPort={viewPort?.current} board={board} />
      </StyledViewPortArea>
      <StyledSideBarArea>
        <ChatBox
          messages={messageList}
          boardUsers={currentBoardUsers}
          error={messageListError}
        />
        <StyledInputWrapper>
          <StyledTextFieldRoot>
            <StyledTextFieldInput
              type="text"
              id="message"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              onKeyDown={handleOnKeyDown}
            />
            <StyledSubmitButton
              disabled={message.length === 0}
              onClick={handleSendMessage}
            >
              Send
            </StyledSubmitButton>
          </StyledTextFieldRoot>
        </StyledInputWrapper>
        <BoardUserList board={board} currentUserList={currentBoardUsers} />
      </StyledSideBarArea>
    </StyledWrapper>
  );
}

export default PageBoard;
