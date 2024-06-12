import React from "react";
import ReactDOM from "react-dom/client";

import { Theme } from "@radix-ui/themes";

import "@radix-ui/themes/styles.css";

import {
  RouteObject,
  RouterProvider,
  createBrowserRouter,
} from "react-router-dom";
import UserCreate from "./pages/User/Create/index.tsx";
import { UserStateProvider } from "./lib/users/context.tsx";
import GetStarted from "./pages/Play/GetStarted/index.tsx";
import Board from "./pages/Play/Board/index.tsx";
import { boardLoader } from "./pages/Play/Board/loader.ts";
import "./index.css";
import AuthenticatedUserGuard from "./components/Guards/AuthenticatedUserGuard/index.tsx";
import PageNotFound from "./pages/Error/NotFound/index.tsx";
import { RoutesEnum } from "./lib/router/routes.ts";
import ToastNotification from "./components/ToastNotification/index.tsx";
import ViewDrawing from "./pages/Play/ViewDrawing/index.tsx";

const routes: RouteObject[] = [
  {
    path: RoutesEnum.ROOT,
    element: <AuthenticatedUserGuard />,
    children: [
      {
        path: RoutesEnum.PLAY_BOARD_ID,
        element: <Board />,
        loader: boardLoader,
        errorElement: <PageNotFound />,
      },
      {
        path: RoutesEnum.PLAY_GET_STARTED,
        element: <GetStarted />,
      },
      {
        path: RoutesEnum.PLAY_VIEW_DRAWING,
        element: <ViewDrawing />,
      },
      {
        path: RoutesEnum.USER_CREATE,
        element: <UserCreate />,
      },
      {
        path: RoutesEnum.ALL,
        element: <PageNotFound />,
      },
    ],
  },
  {
    path: RoutesEnum.ALL,
    element: <PageNotFound />,
  },
];

const router = createBrowserRouter(routes, {
  // basename: "/app",
});

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <UserStateProvider defaultState={null}>
      <Theme
        appearance="dark"
        panelBackground="solid"
        accentColor="violet"
        grayColor="gray"
        radius="full"
      >
        <RouterProvider router={router} />
        <ToastNotification />
      </Theme>
    </UserStateProvider>
  </React.StrictMode>
);
