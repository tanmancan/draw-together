import { BASE_VIEWPORT_HEIGHT, BASE_VIEWPORT_WIDTH } from ".";
import { User } from "../../proto-ts/proto/model/user";
import { PointerStack } from "./types";

type PointerData = [
  OffscreenCanvas,
  OffscreenCanvasRenderingContext2D,
  PointerStack,
  OffscreenCanvas
];

export const getRandomColor = (maxTint: number = 150): string => {
  const r = Math.floor(Math.random() * maxTint);
  const g = Math.floor(Math.random() * maxTint);
  const b = Math.floor(Math.random() * maxTint);
  return `rgb(${r} ${g} ${b})`;
};

export class ViewPortPointer {
  protected primaryCanvas: OffscreenCanvas;
  protected primaryCtx: OffscreenCanvasRenderingContext2D | null;

  protected userCanvasList: Record<string, OffscreenCanvas> = {};
  protected userCtxList: Record<string, OffscreenCanvasRenderingContext2D> = {};
  protected userPointerList: Record<string, PointerStack> = {};
  protected userPointerImage: Record<string, OffscreenCanvas> = {};

  protected pointerCycle = 1000 / 20;
  protected prevTs = 0;
  protected rid: number | null = null;

  constructor(
    protected width = BASE_VIEWPORT_WIDTH,
    protected height = BASE_VIEWPORT_HEIGHT
  ) {
    this.primaryCanvas = new OffscreenCanvas(this.width, this.height);
    this.primaryCtx = this.primaryCanvas.getContext("2d");
  }

  protected addUser(user: User, initPointerStack: PointerStack = []) {
    const canvas = new OffscreenCanvas(this.width, this.height);
    const ctx = canvas.getContext("2d");
    const userID = user?.id?.value;
    if (!userID) return;

    this.userPointerImage[userID] = this.renderPointerImage(user);
    this.userCanvasList[userID] = canvas;
    if (ctx) {
      ctx.fillStyle = getRandomColor();
      ctx.font = "12px sans-serif";
      this.userCtxList[userID] = ctx;
    }
    this.userPointerList[userID] = initPointerStack;
  }

  updatePointer(user: User, pointerStack: PointerStack) {
    const userID = user?.id?.value ?? "";
    if (!this.userCanvasList?.[userID] || !this.userCtxList?.[userID]) {
      this.addUser(user, pointerStack);
      return;
    }
    this.userPointerList[userID] = pointerStack;
  }

  protected getUserInfo(userID: string): PointerData {
    return [
      this.userCanvasList[userID],
      this.userCtxList[userID],
      this.userPointerList?.[userID] ?? [],
      this.userPointerImage?.[userID],
    ];
  }

  protected renderPointerImage(user: User): OffscreenCanvas {
    const x = 0;
    const y = 0;
    const basePointerSize = 15;
    const textHeight = 12;
    const fontStyle = `${textHeight}px sans-serif`;
    const textWidth = 52;
    const canvasWidth = Math.ceil(basePointerSize * 1.5) + textWidth;
    const canvasHeight =
      Math.ceil(basePointerSize * 1.25) + Math.ceil(textHeight / 2);

    const maxNameLen = 10;
    let name = user.name.trim();

    if (name.length > maxNameLen - 3) {
      name = name.slice(0, maxNameLen - 3).trim() + "...";
    }

    const canvas = new OffscreenCanvas(canvasWidth, canvasHeight);
    const ctx = canvas.getContext("2d");

    if (ctx) {
      ctx.fillStyle = getRandomColor();
      ctx.font = fontStyle;

      ctx.beginPath();
      ctx.moveTo(x, y);
      ctx.lineTo(x + basePointerSize, y + basePointerSize);
      ctx.lineTo(x, y + basePointerSize);
      ctx.lineTo(x + Math.floor(basePointerSize / 2.75), y + basePointerSize);
      ctx.lineTo(x, Math.floor(y + basePointerSize * 1.45));
      ctx.lineTo(x, y);
      ctx.fill();
      ctx.closePath();

      ctx.fillText(
        name,
        x + Math.floor(basePointerSize * 1.25),
        y + Math.floor(basePointerSize * 1.25),
        textWidth
      );
    }

    return canvas;
  }

  protected renderUserPointer(userID: string) {
    const userInfo = this.getUserInfo(userID);
    if (!userInfo) return;

    const [, ctx, pointerStack, pointerImg] = userInfo;

    if (pointerStack.length === 0) return;
    const pointer = pointerStack.shift();

    if (!pointer) return;

    const { x, y } = pointer;

    ctx.clearRect(0, 0, this.width, this.height);
    ctx.drawImage(pointerImg, x, y);
  }

  render(now: number) {
    if (!this.primaryCtx) return;

    const elapsed = now - this.prevTs;

    if (elapsed > this.pointerCycle) {
      this.pointerCycle = now - (elapsed % this.pointerCycle);
      this.primaryCtx.clearRect(0, 0, this.width, this.height);

      Object.keys(this.userCanvasList).forEach((u) => {
        this.renderUserPointer(u);
        const [canvas] = this.getUserInfo(u);
        this.primaryCtx?.drawImage(canvas, 0, 0);
      });
    }
  }

  getCanvas(): OffscreenCanvas {
    return this.primaryCanvas;
  }
}
