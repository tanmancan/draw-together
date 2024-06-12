import { LoaderFunctionArgs } from "react-router-dom";
import {
  handleGetBoard,
  handleGetBoardDrawings,
} from "../../../lib/boards/handlers";

export const boardLoader = async ({ params }: LoaderFunctionArgs<null>) => {
  const board = await handleGetBoard(params?.id ?? "none");
  const { id } = board ?? {};
  const { value: boardID } = id ?? {};

  if (!boardID) {
    throw new Error("board not found");
  }

  const boardDrawings = await handleGetBoardDrawings(boardID);
  return [board, boardDrawings];
};
