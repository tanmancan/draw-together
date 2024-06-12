import { ReactNode } from "react";
import { User } from "../../proto-ts/proto/model/user";
import { Timestamp } from "../../proto-ts/google/protobuf/timestamp";

export interface Token {
  token: string;
}
export type UserState = UserLoaded | User | null;

export interface UserLoaded extends User {
  id: {
    value: string;
  };
  name: string;
  createdAt: Timestamp;
}

export interface UserStateProps {
  children: ReactNode | ReactNode[];
  defaultState: UserState;
}

export enum UserActionType {
  ADD_USER,
  REMOVE_USER,
  UPDATE_USER,
}

export interface AddUserAction {
  type: UserActionType.ADD_USER;
  user: UserState;
}

export interface RemoveUserAction {
  type: UserActionType.REMOVE_USER;
}

export interface UpdateUserAction {
  type: UserActionType.UPDATE_USER;
  name: string;
}

export type UserStateAction =
  | AddUserAction
  | RemoveUserAction
  | UpdateUserAction;
