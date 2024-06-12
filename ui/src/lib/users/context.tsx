import {
  Dispatch,
  createContext,
  useContext,
  useReducer,
} from "react";
import {
  AddUserAction,
  UpdateUserAction,
  UserActionType,
  UserLoaded,
  UserState,
  UserStateAction,
  UserStateProps,
} from "./types";

export const isUserLoaded = (user: UserState): user is UserLoaded => {
  if (!user?.id || !user?.id?.value || !user?.createdAt || !user?.name) {
    return false;
  }
  return true;
};

const reducer = (state: UserState, action: UserStateAction): UserState => {
  const { type } = action;

  switch (type) {
    case UserActionType.ADD_USER:
      const { user } = action as AddUserAction;
      return user;
    case UserActionType.REMOVE_USER:
      return null;
    case UserActionType.UPDATE_USER:
      const { name } = action as UpdateUserAction;
      return state?.id ? { ...state, name } : null;
    default:
      throw new Error(`Invalid action type: ${type}`);
  }
};

const userStateContext = createContext<UserState>(null);
const userStateActionContext = createContext<Dispatch<UserStateAction> | null>(
  null
);

const { Provider: StateProvider } = userStateContext;
const { Provider: DispatchProvider } = userStateActionContext;

export const UserStateProvider: React.FunctionComponent<UserStateProps> = ({
  children,
  defaultState = null,
}: UserStateProps): JSX.Element => {
  const [state, dispatch] = useReducer(reducer, defaultState);

  return (
    <StateProvider value={state}>
      <DispatchProvider value={dispatch}>{children}</DispatchProvider>
    </StateProvider>
  );
};

export const useUserState = (): UserState => {
  const state = useContext(userStateContext);

  if (state === undefined) {
    throw new Error(
      "userUserState must be used within a <UserStateProvider />"
    );
  }

  return state;
};

export const useUserStateAction = (): Dispatch<UserStateAction> => {
  const context = useContext(userStateActionContext);

  if (context === undefined || context === null) {
    throw new Error(
      "userUserStateAction must be used within a <UserStateProvider />"
    );
  }

  return context;
};
