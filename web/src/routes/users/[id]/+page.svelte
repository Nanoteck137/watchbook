<script lang="ts">
  import { Card, Popover, Tooltip } from "@nanoteck137/nano-ui";
  import Stat from "./Stat.svelte";
  const { data } = $props();
</script>

<Card.Root>
  <div
    class="grid grid-cols-2 justify-around gap-6 p-6 md:grid-cols-3 lg:grid-cols-6"
  >
    <Stat stat={data.stats.all} />
    <Stat stat={data.stats.completed} />
    <Stat stat={data.stats.inProgress} />
    <Stat stat={data.stats.onHold} />
    <Stat stat={data.stats.dropped} />
    <Stat stat={data.stats.backlog} />

    <!-- <Tooltip.Provider>
      <Tooltip.Root>
        <Tooltip.Trigger class="flex flex-col items-center">
          <p class="text-xl">6</p>
          <p class="text-center text-xs">In Progress</p>
        </Tooltip.Trigger>
        <Tooltip.Content>
          <div class="flex items-center gap-2">
            <p>Anime: 100</p>
          </div>
        </Tooltip.Content>
      </Tooltip.Root>
    </Tooltip.Provider>

    <div class="flex flex-col items-center">
      <p class="text-xl">6</p>
      <p class="text-center text-xs">In Progress</p>
    </div>

    <div class="flex flex-col items-center">
      <p class="text-xl">6</p>
      <p class="text-center text-xs">In Progress</p>
    </div> -->
  </div>
</Card.Root>
