import { User } from "../../proto-ts/proto/model/user";
import {
  AddUserAction,
  RemoveUserAction,
  UpdateUserAction,
  UserActionType,
} from "./types";

export const addUserAction = (user: User): AddUserAction => ({
  type: UserActionType.ADD_USER,
  user,
});

export const removeUserAction = (): RemoveUserAction => ({
  type: UserActionType.REMOVE_USER,
});

export const updateUserAction = (name: string): UpdateUserAction => ({
  type: UserActionType.UPDATE_USER,
  name,
});
