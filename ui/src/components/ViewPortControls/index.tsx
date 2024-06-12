import {
  TrashIcon,
  ResetIcon,
  MixerHorizontalIcon,
  InfoCircledIcon,
  CopyIcon,
  CheckCircledIcon,
  DownloadIcon,
  EyeOpenIcon,
} from "@radix-ui/react-icons";
import {
  Button,
  Flex,
  IconButton,
  Popover,
  Separator,
  Slider,
  Strong,
  TextField,
  Tooltip,
} from "@radix-ui/themes";
import { styled } from "@stitches/react";
import { useState } from "react";
import { HexColorPicker } from "react-colorful";
import {
  BASE_STROKE_STYLE,
  BASE_STROKE_WIDTH,
  ViewPort,
} from "../../lib/viewport";
import { Board } from "../../proto-ts/proto/model/board";
import { useNavigate } from "react-router-dom";
import { RoutesEnum } from "../../lib/router/routes";

const SHOW_DEBUG_CONTROL = false;

export const StyledViewPortControlsWrapper = styled("div", {
  position: "absolute",
  top: 16,
  right: 16,
  zIndex: 2,
  height: "fit-content",
  display: "flex",
  gap: 4,
});

interface IViewPortControls {
  viewPort: ViewPort;
  board: Board;
}

function ViewPortControls(props: IViewPortControls) {
  const { viewPort, board } = props;
  const navigate = useNavigate();
  const [strokeStyle, setStrokeStyle] = useState(BASE_STROKE_STYLE);
  const [strokeWidth, setStrokeWidth] = useState(BASE_STROKE_WIDTH);
  const [copySuccess, setCopySuccess] = useState(false);

  const handleClearViewPort = () => {
    viewPort.clear(true);
  };

  const handleUndoDrawing = () => {
    viewPort.undoHistory();
  };

  const handleColorPicker = (colorPicked: string) => {
    viewPort.setStrokeStyle(colorPicked);
    setStrokeStyle(colorPicked);
  };

  const handleChangeStrokeWidth = (val: number[]) => {
    const width = val?.[0] ?? BASE_STROKE_WIDTH;

    setStrokeWidth(width);
    viewPort.setStrokeWidth(width);
  };

  const handleToggleDebug = () => {
    viewPort.toggleDebug();
  };

  const handleBoardIdCopy = () => {
    if (copySuccess) {
      return;
    }
    (async () => {
      if (
        board.id?.value &&
        typeof window.navigator.clipboard.writeText === "function"
      ) {
        await window.navigator.clipboard.writeText(board.id?.value);
        setCopySuccess(true);
        window.setTimeout(() => {
          setCopySuccess(false);
        }, 2000);
      }
    })();
  };

  function handleSaveImage() {
    viewPort.saveImage();
  }

  function handleViewImage() {
    const imgDataUrl = viewPort.getImageDataUrl();
    if (imgDataUrl) {
      const url = RoutesEnum.PLAY_VIEW_DRAWING.replace(
        ":data",
        btoa(imgDataUrl)
      );
      navigate(url);
    }
  }

  return (
    <StyledViewPortControlsWrapper>
      <Tooltip content="Clear Drawing">
        <IconButton onClick={handleClearViewPort}>
          <TrashIcon />
        </IconButton>
      </Tooltip>
      <Tooltip content="Undo Stroke">
        <IconButton onClick={handleUndoDrawing}>
          <ResetIcon />
        </IconButton>
      </Tooltip>
      <Popover.Root>
        <Tooltip content="Configure Brush">
          <Popover.Trigger>
            <IconButton>
              <MixerHorizontalIcon />
            </IconButton>
          </Popover.Trigger>
        </Tooltip>
        <Popover.Content>
          <Flex direction={"column"} gap={"4"}>
            <Strong>Board ID:</Strong>
            <TextField.Root>
              <TextField.Input defaultValue={board?.id?.value} readOnly />
              <TextField.Slot>
                <Tooltip
                  content={copySuccess ? "Copied" : "Copy to clipboard"}
                  delayDuration={0}
                >
                  <IconButton
                    onClick={handleBoardIdCopy}
                    size="1"
                    variant="ghost"
                  >
                    {copySuccess ? (
                      <CheckCircledIcon color="green" height="14" width="14" />
                    ) : (
                      <CopyIcon height="14" width="14" />
                    )}
                  </IconButton>
                </Tooltip>
              </TextField.Slot>
            </TextField.Root>
            <Separator size={"4"} />
            <Strong>Width: {strokeWidth}px</Strong>
            <Slider
              defaultValue={[strokeWidth]}
              step={1}
              min={1}
              max={100}
              onValueChange={handleChangeStrokeWidth}
            />
            <Separator size={"4"} />
            <Strong>Color: {strokeStyle}</Strong>
            <HexColorPicker color={strokeStyle} onChange={handleColorPicker} />
            <Popover.Close>
              <Button size="1">Close</Button>
            </Popover.Close>
          </Flex>
        </Popover.Content>
      </Popover.Root>
      <Tooltip content="Download Drawing">
        <IconButton onClick={handleSaveImage}>
          <DownloadIcon />
        </IconButton>
      </Tooltip>
      <Tooltip content="View Drawing">
        <IconButton onClick={handleViewImage}>
          <EyeOpenIcon />
        </IconButton>
      </Tooltip>
      {SHOW_DEBUG_CONTROL && (
        <Tooltip content="Show Debug Info">
          <IconButton onClick={handleToggleDebug}>
            <InfoCircledIcon />
          </IconButton>
        </Tooltip>
      )}
    </StyledViewPortControlsWrapper>
  );
}

export default ViewPortControls;
