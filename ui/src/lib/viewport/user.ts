import { BASE_VIEWPORT_HEIGHT, BASE_VIEWPORT_WIDTH } from ".";

type UserImageLayer = [OffscreenCanvas, OffscreenCanvasRenderingContext2D];

export class ViewPortUser {
  protected primaryCanvas: OffscreenCanvas;
  protected primaryCtx: OffscreenCanvasRenderingContext2D | null;

  protected userCanvasList: Record<string, OffscreenCanvas> = {};
  protected userCtxList: Record<string, OffscreenCanvasRenderingContext2D> = {};

  protected updateCycle = 1000 / 20;
  protected prevTs = 0;
  protected rid: number | null = null;

  constructor(
    protected width = BASE_VIEWPORT_WIDTH,
    protected height = BASE_VIEWPORT_HEIGHT
  ) {
    this.primaryCanvas = new OffscreenCanvas(this.width, this.height);
    this.primaryCtx = this.primaryCanvas.getContext("2d");
  }

  protected getUserImageLayer(userID: string): UserImageLayer {
    return [this.userCanvasList[userID], this.userCtxList[userID]];
  }

  addUser(userID: string) {
    if (this.userCanvasList[userID]) return;

    const canvas = new OffscreenCanvas(this.width, this.height);
    const ctx = canvas.getContext("2d");

    this.userCanvasList[userID] = canvas;
    if (ctx) {
      this.userCtxList[userID] = ctx;
    }
  }

  updateExternalUserCanvas(userID: string, imageData: ImageBitmap) {
    const [canvas, ctx] = this.getUserImageLayer(userID);

    if (ctx && canvas) {
      ctx.clearRect(0, 0, this.width, this.height);
      ctx.drawImage(imageData, 0, 0, this.width, this.height);
      imageData.close();
    }
  }

  render(now: number) {
    if (!this.primaryCtx) return;

    const elapsed = now - this.prevTs;

    if (elapsed > this.updateCycle) {
      this.updateCycle = now - (elapsed % this.updateCycle);
      this.primaryCtx.clearRect(0, 0, this.width, this.height);

      Object.keys(this.userCanvasList).forEach((u) => {
        const [canvas] = this.getUserImageLayer(u);

        this.primaryCtx?.drawImage(canvas, 0, 0, this.width, this.height);
      });
    }
  }

  getCanvas(): OffscreenCanvas {
    return this.primaryCanvas;
  }
}
