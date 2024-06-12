import { useEffect, useState } from "react";
// import { Outlet, useNavigate } from "react-router-dom";
import { addUserAction } from "../../../lib/users/action";
import {
  isUserLoaded,
  useUserState,
  useUserStateAction,
} from "../../../lib/users/context";
import { handleGetUser } from "../../../lib/users/handlers";
import App from "../../../App";
import { useNavigate } from "react-router-dom";
import LoadingIndicator from "../../LoadingIndicator";
import { RoutesEnum } from "../../../lib/router/routes";

function AuthenticatedUserGuard() {
  const user = useUserState();
  const [loading, setLoading] = useState(true);
  const dispatch = useUserStateAction();
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      try {
        if (isUserLoaded(user)) {
          setLoading(false);
          return;
        }

        setLoading(true);
        const loadedUser = await handleGetUser();

        if (!loadedUser?.id?.value) {
          console.error("user load error: invalid user id");
          setLoading(false);
          navigate(RoutesEnum.USER_CREATE);
        } else {
          dispatch(addUserAction(loadedUser));
          setLoading(false);
        }
      } catch (error) {
        console.error(error);
        setLoading(false);
        navigate(RoutesEnum.USER_CREATE);
      }
    })();
  }, [user, dispatch, navigate]);

  return <>{loading ? <LoadingIndicator /> : <App />}</>;
}

export default AuthenticatedUserGuard;
