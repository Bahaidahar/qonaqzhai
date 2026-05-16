import type { Block } from "@/features/ai-chat/types";
import { PlanBlock } from "./blocks/plan-block";
import { BudgetBlock } from "./blocks/budget-block";
import { VendorsBlock } from "./blocks/vendors-block";

export function BlockRenderer({ block }: { block: Block }) {
  switch (block.type) {
    case "plan":
      return <PlanBlock data={block} />;
    case "budget":
      return <BudgetBlock data={block} />;
    case "vendors":
      return <VendorsBlock data={block} />;
  }
}
