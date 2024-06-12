enum CircuitBreakerStatusEnum {
  CLOSED = "closed",
  HALF_OPEN = "half-open",
  OPEN = "open",
}

export class CircuitBreaker {
  protected status: CircuitBreakerStatusEnum = CircuitBreakerStatusEnum.CLOSED;

  protected failCount: number = 0;
  protected successCount: number = 0;
  protected lastFailTime: number = 0;

  constructor(
    protected name: string,
    protected debug = false,
    protected readonly failThreshHold: number = 5,
    protected readonly successThreshHold: number = 5,
    protected readonly resetTimeout: number = 30000
  ) {}

  get now() {
    return +new Date();
  }

  get state(): CircuitBreakerStatusEnum {
    if (
      this.failCount >= this.failThreshHold &&
      this.now - this.lastFailTime >= this.resetTimeout
    ) {
      return CircuitBreakerStatusEnum.HALF_OPEN;
    }

    if (this.failCount >= this.failThreshHold) {
      return CircuitBreakerStatusEnum.OPEN;
    }

    return CircuitBreakerStatusEnum.CLOSED;
  }

  get isClosed() {
    return this.status === CircuitBreakerStatusEnum.CLOSED;
  }

  get isHalfOpen() {
    return this.status === CircuitBreakerStatusEnum.HALF_OPEN;
  }

  get isOpen() {
    return this.status === CircuitBreakerStatusEnum.OPEN;
  }

  protected reset() {
    this.failCount = 0;
    this.successCount = 0;
    this.lastFailTime = 0;
  }

  protected fail_state() {
    this.failCount++;
    this.lastFailTime = this.now;
    this.successCount = 0;
  }

  protected success_state() {
    if (this.successCount >= this.successThreshHold) {
      this.reset();
    } else {
      this.successCount++;
    }
    this.failCount = 0;
  }

  protected logState() {
    if (this.debug) {
      console.log(
        `circuit-breaker: ${this.name} | state: ${this.state} | successCount: ${
          this.successCount
        } | errorCount: ${this.failCount} | this.lastFailTime: ${new Date(
          this.lastFailTime
        ).toLocaleDateString()}`
      );
    }
  }

  protect<T>(callable: () => Promise<T>) {
    return async () => {
      switch (this.state) {
        case CircuitBreakerStatusEnum.CLOSED:
        case CircuitBreakerStatusEnum.HALF_OPEN: {
          try {
            this.logState();
            if (typeof callable === "function") {
              const res = await callable();
              this.success_state();
              return res;
            }
            break;
          } catch (error) {
            this.fail_state();
            throw error;
          }
        }

        case CircuitBreakerStatusEnum.OPEN: {
          this.logState();
          break;
        }
        default: {
          this.logState();
          break;
        }
      }
    };
  }
}
