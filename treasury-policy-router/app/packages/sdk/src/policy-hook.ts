import { PolicyHook } from '../contracts/policy-hook';

export class PolicyHookClient {
  private client: PolicyHook;

  constructor(config: { rpcUrl: string; contractId: string }) {
    this.client = new PolicyHook({ rpcUrl: config.rpcUrl, contractId: config.contractId });
  }

  async initialize(admin: string) {
    return this.client.initialize(admin);
  }

  async setJurisdiction(country: string, allowed: boolean, operator: string) {
    return this.client.set_jurisdiction(country, allowed, operator);
  }

  async setDailyCap(cap: string, operator: string) {
    return this.client.set_daily_cap(cap, operator);
  }

  async setTravelRuleThreshold(threshold: string, operator: string) {
    return this.client.set_travel_rule_threshold(threshold, operator);
  }

  async checkPolicy(from: string, to: string, amount: string): Promise<PolicyResult> {
    return this.client.check_policy(from, to, amount);
  }

  async getJurisdiction(country: string): Promise<boolean> {
    return this.client.get_jurisdiction(country);
  }

  async getDailyCap(): Promise<string> {
    return this.client.get_daily_cap();
  }
}
