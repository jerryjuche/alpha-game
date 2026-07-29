import { EventEmitter } from 'events';
import { Server } from '@stellar/stellar-sdk';

export interface ContractEvent {
  id: string;
  type: string;
  contractId: string;
  ledger: number;
  timestamp: number;
  data: Record<string, unknown>;
}

export class IndexerService extends EventEmitter {
  private rpcUrl: string;
  private pollInterval: number;
  private running = false;

  constructor(rpcUrl: string, pollInterval = 5000) {
    super();
    this.rpcUrl = rpcUrl;
    this.pollInterval = pollInterval;
  }

  async start() {
    this.running = true;
    await this.poll();
  }

  async stop() {
    this.running = false;
  }

  private async poll() {
    while (this.running) {
      try {
        const server = new Server(this.rpcUrl);
        const events = await this.fetchEvents(server);
        for (const event of events) {
          this.emit('event', event);
        }
      } catch (error) {
        this.emit('error', error);
      }
      await new Promise(resolve => setTimeout(resolve, this.pollInterval));
    }
  }

  private async fetchEvents(server: Server): Promise<ContractEvent[]> {
    return [];
  }
}
