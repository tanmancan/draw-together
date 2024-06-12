import { EventChatMessage } from "../../proto-ts/proto/model/chat";
import { User } from "../../proto-ts/proto/model/user";

export interface IChatBoxProps {
  boardUsers: User[];
  error?: string;
  messages: EventChatMessage[];
}
