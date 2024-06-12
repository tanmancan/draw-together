import { User } from "../../proto-ts/proto/model/user";
import { ViewPortDebug } from "./debug";
import { DrawingUpdateEvent, PointerUpdateEvent } from "./event";
import { ViewPortPointer } from "./pointer";
import { PointerStack } from "./types";
import { ViewPortUser } from "./user";

export const BASE_VIEWPORT_WIDTH = 512;
export const BASE_VIEWPORT_HEIGHT = 512;
export const BASE_STROKE_WIDTH = 5;
export const BASE_STROKE_STYLE = "#000";

export class ViewPort {
  protected primaryCanvas: HTMLCanvasElement | null = null;
  protected primaryCtx: CanvasRenderingContext2D | null = null;

  protected boardID: string = "";
  protected currentUserID: string = "";
  protected currentUserCanvas: OffscreenCanvas | null = null;
  protected currentUserContext: OffscreenCanvasRenderingContext2D | null = null;

  // Canvas element dimension
  protected width: number = BASE_VIEWPORT_WIDTH;
  protected height: number = BASE_VIEWPORT_HEIGHT;

  // Mouse event state
  protected mouseDown: boolean = false;

  // Drawing event state
  protected beginPath: boolean = false;

  // requestAnimationFrame ID
  protected ridPriv: number | null = null;

  // Mouse offset for canvas element
  protected offsetX: number = 0;
  protected offsetY: number = 0;
  protected prevX: number = 0;
  protected prevY: number = 0;
  protected pointerMoved: boolean = false;

  // Drawing history stack for undo functionality
  protected historyStack: ImageData[] = [];

  // Max number of history allowed for undo
  protected maxHistory: number = 100;

  // History tracking state
  protected canvasChanged: boolean = false;

  // Line styles
  protected strokeWidth = BASE_STROKE_WIDTH;
  protected strokeStyle = BASE_STROKE_STYLE;

  // Render refresh speed controls
  protected prevCycle: number = 0;
  protected fpsCycle: number = 1000 / 60;
  protected currentUserPrevCycle: number = 0;
  protected currentUserFpsCycle: number = 1000 / 60;

  protected pointerStack: PointerStack = [];
  protected prevPointerCycle: number = 0;
  protected pointerUpdateCycle: number = 1000 / 20;
  protected pointerFlushThreshHold: number = 20;
  protected pointerPositionEmit: [number, number] = [0, 0];
  protected pointerEventTid: number | null = null;

  // Render debug view pane
  protected debugViewPort?: ViewPortDebug;
  protected enableDebugging: boolean = false;

  protected pointerViewPort: ViewPortPointer = new ViewPortPointer(
    this.width,
    this.height
  );

  protected userDrawingViewPort: ViewPortUser = new ViewPortUser(
    this.width,
    this.height
  );

  constructor() {
    this.onMouseUp = this.onMouseUp.bind(this);
    this.onMouseMove = this.onMouseMove.bind(this);
    this.render = this.render.bind(this);
    this.debugViewPort = new ViewPortDebug(this.enableDebugging);
  }

  init(
    canvas: HTMLCanvasElement,
    userID: string,
    boardID: string,
    imageData?: ImageBitmap
  ) {
    if (this.getPrimaryCanvas()) return;

    this.currentUserID = userID;
    this.boardID = boardID;

    this.initPrimaryCanvas.call(this, canvas);
    this.initCurrentUserCanvas.call(this);
    if (imageData) {
      this.updateCurrentUserCanvas.call(this, imageData);
    }
    this.initEvent.call(this);

    this.rid = requestAnimationFrame(this.render);
  }

  protected initPrimaryCanvas(canvas: HTMLCanvasElement) {
    this.primaryCanvas = canvas;
    this.updateCanvasSize();
    this.primaryCtx = this.primaryCanvas.getContext("2d", {
      willReadFrequently: true,
      alpha: false,
    }) as unknown as CanvasRenderingContext2D;
  }

  protected initCurrentUserCanvas() {
    this.currentUserCanvas = new OffscreenCanvas(this.width, this.height);
    this.currentUserContext = this.currentUserCanvas.getContext("2d", {
      willReadFrequently: true,
    });
  }

  protected updateCurrentUserCanvas(imageData: ImageBitmap) {
    const canvas = this.getCurrentUserCanvas();
    const ctx = this.getCurrentUserContext();
    if (!ctx || !canvas) return;

    if (ctx && canvas) {
      ctx.drawImage(imageData, 0, 0, this.width, this.height);
      imageData.close();
    }
  }

  protected initEvent() {
    document?.addEventListener("mouseup", this.onMouseUp);
    this.primaryCanvas?.addEventListener("mousemove", this.onMouseMove);
  }

  protected removeEvent() {
    document?.removeEventListener("mouseup", this.onMouseUp);
    this.primaryCanvas?.removeEventListener("mousemove", this.onMouseMove);
  }

