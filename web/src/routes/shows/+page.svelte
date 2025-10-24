<script lang="ts">
  import CollectionCard from "$lib/components/CollectionCard.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import NewCollectionModal from "./NewCollectionModal.svelte";
  import Filter from "./Filter.svelte";
  import StandardPagination from "$lib/components/StandardPagination.svelte";
  import { Plus } from "lucide-svelte";
  import { isRoleAdmin } from "$lib/utils";

  const { data } = $props();

  let openNewCollectionModal = $state(false);
</script>

<Filter fullFilter={data.filter} />

<Spacer size="md" />

<div class="flex items-center justify-between">
  <h2 class="text-bold text-xl">
    Shows
    {#if isRoleAdmin(data.user?.role)}
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openNewCollectionModal = true;
        }}
      >
        <Plus />
      </Button>
    {/if}
  </h2>
  <p class="text-sm">{data.page.totalItems} show(s)</p>
</div>

<Spacer size="md" />

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
  <div
    class="grid grid-cols-[repeat(auto-fit,minmax(280px,1fr))] items-center justify-items-center gap-6"
  >
    {#each data.shows as show}
      <CollectionCard
        href="/shows/{show.id}"
        name={show.name}
        coverUrl={show.coverUrl}
      />
    {/each}
  </div>
</div>

<Spacer size="sm" />

<StandardPagination pageData={data.page} />

<NewCollectionModal bind:open={openNewCollectionModal} />
