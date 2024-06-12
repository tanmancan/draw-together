export class ViewPortDebug {
  protected debugCanvasCtx: CanvasRenderingContext2D | null = null;
  protected debugCanvas: HTMLCanvasElement | null = null;
  protected width: number = 400;
  protected height: number = 600;
  protected lineHeight: number = 18;

  constructor(public enableDebug: boolean = false) {
    if (!enableDebug) return;
    this.debugCanvas = document.createElement("canvas");
    this.debugCanvasCtx = this.debugCanvas.getContext("2d");

    this.debugCanvas.width = this.width;
    this.debugCanvas.height = this.height;
    this.debugCanvas.style.position = "absolute";
    this.debugCanvas.style.top = "80px";
    this.debugCanvas.style.right = "0";
    this.debugCanvas.style.background = "rgba(255,255,255,.5)";
    document.body.appendChild(this.debugCanvas);

    if (this.debugCanvasCtx) {
      this.debugCanvasCtx.font = `700 16px monospace`;
    }
  }

  debug(debugText: string[]) {
    if (!this.enableDebug && this.debugCanvas) {
      this.debugCanvas.style.opacity = "0";
      this.debugCanvas.style.visibility = "hidden";
    }
    if (!this.enableDebug || !this.debugCanvasCtx) return;

    this.debugCanvasCtx.clearRect(0, 0, this.width, this.height);

    debugText.forEach((t, i) => {
      this.debugCanvasCtx?.fillText(t, 8, this.lineHeight * i + 16);
    });
  }
}