  protected getCurrentUserCanvas() {
    return this.currentUserCanvas;
  }

  protected getCurrentUserContext() {
    return this.currentUserContext;
  }

  protected getPrimaryCanvas() {
    return this.primaryCanvas;
  }

  protected getPrimaryContext() {
    return this.primaryCtx;
  }

  protected emitPointerUpdateEvent() {
    if (this.pointerMoved) {
      const { x, y } = this.pos;
      this.pointerStack.push({ x, y });
      if (this.pointerStack.length >= this.pointerFlushThreshHold) {
        const detail = {
          boardID: this.boardID,
          userID: this.currentUserID,
          pointerStack: [...this.pointerStack],
        };
        const ev = new PointerUpdateEvent({
          detail,
        });
        this.primaryCanvas?.dispatchEvent(ev);
        this.pointerStack.length = 0;
      }
    }
  }

  protected emitDrawingUpdateEvent() {
    const canvas = this.getCurrentUserCanvas();
    if (!canvas) return;

    canvas
      .convertToBlob()
      .then((blob) => blob.arrayBuffer())
      .then((buffer) => {
        const data = new Uint8Array(buffer);
        const ev = new DrawingUpdateEvent({
          detail: {
            boardID: this.boardID,
            userID: this.currentUserID,
            imageData: data,
          },
        });
        this.getPrimaryCanvas()?.dispatchEvent(ev);
      });
  }

  protected renderCurrentUserCanvas(now: number) {
    const canvas = this.getCurrentUserCanvas();
    const ctx = this.getCurrentUserContext();
    if (!ctx || !canvas) return;
    const { x, y } = this.pos;

    const elapsed = now - this.currentUserPrevCycle;
    if (elapsed > this.currentUserFpsCycle) {
      this.currentUserPrevCycle = now - (elapsed % this.currentUserFpsCycle);

      if (!this.mouseDown || x > this.width || y > this.height) {
        this.beginPath = false;
        ctx?.closePath();
      }

      if (this.mouseDown) {
        this.canvasChanged = true;
        if (!this.beginPath) {
          ctx?.beginPath();
          this.beginPath = true;
        }

        ctx.lineWidth = this.strokeWidth;
        ctx.strokeStyle = this.strokeStyle;
        ctx.lineJoin = "round";
        ctx.lineCap = "round";
        ctx?.lineTo(x, y);
        ctx?.stroke();
      }
    }
  }

  protected render(now: number) {
    const primaryCtx = this.getPrimaryContext();
    if (!primaryCtx) {
      return;
    }

    this.cancelFrame();
    this.rid = requestAnimationFrame(this.render);

    const { x, y } = this.pos;

    const pointerElapsed = now - this.prevPointerCycle;

    if (pointerElapsed > this.pointerUpdateCycle) {
      this.pointerMoved = false;

      if (this.prevX != x) {
        this.prevX = x;
        this.pointerMoved = true;
      }
      if (this.prevY != y) {
        this.prevY = y;
        this.pointerMoved = true;
      }
      this.prevPointerCycle = now - (pointerElapsed & this.pointerUpdateCycle);
      this.emitPointerUpdateEvent();
    }

    const elapsed = now - this.prevCycle;
    if (elapsed > this.fpsCycle) {
      this.prevCycle = now - (elapsed % this.fpsCycle);
      this.debug(elapsed);

      this.renderCurrentUserCanvas(now);
      this.pointerViewPort.render(now);
      this.userDrawingViewPort.render(now);

      primaryCtx.fillStyle = "#fff";
      primaryCtx.fillRect(0, 0, this.width, this.height);

      primaryCtx.drawImage(this.pointerViewPort.getCanvas(), 0, 0);
      primaryCtx.drawImage(this.userDrawingViewPort.getCanvas(), 0, 0);

      const currentUserCanvas = this.getCurrentUserCanvas();
      if (currentUserCanvas) {
        primaryCtx.drawImage(currentUserCanvas, 0, 0);
      }
    }
  }

  protected updateCanvasSize() {
    if (!this.primaryCanvas || !this.primaryCtx) return;

    const rect = this.primaryCanvas.getBoundingClientRect();
    this.primaryCanvas.width = rect.width;
    this.primaryCanvas.height = rect.height;
    this.width = rect.width;
    this.height = rect.height;
  }

  protected saveHistory() {
    const canvas = this.getCurrentUserCanvas();
    const ctx = this.getCurrentUserContext();
    if (!ctx || !canvas) return;

    const imgData = ctx?.getImageData(0, 0, this.width, this.height);

    if (this.historyStack.length >= 100) {
      this.historyStack.shift();
    }

    this.historyStack.push(imgData);
    this.canvasChanged = false;

    this.emitDrawingUpdateEvent();
  }

  protected onMouseUp() {
    if (this.canvasChanged) {
      this.canvasChanged = false;
      this.mouseDown = false;
      requestAnimationFrame(this.saveHistory.bind(this));
    }
  }

