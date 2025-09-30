<script lang="ts">
  import Stat from "$lib/components/Stat.svelte";
  import { Card } from "@nanoteck137/nano-ui";
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
  </div>
</Card.Root>