  protected onMouseMove(e: MouseEvent) {
    this.offsetX = e.offsetX;
    this.offsetY = e.offsetY;
    this.mouseDown = e.buttons > 0 && this.mouseIsBound;
  }

  protected get canvasAspect(): number {
    const rect = this.primaryCanvas?.getBoundingClientRect();
    return rect?.width ? this.width / rect?.width : 1;
  }

  protected get pos(): { x: number; y: number } {
    const x = Math.floor(this.offsetX * this.canvasAspect);
    const y = Math.floor(this.offsetY * this.canvasAspect);

    return { x: x >= 0 ? x : 0, y: y >= 0 ? y : 0 };
  }

  get mouseIsBound(): boolean {
    const { x, y } = this.pos;
    return x <= this.width && y <= this.height && x >= 0 && y >= 0;
  }

  protected getRid() {
    return this.rid;
  }

  teardown() {
    this.removeEvent.call(this);

    const rid = this.getRid();
    if (rid) {
      cancelAnimationFrame(rid);
    }
  }

  addExternalUserLayer(userID: string) {
    if (this.currentUserID && userID !== this.currentUserID) {
      this.userDrawingViewPort.addUser(userID);
    }
  }

  updateExternalPointer(user: User, pointerStack: PointerStack) {
    if (this.currentUserID && user.id?.value === this.currentUserID) return;
    this.pointerViewPort.updatePointer(user, pointerStack);
  }

  updateExternalUserCanvas(userID: string, imageData: ImageBitmap) {
    if (this.currentUserID && userID === this.currentUserID) return;
    this.userDrawingViewPort.updateExternalUserCanvas(userID, imageData);
  }

  undoHistory() {
    const canvas = this.getCurrentUserCanvas();
    const ctx = this.getCurrentUserContext();
    if (!ctx || !canvas) return;

    this.historyStack.pop();
    if (this.historyStack.length > 0) {
      ctx.putImageData(this.historyStack[this.historyStack.length - 1], 0, 0);
      this.emitDrawingUpdateEvent();
    } else {
      this.clear(true);
    }
  }

  clear(force: boolean = false, emit: boolean = true) {
    if (this.historyStack.length === 0 && !force) {
      return;
    }

    const canvas = this.getCurrentUserCanvas();
    const ctx = this.getCurrentUserContext();
    if (canvas && ctx) {
      ctx.closePath();
      ctx.clearRect(0, 0, this.width, this.height);
      this.historyStack.length = 0;
      this.beginPath = false;
      this.mouseDown = false;

      if (emit) {
        this.emitDrawingUpdateEvent();
      }

      this.cancelFrame();
      this.rid = requestAnimationFrame(this.render);
    }
  }

  setStrokeStyle(style: string) {
    this.strokeStyle = style;
  }

  setStrokeWidth(width: number) {
    this.strokeWidth = width;
  }

  toggleDebug() {
    if (this.debugViewPort) {
      this.debugViewPort.enableDebug = !this.debugViewPort?.enableDebug;
    }
  }

  saveImage() {
    const canvas = this.getPrimaryCanvas();
    if (!canvas) return;

    const dl = document.createElement("a");
    dl.download = "draw-together.png";
    dl.href = canvas.toDataURL();
    dl.click();
  }

  getImageDataUrl() {
    const canvas = this.getPrimaryCanvas();
    if (!canvas) return;

    return canvas.toDataURL();
  }

  protected cancelFrame() {
    const rid = this.getRid();
    if (rid) {
      cancelAnimationFrame(rid);
      this.rid = null;
    }
  }

  protected set rid(value: number | null) {
    this.ridPriv = value;
  }

  protected get rid() {
    return this.ridPriv;
  }

  protected debug(elapsed: number) {
    const { x, y } = this.pos;
    this.debugViewPort?.debug([
      `posX: ${x}, posY: ${y}`,
      `width: ${this.width}`,
      `height: ${this.height}`,
      `offsetX: ${this.offsetX}`,
      `offsetY: ${this.offsetY}`,
      `-----`,
      `historyStackLength: ${this.historyStack.length}`,
      `maxHistory: ${this.maxHistory}`,
      `-----`,
      `strokeWidth: ${this.strokeWidth}`,
      `strokeStyle: ${this.strokeStyle}`,
      `-----`,
      `canvasChanged: ${this.canvasChanged}`,
      `mouseDown: ${this.mouseDown}`,
      `beginPath: ${this.beginPath}`,
      `-----`,
      `fpsCycle: ${this.fpsCycle}`,
      `FPS: ${Math.floor(1000 / elapsed)}`,
      `-----`,
      `pointerStack: ${this.pointerStack.join(",")}`,
      `pointerStackLength: ${this.pointerStack.length}`,
      `pointerLastHistory: ${this.pointerPositionEmit.join(",")}`,
      `prevX: ${this.prevX}`,
      `prevY: ${this.prevY}`,
      `pointerMoved: ${this.pointerMoved}`,
    ]);
  }
}
